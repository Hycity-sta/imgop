package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"imgop/internal/models"
	"imgop/internal/utils"
)

type GetFriendshipResponse struct {
	Friends []*models.Friendship `json:"friends"`
}

// 获取用户所有的好友
func GetAllFriends(c *gin.Context) {
	user_id, err := utils.GetUserIdStr(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户id字符串时出错"})
	}

	friends := models.FindFriendShipById(user_id)

	c.JSON(http.StatusOK, GetFriendshipResponse{
		Friends: friends,
	})
}

type SearchFriendsRequest struct {
	FriendEmail string `json:"friend_email" binding:"required"`
}

func SearchFriend(c *gin.Context) {
	// 获取用户id
	user_id, err := utils.GetUserIdStr(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户id字符串时出错"})
	}

	// 获取好友邮箱
	var req SearchFriendsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的请求"})
	}

	// 通过好友邮箱来获取好友id
	friend := models.FindUserByEmail(req.FriendEmail)
	friend_id := friend.ID.Hex()

	// 通过好友id和用户id获取好友关系
	friendship := models.FindFriendShipBy2Id(user_id, friend_id)
	if friendship != nil {
		c.JSON(http.StatusOK, gin.H{"friendship": friendship})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到此好友"})
	}
}

type AddFriendRequest struct {
	Email string `json:"email" binding:"required"`
}

func AddFriend(c *gin.Context) {
	// 获取用户的id
	user_id_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt验证出错"})
		return
	}

	// 转换字符串
	user_id_str, ok := user_id_.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "user_id is not a string"})
		return
	}

	// 转换成objectID
	user_id, err1 := bson.ObjectIDFromHex(user_id_str)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "转换失败"})
		return
	}

	// 获取前端发送的好友邮箱和申请信息
	var req AddFriendRequest
	err2 := c.BindJSON(&req)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的请求格式"})
		return
	}

	// 查询好友的id
	friend := models.FindUserByEmail(req.Email)
	friend_id := friend.ID

	// 添加一个新的好友关系
	friendship_a := models.Friendship{
		UserId:   user_id,
		FriendId: friend_id,
		Status:   "pending",
	}

	// 同时，被申请的人也要添加一个新的好友关系
	friendship_b := models.Friendship{
		UserId:   friend_id,
		FriendId: user_id,
		Status:   "pending",
	}

	if err := models.InsertNewFriendShip(friendship_a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库插入错误"})
		return
	}

	if err := models.InsertNewFriendShip(friendship_b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库插入错误"})
		return
	}

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

type AcceptFriendRequest struct {
	Email string
}

// 接受好友申请
func AcceptFriend(c *gin.Context) {
	// 获取用户id

	// 获取发送请求的好友id

	// 更新用户这一方的好友关系的状态

	// 更新好友这一方的好友关系的状态

	// 返回成功
}

type RejectFriendRequest struct {
}

// 拒绝好友申请
func RejectFriend(c *gin.Context) {

}

type GetAllPendingFriendResponse struct {
	Friends []*models.Friendship `json:"friends"`
}

// 获取所有待同意的好友
func GetAllPendingFriend(c *gin.Context) {
	// 获取用户id
	user_id, err := utils.GetUserIdStr(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户id字符串时出错"})
	}

	// 数据库查询
	friends := models.FindPendingFriendShipById(user_id)

	if friends != nil {
		c.JSON(http.StatusOK, GetAllPendingFriendResponse{Friends: friends})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到待通过好友申请的好友"})
	}
}
