package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"imgop/internal/db"
)

type UserStatus struct {
	UserId     bson.ObjectID `json:"user_id" bson:"user_id"`
	Online     bool          `json:"online" bson:"online"`
	LastActive time.Time     `json:"last_active" bson:"last_active"`
}

func InsertNewUserStatus(user_status UserStatus) error {
	coll := db.Imgop.Collection("user_status_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, error := coll.InsertOne(ctx, &user_status)
	if error != nil {
		return error
	}

	return nil
}

func FindUserStatusById(user_id string) []*UserStatus {
	// 转换id为字符串
	id_, err1 := bson.ObjectIDFromHex(user_id)
	if err1 != nil {
		return nil
	}

	// 打开集合并设置ctx
	coll := db.Imgop.Collection("user_status_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 数据库查询
	filter := bson.M{"user_id": id_}
	cursor, err2 := coll.Find(ctx, filter)
	if err2 != nil {
		return nil
	}

	// 解码数据到对象上
	var userStatusList []*UserStatus
	err3 := cursor.All(ctx, &userStatusList)
	if err3 != nil {
		return nil
	}

	return userStatusList
}

func UpdateUserStatusById(user_id string, update bson.M) error {
	id_, err1 := bson.ObjectIDFromHex(user_id)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("user_status_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err2 := coll.UpdateByID(ctx, bson.D{{Key: "user_id", Value: id_}}, update)
	if err2 != nil {
		return err2
	}

	return nil
}

func DeleteUserStatusById(user_id string) error {
	id_, err1 := bson.ObjectIDFromHex(user_id)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("user_status_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err2 := coll.DeleteOne(ctx, bson.M{"user_id": id_})
	if err2 != nil {
		return err2
	}

	return nil
}
