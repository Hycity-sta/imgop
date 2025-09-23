package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserIdStr(c *gin.Context) (string, error) {
	id, exist := c.Get("user_id")
	if !exist {

		return "", fmt.Errorf("意外出错")
	}

	// 类型断言为 string
	user_id, ok := id.(string)
	if !ok {
		return "", fmt.Errorf("类型转换出错")
	}

	return user_id, nil
}
