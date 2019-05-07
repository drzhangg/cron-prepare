package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var(
		conf clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
	)

	conf = clientv3.Config{
		Endpoints:[]string{"47.99.240.52:2379"},
		DialTimeout:5 * time.Second,
	}

	client,err = clientv3.New(conf)
	if err!=nil{
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	 kv.Put(context.TODO(),"/cron/jobs/job2","world")

	getResp,err = kv.Get(context.TODO(),"/cron/jobs/",clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(getResp.Kvs)
}
