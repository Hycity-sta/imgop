package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"imgop/internal/db"
)

type Friendship struct {
	UserId   bson.ObjectID `json:"user_id" bson:"user_id"`
	FriendId bson.ObjectID `json:"friend_id" bson:"friend_id"`
	Status   string        `json:"status" bson:"status"` // 好友关系 accepted rejected pending
}

// 插入一条新的好友关系
func InsertNewFriendShip(fs Friendship) error {
	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := coll.InsertOne(ctx, fs)
	return err
}

// 通过user_id查找好友关系
func FindFriendShipById(user_id string) []*Friendship {
	// 打开集合
	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 字符串转换成数据库id类型
	id_, err := bson.ObjectIDFromHex(user_id)
	if err != nil {
		return nil
	}

	// 查找好友关系
	filter := bson.M{"user_id": id_}
	cursor, err1 := coll.Find(ctx, filter)
	if err1 != nil {
		return nil
	}

	// 编码到列表中
	var friendship_list []*Friendship
	err2 := cursor.All(ctx, &friendship_list)
	if err2 != nil {
		return nil
	}

	return friendship_list
}

// 查找两个id确认唯一一条好友关系
func FindFriendShipBy2Id(user_id string, friend_id string) *Friendship {
	// 打开集合
	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 字符串转换成数据库id类型
	user_id_, err1 := bson.ObjectIDFromHex(user_id)
	if err1 != nil {
		return nil
	}

	friend_id_, err2 := bson.ObjectIDFromHex(friend_id)
	if err2 != nil {
		return nil
	}

	// 查找好友关系
	filter := bson.M{"user_id": user_id_, "friend_id": friend_id_}
	cursor, err3 := coll.Find(ctx, filter)
	if err3 != nil {
		return nil
	}

	// 编码到列表中
	var friendship *Friendship
	err4 := cursor.All(ctx, &friendship)
	if err4 != nil {
		return nil
	}

	return friendship
}

// 查找当前用户的待同意的好友
func FindPendingFriendShipById(user_id string) []*Friendship {
	// 打开集合
	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 字符串转换成数据库id类型
	id_, err := bson.ObjectIDFromHex(user_id)
	if err != nil {
		return nil
	}

	// 查找好友关系
	filter := bson.M{"user_id": id_, "status": "pending"}
	cursor, err1 := coll.Find(ctx, filter)
	if err1 != nil {
		return nil
	}

	// 编码到列表中
	var friendship_list []*Friendship
	err2 := cursor.All(ctx, &friendship_list)
	if err2 != nil {
		return nil
	}

	return friendship_list
}

// 通过用户id更新好友关系
func UpdateFriendShipById(user_id string, update any) error {
	id_, err1 := bson.ObjectIDFromHex(user_id)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err2 := coll.UpdateByID(ctx, bson.D{{Key: "_id", Value: id_}}, update)
	if err2 != nil {
		return err2
	}

	return nil
}

// 通过用户id删除好友关系
func DeleteFriendShipById(user_id string) error {
	id_, err1 := bson.ObjectIDFromHex(user_id)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("friendship_coll")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err2 := coll.DeleteOne(ctx, bson.M{"_id": id_})
	if err2 != nil {
		return err2
	}

	return nil
}
