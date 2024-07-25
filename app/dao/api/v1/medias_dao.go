package v1

import (
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/apiv1/medias"
	"XArr-Rss/app/model/dbmodel"
	match2 "XArr-Rss/app/model/services/match"
	"XArr-Rss/app/service/match"
	"XArr-Rss/app/service/matchv2"
	medias2 "XArr-Rss/app/service/medias"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	"strings"
)

type ApiV1Medias struct {
}

func (this ApiV1Medias) Get(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	ret := []dbmodel.Media{}
	mediaList := medias2.MediaService{}.GetMediaList()
	for _, v := range mediaList {
		ret = append(ret, v)
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return (ret[i].SonarrId) > (ret[j].SonarrId)
	})
	for k, v := range ret {
		if v.CnTitle == "" {
			ret[k].CnTitle = v.OriginalTitle
		}
	}

	c.JSON(200, gin.H{
		"data":    ret,
		"message": "success",
		"code":    0,
	})
}

// 测试正则匹配
func (this ApiV1Medias) TestReg(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	// install error handler
	defer func() {
		if e := recover(); e != nil {
			logsys.Error("err:%v", "TestReg", e)
		}
	}()
	var req medias.Apiv1MediasTestReq
	err = c.ShouldBind(&req)
	if err != nil {
		c.Status(200)
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求参数异常:" + err.Error(),
		})
		return
	}
	if req.SonarrId == 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请选择Sonarr媒体",
		})
		return
	}

	if req.MatchType != "auto" && req.Reg == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请填写正则匹配内容",
		})
		return
	}
	var ret []match2.MediaMatchInfo

	groupMedia := medias2.MediaService{}.GetMediaInfo(req.SonarrId)
	//var items []model.RssResultItem

	//titles := []string{}
	sourceItems := sources.GetSourceItems(req.UseSource, false)

	if len(sourceItems) == 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "数据源中暂无数据",
		})
		return
	}

	// 过滤字幕组
	items := match.FilterPushGroups(sourceItems, req.FilterPushGroup)

	// 获取匹配内容
	dbmGroupMeidia := &dbmodel.GroupMedia{
		SonarrId: req.SonarrId,
		Language: req.Language,
		Quality:  req.Quality,
		Regex: []dbmodel.GroupRegex{
			{
				MatchType: req.MatchType,
				Reg:       req.Reg,
				RegType:   req.RegType,
				Season:    req.Season,
				Offset:    req.Offset,
			},
		},
		UseSource:      nil,
		MediaInfo:      *groupMedia,
		EchoTitleAnime: req.EchoTitleAnime,
		EchoTitleTv:    req.EchoTitleTv,
	}

	groupCacheChannel := model.RssResult{}
	dbmGroupMeidia.MediaInfo.SearchTitle = strings.Join(matchv2.MatchV2{}.JoinSonarrTitle(dbmGroupMeidia.MediaInfo), "1666666666663")

	err = matchv2.MatchV2{}.Parse(dbmGroupMeidia, items, &groupCacheChannel)

	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	matchString := ""
	for _, v := range groupCacheChannel.Item {
		matchString += "新标题:<br>" + v.Title.Text + "<br>"
		matchString += "源标题:<br>" + v.OriginalTitle.Text + "<br>"
		if v.MinEpisode > 0 {
			if v.MaxEpisode != 0 && v.MinEpisode != v.MaxEpisode {
				matchString += "匹配最小集数:" + strconv.Itoa(v.MinEpisode)
				matchString += " 匹配最大集数:" + strconv.Itoa(v.MaxEpisode)
			} else {
				matchString += "匹配集数:" + strconv.Itoa(v.MinEpisode)
			}
		}

		if v.Season > 0 {
			matchString += " 匹配季数:" + strconv.Itoa(v.Season)
		}
		//if v.OtherInfo.Language != "" {
		//	matchString += " 匹配语言:" + (v.OtherInfo.Language)
		//}
		//if v.OtherInfo.QualityResolution != "" {
		//	matchString += " 匹配尺寸:" + (v.OtherInfo.QualityResolution)
		//}

		if req.Language == "-1" && v.OtherInfo.Language != "" {
			matchString += " 匹配语言:" + v.OtherInfo.Language
		}
		if req.Quality == "-1" && v.OtherInfo.QualityResolution != "" {
			matchString += " 匹配质量:" + v.OtherInfo.QualityResolution
		}

		matchString += "<hr>"
		ret = append(ret, match2.MediaMatchInfo{
			Title:         v.Title.Text,
			OriginalTitle: v.OriginalTitle.Text,
			OtherTitle:    v.OtherTitle.Text,
			MinEpisode:    v.MinEpisode,
			MaxEpisode:    v.MaxEpisode,
			OldMinEpisode: v.OldMinEpisode,
			OldMaxEpisode: v.OldMaxEpisode,
			Language:      v.OtherInfo.Language,
			Quality:       v.OtherInfo.QualityResolution,
			Season:        v.Season,
		})
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    ret,
		"message": "测试匹配完成,共匹配数量:" + strconv.Itoa(len(ret)) + "<br>" + matchString,
	})
}

type ApiV1MediaInfoRes struct {
	MediaInfo *dbmodel.Media                `json:"media_info,omitempty"`
	Season    map[int][]dbmodel.MediaSeason `json:"season,omitempty"`
	Episode   []dbmodel.MediaEpisodeList    `json:"episode,omitempty"`
}

func (this ApiV1Medias) Info(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	// 获取ID
	sonarrId := c.Query("id")
	if sonarrId == "" {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "请求参数丢失:id",
		})
		return
	}

	ret := &ApiV1MediaInfoRes{}
	ret.MediaInfo = medias2.MediaService{}.GetMediaInfo(helper.StrToInt(sonarrId))

	if ret.MediaInfo.CnTitle == "" {
		ret.MediaInfo.CnTitle = ret.MediaInfo.OriginalTitle
	}

	// 查询剧集信息
	ret.Episode = medias2.MediaEpisodeService{}.GetEpisodes(helper.StrToInt(sonarrId))

	c.JSON(200, gin.H{
		"data":    ret,
		"message": "success",
		"code":    0,
	})
}

func (this ApiV1Medias) Rename(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	// 获取ID
	sonarrId := c.PostForm("id")
	if sonarrId == "" {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "请求参数丢失:id",
		})
		return
	}

	cn_title := c.PostForm("cn_title")
	if cn_title == "" {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "请求参数丢失:cn_title",
		})
		return
	}
	titles := c.PostForm("titles")

	// 找到信息
	mediaInfo := medias2.MediaService{}.GetMediaInfo(helper.StrToInt(sonarrId))
	if mediaInfo == nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "没有找到媒体信息",
		})
		return
	}
	mediaInfo.CnTitle = cn_title
	mediaInfo.Titles = strings.Split(titles, "|")
	err = medias2.MediaService{}.Save(mediaInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "保存失败:" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data":    nil,
		"message": "success",
		"code":    0,
	})
}
