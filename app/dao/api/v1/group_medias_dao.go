package v1

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/apiv1/group"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/groups"
	"XArr-Rss/app/service/match"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/util/array"
	"XArr-Rss/util/helper"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
)

type Apiv1GroupMediasDao struct {
}

func (this Apiv1GroupMediasDao) AutoGenReg(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	if !variable.ServerState.IsVip {
		c.JSON(200, gin.H{
			"code":    564,
			"data":    nil,
			"message": "非赞助会员用户无法使用当前功能",
		})
		return
	}

	regStr := c.PostForm("req_str")
	if regStr == "" {
		c.JSON(200, gin.H{
			"code":    564,
			"data":    nil,
			"message": "请输入需要匹配的正则规则",
		})
		return
	}
	sonarrId := c.PostForm("sonarr_id")

	mediaInfo := medias.MediaService{}.GetMediaInfo(helper.StrToInt(sonarrId))
	if mediaInfo == nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "未能找到Sonarr媒体信息",
		})
		return
	}
	// 提取中文标题
	sonarrTitle := mediaInfo.CnTitle
	if mediaInfo.OriginalTitle != "" {
		sonarrTitle += "-|-" + mediaInfo.OriginalTitle
	}
	sonarrTitle += "-|-" + strings.Join(mediaInfo.Titles, "-|-")

	if sonarrTitle == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "未能找到Sonarr媒体标题信息",
		})
		return
	}

	// 把标题都替换为占位符
	///////////////////////////////
	var zwf = []string{}
	sonarrTitleArr := strings.Split(sonarrTitle, "-|-")
	var zwfArr = strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	for _, v := range sonarrTitleArr {
		v = helper.ReplaceRegString(v)
		if v == "" {
			continue
		}

		// 模糊匹配
		hasMatch, matchString := helper.MatchString(strings.ToLower(v), strings.ToLower(regStr))
		if hasMatch {
			zwf = append(zwf, matchString)
			regStr = strings.Replace(regStr, matchString, strings.Repeat(zwfArr[len(zwf)-1], 10), -1)
			continue
		}
	}
	/////////////////////////////

	err, retRegexStr, _, _, _ := match.ParseMediaTitleEpisode(regStr)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    564,
			"data":    nil,
			"message": "匹配错误:" + err.Error(),
		})
		return
	}
	if retRegexStr == "" {
		c.JSON(200, gin.H{
			"code":    564,
			"data":    nil,
			"message": "匹配错误:" + retRegexStr,
		})
		return
	}
	for i, v := range zwfArr {
		if i > len(zwf)-1 {
			continue
		}

		retRegexStr = strings.Replace(retRegexStr, strings.Repeat(v, 10), zwf[i], -1)
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    retRegexStr,
		"message": "匹配成功",
	})
	return
}

// 批量删除
func (this Apiv1GroupMediasDao) BatchRemove(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	ids := c.Query("ids")
	idsArr := strings.Split(ids, ",")
	if len(idsArr) == 0 {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "请选择要删除的媒体数据",
		})
		return
	}

	groupId := strings.TrimLeft(c.Param("groupId"), "/")
	if groupId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return

	}

	groupInfo, err := groups.GroupsService{}.GroupInfo((helper.StrToInt(groupId)))
	if groupInfo == nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}

	// 搜索分组媒体数据
	db.MainDb.Model(dbmodel.GroupMedia{}).Where("id in ?", idsArr).Where("group_id = ?", groupId).Delete(&dbmodel.GroupMedia{})

	for _, id := range idsArr {
		// 删除对应媒体数据
		os.RemoveAll(appconf.AppConf.ConfDir + "/trans/group_" + groupId + "/" + (id))
		os.Remove(appconf.AppConf.ConfDir + "/trans/group_" + groupId + "/" + (id) + ".xml")
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "删除成功",
	})

}

// 增加媒体
func (this Apiv1GroupMediasDao) AddMedias(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	groupId := strings.TrimLeft(c.Param("groupId"), "/")
	if groupId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return

	}
	groupInfo, err := groups.GroupsService{}.GroupInfo(helper.StrToInt(groupId))
	if groupInfo == nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "没有找到分组信息",
		})
		return
	}

	var req group.Apiv1GroupMediaAdd
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}

	if len(req.Regex) == 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请添加规则",
		})
		return
	}

	if len(req.UseSource) == 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请选择数据源",
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
	if !variable.ServerState.IsVip && req.Language == "-1" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用自动搜索语言",
		})
		return
	}
	if !variable.ServerState.IsVip && req.Quality == "-1" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用自动搜索质量",
		})
		return
	}
	for _, v := range req.Regex {
		if v.MatchType == "auto" && !variable.ServerState.IsVip {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"message": "赞助会员可使用全自动匹配功能",
			})
			return
		}
	}

	if !variable.ServerState.IsVip && req.EchoTitleAnime != "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用分组内控制输出标题格式",
		})
		return
	}
	if !variable.ServerState.IsVip && req.EchoTitleTv != "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用分组内控制输出标题格式",
		})
		return
	}
	if array.InArray("-1", req.UseSource) {
		req.UseSource = []string{"-1"}
	}

	if err := (groups.GroupMediaService{}.AddGroupMedias(helper.StrToInt(groupId), req)); err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "保存错误:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "新增成功",
	})
	return
}

