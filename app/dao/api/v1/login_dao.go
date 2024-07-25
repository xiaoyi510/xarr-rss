package v1

import (
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/apiv1/login"
	"XArr-Rss/util/hash"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type Apiv1LoginDao struct {
}

// Login 登录后台
func (this Apiv1LoginDao) Login(c *gin.Context) {
	req := login.ApiV1LoginReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Username) == 0 || len(req.Password) == 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请输入账号密码",
		})
		return
	}

	variable.ServerState.ExpireTime = time.Now().Unix() + 365*60*60*24
	variable.ServerState.Email = "admin@52nyg.com"
	variable.ServerState.Username = "admin"
	variable.ServerState.LicenseKey = "asfasf"

	variable.ServerState.RefreshTime = time.Now().Unix()

	variable.ServerState.IsVip = true

	// 账号密码正确 返回cookie
	token := genToken()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "xarn_rss_token" + variable.HttpTokenAfter,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"token": token,
		},
		"message": "登录成功",
	})
}

func genToken() string {

	text := "奥斯房间奥斯房间奥+++++"
	text = text + time.Now().String() + "asfjoasfij"

	k := hash.Md5{}.Hash([]byte(text))

	variable.LoginToken.Store(k, 60*60*24)

	return k
}
