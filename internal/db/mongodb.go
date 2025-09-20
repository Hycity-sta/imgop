package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoDB *mongo.Client
var Imgop *mongo.Database

func ConnectMongoDB() {
	opts := options.Client().ApplyURI("mongodb://0.0.0.0:27017")

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	MongoDB = client
	Imgop = MongoDB.Database("imgop")

	initUserCollection()

	log.Println("Monogdb 连接成功")

}

func DisconnectMongoDB() {
	err := MongoDB.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}

	log.Println("Monogdb 断开连接成功")
}

func initUserCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查索引是否已存在
	cur, err := Imgop.Collection("users").Indexes().List(ctx)
	if err != nil {
		log.Fatalf("Failed to list indexes: %v", err)
	}

	var indexes []bson.M
	if err = cur.All(ctx, &indexes); err != nil {
		log.Fatalf("Failed to decode indexes: %v", err)
	}

	emailIndexExists := false
	for _, idx := range indexes {
		if key, ok := idx["key"].(bson.M); ok {
			if _, exists := key["email"]; exists {
				emailIndexExists = true
				break
			}
		}
	}

	// 设置email为唯一索引
	if !emailIndexExists {
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		}

		_, err := Imgop.Collection("users").Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			log.Fatalf("Failed to create unique index: %v", err)
		}
		log.Println("Created unique index on email field")
	}
}
