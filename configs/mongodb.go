package configs

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoDB *mongo.Client

func ConnectMongoDB() {
	opts := options.Client().ApplyURI("mongodb://0.0.0.0:27017")

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	MongoDB = client

	log.Println("Monogdb 连接成功")

}

func DisconnectMongoDB() {
	err := MongoDB.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}

	log.Println("Monogdb 断开连接成功")
}
