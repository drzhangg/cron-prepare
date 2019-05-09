package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

//任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

//记录一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"` //任务名
	Command   string    `bson:"command"` //shell命令
	Err       string    `bson:"err"`     //脚本错误
	Content   string    `bson:"content"` //脚本输出内容
	TimePoint TimePoint `bson:"timePoint""`
}

func main() {

	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		record     *LogRecord
		result     *mongo.InsertOneResult
		objId      objectid.ObjectID
	)

	//创建连接
	client, err = mongo.Connect(context.TODO(), "mongodb://47.99.240.52:27017",clientopt.ConnectTimeout(5 * time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	//选择数据库
	database = client.Database("cron")

	//选择表
	collection = database.Collection("log")

	record = &LogRecord{
		JobName: "job10",
		Command: "echo hello 10",
		Err:     "",
		Content: "hello 10",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	result, err = collection.InsertOne(context.TODO(), record)
	if err != nil {
		fmt.Println(err)
		return
	}

	objId = result.InsertedID.(objectid.ObjectID)
	fmt.Println("自增ID：", objId.Hex())
}
