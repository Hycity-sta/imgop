package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"imgop/internal/db"
)

type OnlineUser struct {
	UserId bson.ObjectID `json:"user_id" bson:"user_id"`
}

func InsertOnlineUser(user_id_str string) error {
	user_id, err1 := bson.ObjectIDFromHex(user_id_str)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("online_user_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err2 := coll.InsertOne(ctx, user_id)
	if err2 != nil {
		return err2
	}

	return nil
}

func DeleteOnlineUser(user_id_str string) error {
	user_id, err1 := bson.ObjectIDFromHex(user_id_str)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("online_user_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.M{"user_id": user_id}
	_, err2 := coll.DeleteOne(ctx, filter)
	if err2 != nil {
		return err2
	}

	return nil
}

// 查找好友是否在线
func FindUserIsOnline(user_id_str string) (bool, error) {
	user_id, err1 := bson.ObjectIDFromHex(user_id_str)
	if err1 != nil {
		return false, err1
	}

	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: user_id}}
	count, err2 := coll.CountDocuments(ctx, filter)
	if err2 != nil {
		return false, err2
	}

	if count != 1 {
		return false, fmt.Errorf("在线好友树出现错误")
	}

	return true, nil
}
