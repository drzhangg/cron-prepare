package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		conf               clientv3.Config
		client             *clientv3.Client
		err                error
		kv                 clientv3.KV
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
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

	kv = clientv3.NewKV(client)

	go func() {
		for {

			kv.Put(context.TODO(), "/cron/jobs/job7", "it is job7")

			kv.Delete(context.TODO(), "/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7")
	if err != nil {
		fmt.Println(err)
		return
	}

	//现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值：", string(getResp.Kvs[0].Value))
	}

	watchStartRevision = getResp.Header.Revision + 1

	//创建一个watcher
	watcher = clientv3.NewWatcher(client)

	//启动监听
	fmt.Println("从该版本向后监听：", watchStartRevision)

	watchRespChan = watcher.Watch(context.TODO(), "/corn/jobs/job7", clientv3.WithRev(watchStartRevision))

	//处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}
