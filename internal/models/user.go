package models

import (
	"context"
	"log"
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
func InsertUser(user User) error {
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

// 通过用户id查找用户
func FindUserById(id string) *User {
	id_, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}

	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id_}}
	result := coll.FindOne(ctx, filter)

	var user User
	err = result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(err)
			return nil
		} else {
			log.Println(err)
		}
	}

	return &user
}

// 通过用户名查找用户
func FindUserByName(name string) []*User {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{Key: "name", Value: name}}
	cursor, err1 := coll.Find(ctx, filter)
	if err1 != nil {
		log.Println(err1)
	}

	var users []*User
	err2 := cursor.All(ctx, &users)
	if err2 != nil {
		return nil
	}

	return users
}

// 通过邮件来查找用户
func FindUserByEmail(email string) *User {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	result := coll.FindOne(ctx, filter)

	var user User
	err := result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(err)
			return nil
		} else {
			log.Println(err)
		}
	}

	return &user
}

// 查找所有用户
func FindAllUser() []*User {
	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cursor, err1 := coll.Find(ctx, bson.D{})
	if err1 != nil {
		return nil
	}

	var users []*User
	err2 := cursor.All(ctx, &users)
	if err2 != nil {
		return nil
	}

	return users
}

// 通过用户id更新用户
func UpdateUserById(id string, update bson.M) error {
	id_, err1 := bson.ObjectIDFromHex(id)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err2 := coll.UpdateByID(ctx, bson.D{{Key: "_id", Value: id_}}, update)
	if err2 != nil {
		return err2
	}

	return nil
}

// 通过用户id删除用户
func DeleteUserById(id string) error {
	id_, err1 := bson.ObjectIDFromHex(id)
	if err1 != nil {
		return err1
	}

	coll := db.Imgop.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err2 := coll.DeleteOne(ctx, bson.M{"_id": id_})
	if err2 != nil {
		return err2
	}

	return nil
}
