package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

func main() {

	//1.建立连接
	client, err := mongo.Connect(context.TODO(), "mongodb://47.99.240.52:27017", clientopt.ConnectTimeout(5*time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	//2.选择数据库my_db
	database := client.Database("my_db")

	//3.选择表my_collection
	collection := database.Collection("my_collection")
	collection = collection
	//fmt.Println(collection)
}
