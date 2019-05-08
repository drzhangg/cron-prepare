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
		leaseID        clientv3.LeaseID
		kv             clientv3.KV
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
	)

	conf = clientv3.Config{
		Endpoints:   []string{"47.99.240.52:2379"},
		DialTimeout: 5 * time.Second,
	}

	client, err = clientv3.New(conf)
	if err != nil {
		fmt.Println(err)
	}

	//申请一个lease(租约)
	lease = clientv3.NewLease(client)

	//申请一个10秒的租约
	leaseGrantResp, err = lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约的ID
	leaseID = leaseGrantResp.ID


	ctx,_ := context.WithTimeout(context.TODO(),5 * time.Second)
	//开始续租
	keepRespChan, err = lease.KeepAlive(ctx, leaseID)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				}else {
					//每秒会续租一次，所以就会收到一次应答
					fmt.Println("收到自动续租应答：", keepResp.ID)
				}

			}
		}
	END:
	}()

	kv = clientv3.NewKV(client)

	//Put一个kv，让它与租约关联起来，从而实现10秒后自动过期
	putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseID))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功：", putResp.Header.Revision)

	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		//判断租约是否过期，过期的话Count为0
		if getResp.Count == 0 {
			fmt.Println("过期了")
			break
		}

		fmt.Println("还没过期：", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}

}