type Apiv1GetMediasRes struct {
	Id              int                  `json:"id"`                // id
	Title           string               `json:"title"`             // 资源名称 - 海贼王
	OriginalTitle   string               `json:"name"`              // 资源名称 - Sonarr标题
	ImdbId          string               `json:"imdbid"`            // 媒体ID = imdbid
	TvdbId          int                  `json:"tvdbid"`            // 媒体ID = imdbid
	TmdbId          string               `json:"tmdbId"`            // 媒体ID = imdbid
	SonarrId        int                  `json:"sonarr_id"`         // SonarrId
	TitleSlug       string               `json:"title_slug"`        // 用于跳转到Sonarr
	Regex           []dbmodel.GroupRegex `json:"regex"`             // 正则匹配规则
	Language        string               `json:"language"`          //  语言
	Quality         string               `json:"quality"`           // 质量
	UseSource       []string             `json:"use_source"`        // 使用那些数据源
	UseSourceText   []string             `json:"use_source_text"`   // 使用那些数据源
	FilterPushGroup []string             `json:"filter_push_group"` // 过滤发布组
}

func (this Apiv1GroupMediasDao) GetMediasInfo(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	groupId := strings.TrimLeft(c.Param("groupId"), "/")
	mediaId := strings.TrimLeft(c.Param("mediaId"), "/")
	if groupId == "" || mediaId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}
	groupInfo, err := groups.GroupsService{}.GroupInfo((helper.StrToInt(groupId)))

	if groupInfo == nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}

	groupMedia, err := groups.GroupMediaService{}.GetGroupMediaInfo(helper.StrToInt(mediaId))

	// 找到资源
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "没有找到媒体信息",
		})
		return
	}

	// 返回媒体数据

	c.JSON(200, gin.H{
		"code":    0,
		"data":    groupMedia,
		"message": "获取成功",
	})
}

func (this Apiv1GroupMediasDao) EditMedias(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	groupId := strings.TrimLeft(c.Param("groupId"), "/")
	mediaId := strings.TrimLeft(c.Param("mediaId"), "/")
	if groupId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return

	}
	groupInfo, err := groups.GroupsService{}.GroupInfo(helper.StrToInt(groupId))

	if groupInfo == nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return
	}

	groupMedia, err := groups.GroupMediaService{}.GetGroupMediaInfo(helper.StrToInt(mediaId))

	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return
	}

	var req group.Apiv1GroupMediaAdd
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}
	mediaInfo := medias.MediaService{}.GetMediaInfo(req.SonarrId)
	if mediaInfo == nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "没有找到媒体信息",
		})
		return
	}
	if array.InArray("-1", req.UseSource) {
		req.UseSource = []string{"-1"}
	}

	if !variable.ServerState.IsVip && req.Language == "-1" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用自动搜索语言",
		})
		return
	}
	if !variable.ServerState.IsVip && req.Quality == "-1" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用自动搜索质量",
		})
		return
	}
	for _, v := range req.Regex {
		if v.MatchType == "auto" && !variable.ServerState.IsVip {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"message": "赞助会员可使用全自动匹配功能",
			})
			return
		}
	}

	if !variable.ServerState.IsVip && req.EchoTitleAnime != "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用分组内控制输出标题格式",
		})
		return
	}
	if !variable.ServerState.IsVip && req.EchoTitleTv != "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "赞助会员可使用分组内控制输出标题格式",
		})
		return
	}
	groupMedia.Regex = req.Regex
	groupMedia.SonarrId = req.SonarrId
	groupMedia.UseSource = req.UseSource
	groupMedia.Quality = req.Quality
	groupMedia.Language = req.Language
	groupMedia.FilterPushGroup = req.FilterPushGroup
	groupMedia.EchoTitleAnime = req.EchoTitleAnime
	groupMedia.EchoTitleTv = req.EchoTitleTv

	if err := (groups.GroupMediaService{}.Save(groupMedia)); err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "保存错误:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "保存成功",
	})
}

func (this Apiv1GroupMediasDao) DeleteMedias(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	groupId := strings.TrimLeft(c.Param("groupId"), "/")
	mediaId := strings.TrimLeft(c.Param("mediaId"), "/")
	if groupId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return

	}
	groupInfo, err := groups.GroupsService{}.GroupInfo(helper.StrToInt(groupId))

	if groupInfo == nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}

	groupMedia, err := groups.GroupMediaService{}.GetGroupMediaInfo(helper.StrToInt(mediaId))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}

	err = groups.GroupMediaService{}.DeleteMedia(groupMedia)
	if err == nil {
		// 删除对应媒体数据
		os.RemoveAll(appconf.AppConf.ConfDir + "/trans/group_" + groupId + "/" + strconv.Itoa(groupMedia.Id))
		os.Remove(appconf.AppConf.ConfDir + "/trans/group_" + groupId + "/" + strconv.Itoa(groupMedia.Id) + ".xml")
		c.JSON(200, gin.H{
			"code":    0,
			"data":    nil,
			"message": "删除成功",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    500,
		"data":    nil,
		"message": "删除错误:" + err.Error(),
	})
}
