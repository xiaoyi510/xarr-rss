package v1

import (
	"github.com/gin-gonic/gin"
)

type ApiV1OneSaidDao struct {
}

func (this ApiV1OneSaidDao) Get(c *gin.Context) {
	//curl := helper.GetCurlHttpHelperDefault()
	//_, response := curl.Get("https://v1.hitokoto.cn/", nil, false)
	//
	//res := &resHitokoto{}
	//err := json.Unmarshal(response, res)
	//if err != nil {
	//	return
	//}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": "美好的一天从现在开始", //res.Hitokoto,
	})
}

type resHitokoto struct {
	Id         int         `json:"id"`
	Uuid       string      `json:"uuid"`
	Hitokoto   string      `json:"hitokoto"`
	Type       string      `json:"type"`
	From       string      `json:"from"`
	FromWho    interface{} `json:"from_who"`
	Creator    string      `json:"creator"`
	CreatorUid int         `json:"creator_uid"`
	Reviewer   int         `json:"reviewer"`
	CommitFrom string      `json:"commit_from"`
	CreatedAt  string      `json:"created_at"`
	Length     int         `json:"length"`
}
