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

	//初始化连接
	conf = clientv3.Config{
		Endpoints:[]string{"47.99.240.52:2379"},
		DialTimeout:5 * time.Second,
	}

	//创建etcd连接
	client,err = clientv3.New(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	getResp,err = kv.Get(context.TODO(),"/cron/jobs/job1")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(getResp.Kvs)
}
