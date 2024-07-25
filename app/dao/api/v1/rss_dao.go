package v1

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/model"
	"XArr-Rss/util/logsys"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type RssDao struct {
}

func (this RssDao) Result(c *gin.Context) {
	groupFile := c.Param("groupId.xml")
	mediaId := c.Query("mediaId")

	// 找到group Id
	groupId := strings.Replace(groupFile, ".xml", "", -1)
	if mediaId != "" {
		groupFile = groupId + "/" + mediaId + ".xml"
	}
	// 判断是否需要刷新
	refresh := c.Query("refresh")
	if refresh == "1" {
		// 先来波刷新
		err := ApiV1Groups{}._refresh(groupId)
		if err != nil {
			//continue
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"message": "刷新失败:" + err.Error(),
			})
			return
		}
	}

	_, err := os.Stat("./conf/trans/group_" + groupFile)
	if os.IsNotExist(err) {
		c.String(404, "文件不存在")
		return
	}
	if err != nil {
		logsys.Error("访问Rss文件错误:%s", err.Error())
		c.String(200, "错误的访问:"+err.Error())
		return
	}

	// 返回文件
	c.File("./conf/trans/group_" + groupFile)
}

func (this RssDao) getDefaultXml() *model.RssRoot {

	groupAll := &model.RssRoot{Version: "2.0"}
	groupAll.Channel.Link = appconf.AppConf.System.HttpAddr + "/rss/group_all.xml"
	groupAll.Channel.Description = "XArr-Rss"
	groupAll.Channel.Title = "XArr-Rss"
	groupAll.Channel.Item = append(groupAll.Channel.Item, model.RssResultItem{
		Title: model.CDATA{
			Text: "XArr-Rss 没有找到",
		},
		OriginalTitle: model.CDATA{
			Text: "XArr-Rss 没有找到",
		},
		OtherTitle: model.CDATA{
			Text: "",
		},
		PubDate: "2022-04-21T05:38:00",
		Enclosure: model.RssResultItemEnclosure{
			Type:   "application/x-bittorrent",
			Length: "329672288",
			Url:    "https://xarr.52nyg.com",
		},
		Link: "https://xarr.52nyg.com",
		Guid: model.RssResultItemGuid{
			IsPermaLink: false,
			Text:        "https://xarr.52nyg.com",
		},
		XArrRssIndexer: model.RssResultItemXArrRssIndexer{
			Text: "",
			ID:   "",
		},
	})

	return groupAll

}
