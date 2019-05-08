package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		conf           clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		ctx            context.Context
		cancelFunc     context.CancelFunc
		kv             clientv3.KV
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
	)

	conf = clientv3.Config{
		Endpoints:   []string{"47.99.240.52:2379"},
		DialTimeout: 5 * time.Second,
	}

	client, err = clientv3.New(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	//lease实现锁自动过期：
	//op操作
	//txn事务：if else then

	//1.上锁（创建租约，自动续租，拿着租约去抢占一个key）
	//创建一个租约
	lease = clientv3.NewLease(client)

	//申请一个5秒租约
	leaseGrantResp, err = lease.Grant(context.TODO(), 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约ID
	leaseId = leaseGrantResp.ID

	//准备一个用于取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())

	//确保函数退出后，自动续租会停止
	defer cancelFunc()
	//函数退出后，立刻销毁租约
	defer lease.Revoke(context.TODO(), leaseId)

	keepRespChan, err = lease.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}

	//处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else {
					fmt.Println("收到自动续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()

	//创建事务
	kv = clientv3.NewKV(client)

	//创建事务
	txn = kv.Txn(context.TODO())

	//如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/kob9")) //否则抢锁失败


	txnResp,err = txn.Commit()
	if err != nil {
		fmt.Println(err)
		return
	}

	if !txnResp.Succeeded {
		fmt.Println("🔐被占用：",string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}


	//2.处理业务
	fmt.Println("处理任务")
	time.Sleep(5 * time.Second)

	//3.释放锁(取消自动续租，释放租约)

}
