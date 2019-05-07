package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {

	var(
		conf clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		delResp *clientv3.DeleteResponse
		kvpair *mvccpb.KeyValue
	)

	conf = clientv3.Config{
		Endpoints:[]string{"47.99.240.52:2379"},
		DialTimeout:5 * time.Second,
	}

	client,err = clientv3.New(conf)
	if err != nil {
		fmt.Println(err)
	}

	kv = clientv3.NewKV(client)

	delResp,err = kv.Delete(context.TODO(),"/cron/jobs/job2",clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}

	//被删除之前的value是什么
	if len(delResp.PrevKvs) != 0 {
		for _,kvpair = range delResp.PrevKvs{
			fmt.Println("删除了：",string(kvpair.Key),string(kvpair.Value))
		}
	}
}
