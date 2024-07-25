package v1

import (
	"XArr-Rss/app/service/download"
	"XArr-Rss/app/service/groups"
	"XArr-Rss/app/service/medias"
	"github.com/gin-gonic/gin"
)

type Apiv1StaticsDao struct {
}

func (this Apiv1StaticsDao) Info(c *gin.Context) {
	mediaCount := medias.MediaService{}.GetMediaCount()

	////////////// 下载统计
	down := this.getStaticsDown()
	////////////// 分组统计
	group := this.getStaticsGroup()

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"medias": gin.H{
				"sonarr": gin.H{
					"count": mediaCount,
				},
			},
			"down":  down,
			"group": group,
			"group_medias": gin.H{
				"total":     groups.GroupMediaService{}.GetTotalCount(),
				"today_new": groups.GroupMediaService{}.GetTodayNewCount(),
			},
		},
	})
}

func (this Apiv1StaticsDao) getStaticsDown() map[string]int {
	down := make(map[string]int)
	// 统计今日新增种子量
	down["today_new"] = download.DownloadService{}.GetTodayNew()
	// 统计今日种子下载完成量
	down["today_new_down"] = download.DownloadService{}.GetTodayNewDownSuccess()
	// 统计总监控数量
	down["total"] = download.DownloadService{}.GetTotalCount()
	return down
}
func (this Apiv1StaticsDao) getStaticsGroup() map[string]int {
	down := make(map[string]int)
	// 统计总监控数量
	down["total"] = groups.GroupsService{}.GetTotalCount()
	return down
}
