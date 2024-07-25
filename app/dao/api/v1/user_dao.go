package v1

import (
	"XArr-Rss/app/global/variable"
	"github.com/gin-gonic/gin"
	"time"
)

type Apiv1UserDao struct {
}

func (this Apiv1UserDao) Info(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	expireTime := ""
	if variable.ServerState.ExpireTime > 1956564045 {
		expireTime = "永久赞助会员"
	} else {
		expireTime = time.Unix(variable.ServerState.ExpireTime, 0).Format("2006-01-02 15:04:05")
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"version":           variable.ServerVersion,
			"program_file_time": variable.ProgramFileTime,
			"username":          variable.ServerState.Username,
			"email":             variable.ServerState.Email,
			"is_vip":            variable.ServerState.IsVip,
			"expire_time":       expireTime,
		},
		"message": "success",
	})
}
