package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"imgop/internal/models"
	"imgop/internal/utils"
)

// 添加好友申请
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

	type Req struct {
		Email string `json:"email" binding:"required"`
	}

	// 获取前端发送的好友邮箱和申请信息
	var req Req
	err2 := c.BindJSON(&req)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的请求格式"})
		return
	}

	// 查询好友的id
	friend := models.FindUserByEmail(req.Email)
	friend_id := friend.ID
	friend_id_str := friend.ID.Hex()

	// 数据库查验有没有现存好友关系，有的话就不创建了
	fs := models.FindFriendRequest(user_id_str, friend_id_str)
	if fs != nil {
		c.JSON(http.StatusOK, gin.H{"success": "ok", "message": "好友已存在或已申请"})
		return
	}

	// 添加一个新的好友关系
	friend_request := models.FriendRequest{FromId: user_id, ToId: friend_id}
	if err := models.InsertNewFriendRequest(friend_request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库插入错误"})
		return
	}

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

// 获取入站请求
func GetIncomingFriendRequests(c *gin.Context) {
	user_id_str, err1 := utils.GetUserIdStr(c) // 从认证信息中获取当前用户ID
	if err1 != nil {
		return
	}

	// 从数据库取出入站好友
	incominglist := models.FindIncomingFriendRequest(user_id_str)

	// 没找到
	if incominglist == nil {
		c.JSON(http.StatusOK, gin.H{"success": "ok", "message": "没有找到任何入站好友"})
		return
	}

	type friendInfo struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var incoming []friendInfo

	// 找到了进一步提出去入站请求中的好友信息
	for _, f := range incominglist {
		friend_id := f.FromId
		friend_id_str := friend_id.Hex()

		friend := models.FindUserById(friend_id_str)
		incoming = append(incoming, friendInfo{friend.Name, friend.Email})
	}

	// 返回包含好友信息的列表
	c.JSON(http.StatusOK, gin.H{
		"success":  "ok",
		"incoming": incoming,
	})
}

// 获取出站请求
func GetOutgoingFriendRequests(c *gin.Context) {
	user_id_str, err1 := utils.GetUserIdStr(c) // 从认证信息中获取当前用户ID
	if err1 != nil {
		return
	}

	// 从数据库取出出站好友
	outgoinglist := models.FindOutgoingFriendRequest(user_id_str)

	// 没找到
	if outgoinglist == nil {
		c.JSON(http.StatusOK, gin.H{"success": "ok", "message": "没有找到任何出站好友"})
		return
	}

	type friendInfo struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var outgoing []friendInfo

	// 找到了进一步提出去入站请求中的好友信息
	for _, f := range outgoinglist {
		friend_id := f.FromId
		friend_id_str := friend_id.Hex()

		friend := models.FindUserById(friend_id_str)
		outgoing = append(outgoing, friendInfo{friend.Name, friend.Email})
	}

	// 返回包含好友信息的列表
	c.JSON(http.StatusOK, gin.H{
		"success":  "ok",
		"outgoing": outgoing,
	})
}

// 同意好友申请
func AcceptFriendRequests(c *gin.Context) {
	// 获取当前用户id
	user_id_str, err1 := utils.GetUserIdStr(c)
	if err1 != nil {
		return
	}

	// 从前端请求中获取好友的邮箱
	type Req struct {
		Email string `json:"email"`
	}
	var req Req

	err2 := c.BindJSON(&req)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体出错"})
		return
	}

	friend_email := req.Email

	// 获取好友的id
	friend := models.FindUserByEmail(friend_email)
	if friend == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "找不到好友"})
		return
	}

	friend_id := friend.ID
	friend_id_str := friend_id.Hex()

	// 删除对应的好友申请
	err3 := models.DeleteFriendRequest(friend_id_str, user_id_str)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除请求失败"})
		return
	}

	// 将好友id添加到当前用户的好友列表
	err4 := models.UpdateUserById(user_id_str, bson.M{"$push": bson.M{"friend_list": friend_id}})
	if err4 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	user_id, err6 := bson.ObjectIDFromHex(user_id_str)
	if err6 != nil {
		return
	}

	// 将自己添加到用户的好友列表
	err5 := models.UpdateUserById(friend_id_str, bson.M{"$push": bson.M{"friend_list": user_id}})
	if err5 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

// 拒绝好友申请
func RejectFriendRequests(c *gin.Context) {
	// 获取当前用户id
	user_id_str, err1 := utils.GetUserIdStr(c)
	if err1 != nil {
		return
	}

	// 从前端请求中获取好友的邮箱
	type Req struct {
		Email string `json:"email"`
	}
	var req Req

	err2 := c.BindJSON(&req)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体出错"})
		return
	}

	friend_email := req.Email

	// 获取好友的id
	friend := models.FindUserByEmail(friend_email)
	if friend == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "找不到好友"})
		return
	}

	friend_id := friend.ID
	friend_id_str := friend_id.Hex()

	// 删除对应的好友申请
	err3 := models.DeleteFriendRequest(user_id_str, friend_id_str)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除请求失败"})
		return
	}

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

// 获取用户所有的好友
func GetAllFriends(c *gin.Context) {
	// 获取用户id
	user_id_str, err := utils.GetUserIdStr(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户id字符串时出错"})
		return
	}

	// 获取当前用户以及好友列表
	user := models.FindUserById(user_id_str)
	friend_list := user.FriendList
	if len(friend_list) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": "ok", "message": "没有好友"})
		return
	}

	// 遍历好友列表，提取好友信息，塞入应答字段
	type FriendInfo struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var friends []FriendInfo
	for _, f_id := range friend_list {
		f_id_str := f_id.Hex()
		f := models.FindUserById(f_id_str)
		friends = append(friends, FriendInfo{f_id_str, f.Name, f.Email})
	}

	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"success": "ok",
		"friends": friends,
	})
}

// 搜索某个好友
func SearchFriend(c *gin.Context) {
	// 获取用户id
	user_id_str, err := utils.GetUserIdStr(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户id字符串时出错"})
	}

	// 获取好友邮箱
	type Req struct {
		Email string `json:"email"`
	}
	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的请求"})
	}

	// 通过好友邮箱来获取好友id
	friend := models.FindUserByEmail(req.Email)
	friend_id_str := friend.ID.Hex()

	// 查询好友是否在好友列表里
	if !models.IsFriend(user_id_str, friend_id_str) {
		c.JSON(http.StatusOK, gin.H{"success": "ok", "message": "不存在这个好友"})
		return
	}

	// 将好友信息提取出来，然后返回
	c.JSON(http.StatusOK, gin.H{
		"success":      "ok",
		"friend_id":    friend_id_str,
		"friend_name":  friend.Name,
		"friend_email": friend.Email,
	})
}
