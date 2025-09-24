package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"imgop/internal/db"
)

type FriendRequest struct {
	FromId bson.ObjectID `json:"from_id" bson:"from_id"`
	ToId   bson.ObjectID `json:"to_id" bson:"to_id"`
}

func InsertNewFriendRequest(fr FriendRequest) error {
	coll := db.Imgop.Collection("friend_request_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := coll.InsertOne(ctx, fr)
	return err
}

func DeleteFriendRequest(from_id_str string, to_id_str string) error {
	from_id, err1 := bson.ObjectIDFromHex(from_id_str)
	if err1 != nil {
		return err1
	}

	to_id, err2 := bson.ObjectIDFromHex(to_id_str)
	if err2 != nil {
		return err2
	}

	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.M{"from_id": from_id, "to_id": to_id}
	_, err3 := coll.DeleteOne(ctx, filter)
	if err3 != nil {
		return err3
	}

	return nil
}

func FindFriendRequest(from_id_str string, to_id_str string) *FriendRequest {
	// 打开集合
	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 字符串转换成数据库id类型
	from_id, err1 := bson.ObjectIDFromHex(from_id_str)
	if err1 != nil {
		return nil
	}

	to_id, err2 := bson.ObjectIDFromHex(to_id_str)
	if err2 != nil {
		return nil
	}

	// 查找
	filter := bson.M{"from_id": from_id, "to_id": to_id}
	cursor, err3 := coll.Find(ctx, filter)
	if err3 != nil {
		return nil
	}

	// 编码到列表中
	var friend_request *FriendRequest
	err4 := cursor.All(ctx, &friend_request)
	if err4 != nil {
		return nil
	}

	return friend_request
}

func FindIncomingFriendRequest(user_id_str string) []*FriendRequest {
	// 打开集合
	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 字符串转换成数据库id类型
	user_id, err1 := bson.ObjectIDFromHex(user_id_str)
	if err1 != nil {
		return nil
	}

	// 查找
	filter := bson.M{"to_id": user_id}
	cursor, err2 := coll.Find(ctx, filter)
	if err2 != nil {
		return nil
	}

	// 编码到列表中
	var friend_request_list []*FriendRequest
	err3 := cursor.All(ctx, &friend_request_list)
	if err3 != nil {
		return nil
	}

	return friend_request_list
}

func FindOutgoingFriendRequest(user_id_str string) []*FriendRequest {
	// 打开集合
	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 字符串转换成数据库id类型
	user_id, err1 := bson.ObjectIDFromHex(user_id_str)
	if err1 != nil {
		return nil
	}

	// 查找
	filter := bson.M{"from_id": user_id}
	cursor, err2 := coll.Find(ctx, filter)
	if err2 != nil {
		return nil
	}

	// 编码到列表中
	var friend_request_list []*FriendRequest
	err3 := cursor.All(ctx, &friend_request_list)
	if err3 != nil {
		return nil
	}

	return friend_request_list
}
