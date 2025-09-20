package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"imgop/internal/db"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty"` // 用户的唯一标识符（MongoDB自动生成）
	Username     string        `bson:"username"`      // 用户登录名（唯一）
	Email        string        `bson:"email"`         // 用户邮箱（唯一，用于找回密码等）
	PasswordHash string        `bson:"password_hash"` // 加密后的密码（绝对不能存明文）
	Status       string        `bson:"status"`        // 在线状态：online/offline/away
	LastSeen     time.Time     `bson:"last_seen"`     // 最后活动时间（用于判断离线超时）
	CreatedAt    time.Time     `bson:"created_at"`    // 账号创建时间
	UpdatedAt    time.Time     `bson:"updated_at"`    // 资料最后更新时间
}

// 插入一个新的用户
func InsertUser(user User) (*mongo.InsertOneResult, error) {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, user)
	return result, err
}

// 插入许多用户
func InsertManyUser(users []User) (*mongo.InsertManyResult, error) {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := coll.InsertMany(ctx, users)
	return result, err
}

// 通过用户id查找用户
func FindUserById(id bson.ObjectID) *mongo.SingleResult {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}})

	return result
}

// 通过用户名查找用户
func FindUserByName(name string) (*mongo.Cursor, error) {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := coll.Find(ctx, bson.D{{Key: "name", Value: name}})
	return result, err
}

// 查找所有用户
func FindAllUser() (*mongo.Cursor, error) {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := coll.Find(ctx, bson.D{})
	return result, err
}

// 通过用户id更新用户
func UpdateUserById(id bson.ObjectID, update any) *mongo.SingleResult {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := coll.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: id}}, update)
	return result
}

// 通过用户id删除用户
func DeleteUserById(id bson.ObjectID) *mongo.SingleResult {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := coll.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: id}})
	return result
}
