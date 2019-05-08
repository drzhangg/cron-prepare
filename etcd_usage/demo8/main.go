package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {

	var (
		conf   clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		opResp clientv3.OpResponse
		getOp  clientv3.Op
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

	//创建OP：operation
	putOp = clientv3.OpPut("/cron/jobs/job8", "123123123")

	//执行op
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入Revision：", opResp.Put().Header.Revision)

	//创建op
	getOp = clientv3.OpGet("/cron/jobs/job8")

	//执行op
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("数据Revision:", opResp.Get().Kvs[0].ModRevision)	// create rev == mod rev
	fmt.Println("数据value:", string(opResp.Get().Kvs[0].Value))
}
