package v1

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/model/webmenu"
	"XArr-Rss/app/service/groups"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

type ApiV1Menu struct {
}

func (this ApiV1Menu) Get(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	//file, err := os.ReadFile("./web/data/menu.json")
	//if err != nil {
	//	c.JSON(200, gin.H{
	//		"data": nil,
	//		"message":  "找不到菜单信息",
	//		"code": 500,
	//	})
	//	return
	//}
	//var menu []*webmenu.WebMenu
	menu := appconf.AppConf.Menu
	// 添加group菜单
	for _, v := range menu {
		of := reflect.TypeOf(v.Id)
		if of.Kind() == reflect.Float64 {
			if v.Id.(float64) == 300 {
				// 查询groupList
				groupList := groups.GroupsService{}.GetGroupList()
				v.Children = v.Children[:2]
				for _, item := range groupList {
					v.Children = append(v.Children, webmenu.MenuChildren{
						Id:       302 + item.Id,
						Title:    item.Name,
						Icon:     "",
						Type:     1,
						OpenType: "_iframe",
						Href:     "view/groups/medias-list.html?groupId=" + strconv.Itoa(int(item.Id)),
					})
				}
			}
		}

	}

	c.JSON(200, menu)
}
