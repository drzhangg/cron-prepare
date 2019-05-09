package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		record     *LogRecord
		logArr     []interface{}
		result     *mongo.InsertManyResult
		insertId   interface{}
		docId      objectid.ObjectID
	)

	client, err = mongo.Connect(context.TODO(), "mongodb://47.99.240.52:27017", clientopt.ConnectTimeout(5*time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	//连接数据库
	database = client.Database("cron")

	//连接表
	collection = database.Collection("log")

	record = &LogRecord{
		JobName: "job11",
		Command: "echo i want a job",
		Err:     "no job",
		Content: "i need a job",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	logArr = []interface{}{record, record, record}

	//插入多条数据
	result, err = collection.InsertMany(context.TODO(), logArr)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, insertId = range result.InsertedIDs {
		docId = insertId.(objectid.ObjectID)
		fmt.Println("自增ID：",docId.Hex())
	}
}
