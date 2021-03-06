package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		conf    clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		putResp *clientv3.PutResponse
	)

	conf = clientv3.Config{
		Endpoints:   []string{"47.99.240.52:2379"},
		DialTimeout: 5 * time.Second,
	}

	//建立一个客户端
	if client, err = clientv3.New(conf); err != nil {
		fmt.Println(err)
		return
	}

	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hello",clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Revision:",putResp.Header.Revision)
	if putResp.PrevKv != nil {
		fmt.Println("PrevValue:",string(putResp.PrevKv.Value))
	}
}
