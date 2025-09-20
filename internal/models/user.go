package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"imgop/internal/db"
)

type User struct {
	ID           bson.ObjectID `json:"id" bson:"_id,omitempty"`            // 用户的唯一标识符（MongoDB自动生成）
	Name         string        `json:"name" bson:"name"`                   // 用户登录名
	Email        string        `json:"email" bson:"email"`                 // 用户邮箱（唯一，用于找回密码等）
	PasswordHash string        `json:"password_hash" bson:"password_hash"` // 加密后的密码（绝对不能存明文）
	Status       string        `json:"status" bson:"status"`               // 在线状态：online/offline/away
	LastSeen     time.Time     `json:"last_seen" bson:"last_seen"`         // 最后活动时间（用于判断离线超时）
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`       // 账号创建时间
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`       // 资料最后更新时间
}

// 插入一个新的用户
func InsertUser(user *User) error {
	// 设置时间戳
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// 确保ID为空时会自动生成
	if user.ID == bson.NilObjectID {
		user.ID = bson.NewObjectID()
	}

	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := coll.InsertOne(ctx, user)
	return err
}

// 插入许多用户
func InsertManyUser(users *[]User) (*mongo.InsertManyResult, error) {
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
