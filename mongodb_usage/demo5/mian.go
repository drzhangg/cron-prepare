package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

//startTime小于某时间
//{"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		delCond    *DeleteCond
	)

	client, err = mongo.Connect(context.TODO(), "mongodb://47.99.240.52:27017", clientopt.ConnectTimeout(5*time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("cron")

	collection = database.Collection("log")

	delCond = &DeleteCond{TimeBeforeCond{time.Now().Unix()}}

	delResult,err := collection.DeleteMany(context.TODO(),delCond)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除的总行数：",delResult.DeletedCount)

}
