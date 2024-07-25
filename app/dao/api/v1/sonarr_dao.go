package v1

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/util/helper"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	"strings"
)

type Apiv1SonarrDao struct {
}

func (this Apiv1SonarrDao) Medias(c *gin.Context) {

	name := strings.Trim(c.Query("name"), " ")

	ret := []dbmodel.Media{}
	mediaList := medias.MediaService{}.GetMediaListForTitle(name)

	for _, v := range mediaList {
		// 处理image
		image := v.Image
		if image == "" {
			image = "/api/v1/sonarr/image/" + strconv.Itoa(v.SonarrId) + "/poster.jpg"
		}
		if v.CnTitle == "" {
			v.CnTitle = v.OriginalTitle
		}

		ret = append(ret, dbmodel.Media{
			CnTitle:       v.CnTitle,
			OriginalTitle: v.OriginalTitle,
			ImdbId:        v.ImdbId,
			TvdbId:        v.TvdbId,
			TmdbId:        v.TmdbId,
			Overview:      v.Overview,
			Image:         image,
			Year:          v.Year,
			SonarrId:      v.SonarrId,
			TitleSlug:     appconf.AppConf.Service.Sonarr.Host + "/series/" + v.TitleSlug,
		})
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return (ret[i].SonarrId) > (ret[j].SonarrId)
	})

	c.JSON(200, gin.H{
		"code":    0,
		"count":   len(ret),
		"data":    ret,
		"message": "success",
	})
}

func (this Apiv1SonarrDao) Image(c *gin.Context) {
	//http://127.0.0.1:20001/api/v3/MediaCover/1/poster.jpg?apikey=
	sonarrId := c.Param("sonarrId")
	image := c.Param("image")
	// 下载sonarr 图片
	uri := appconf.AppConf.Service.Sonarr.Host + "/api/v3/MediaCover/" + sonarrId + "/" + image + "?apikey=" + appconf.AppConf.Service.Sonarr.Apikey
	_, res, statusCode := helper.CurlHelper{}.GetUri(uri, nil, nil, false)
	if statusCode != 200 {
		//c.Status(statusCode)
		c.String(statusCode, "sonarr 服务异常")
		return
	}

	c.Data(200, "image/jpg", res)
}
