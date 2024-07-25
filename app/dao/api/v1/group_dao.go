package v1

import (
	"XArr-Rss/app/dao/monitor"
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/apiv1/group"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/sdk/sonarr"
	"XArr-Rss/app/service/groups"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/util/array"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ApiV1Groups struct {
}

type ApiV1GroupGetRes struct {
	Id                int                            `json:"id"`
	Name              string                         `json:"name"`
	Medias            map[string]*dbmodel.GroupMedia `json:"medias"`
	Url               string                         `json:"url"`
	TorznabUrl        string                         `json:"torznab_url"`
	Tags              string                         `json:"tags"`
	AutoInsertSonarr  int                            `json:"auto_insert_sonarr"`
	UpdatedAt         int64                          `json:"updated_at"`
	CreatedAt         int64                          `json:"created_at"`
	GroupTemplateId   int32                          `json:"group_template_id"`   // 创建时使用的模板ID
	GroupTemplateName string                         `json:"group_template_name"` // 创建时使用的模板名称
}

func (this ApiV1Groups) Get(c *gin.Context) {

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

	var ret []ApiV1GroupGetRes
	groupList := groups.GroupsService{}.GetGroupList(helper.StrToInt(groupId))
	templateIds := []int32{}
	for _, v := range groupList {
		templateIds = append(templateIds, v.GroupTemplateId)
	}

	allTemplate := groups.GroupTemplateService{}.FindALl(templateIds...)

	for _, v := range groupList {
		groupTempleteStr := ""
		if tpl, ok := allTemplate[int(v.GroupTemplateId)]; ok {
			groupTempleteStr = tpl.Name
		} else {
			v.GroupTemplateId = 0
		}
		if groupId != "" {
			if v.Id == (helper.StrToInt(groupId)) {
				ret = append(ret, ApiV1GroupGetRes{
					Id:        v.Id,
					Name:      v.Name,
					CreatedAt: v.CreatedAt.Unix(),
					UpdatedAt: v.UpdatedAt.Unix(),
					//Medias:           v.Medias,
					AutoInsertSonarr:  v.AutoInsertSonarr,
					Url:               appconf.AppConf.System.HttpAddr + "/rss/group/group_" + strconv.Itoa(int(v.Id)) + ".xml",
					TorznabUrl:        appconf.AppConf.System.HttpAddr + "/torznab/" + strconv.Itoa(int(v.Id)),
					Tags:              v.Tags,
					GroupTemplateId:   v.GroupTemplateId,
					GroupTemplateName: groupTempleteStr,
				})
			}
		} else {
			ret = append(ret, ApiV1GroupGetRes{
				Id:        v.Id,
				Name:      v.Name,
				CreatedAt: v.CreatedAt.Unix(),
				UpdatedAt: v.UpdatedAt.Unix(),
				//Medias:           v.Medias,
				AutoInsertSonarr:  v.AutoInsertSonarr,
				Url:               appconf.AppConf.System.HttpAddr + "/rss/group/group_" + strconv.Itoa(int(v.Id)) + ".xml",
				TorznabUrl:        appconf.AppConf.System.HttpAddr + "/torznab/" + strconv.Itoa(int(v.Id)),
				Tags:              v.Tags,
				GroupTemplateId:   v.GroupTemplateId,
				GroupTemplateName: groupTempleteStr,
			})
		}
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return (ret[i].Id) > (ret[j].Id)
	})

	c.JSON(200, gin.H{
		"code":    0,
		"data":    ret,
		"count":   len(ret),
		"message": "success",
	})
}

func (this ApiV1Groups) Add(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	groupsCount := groups.GroupsService{}.GetGroupCount()

	if groupsCount > 3 && !variable.ServerState.IsVip {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "非赞助会员用户只能创建3条分组信息哦!",
		})
		return
	}

	var req group.ApiV1GroupsAddReq
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}

	if len(req.Name) <= 3 {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "分组名称必须大于3个字符",
		})
		return
	}
	//conf_old.AppConfig.MaxIds.Groups++
	//nextId := conf_old.AppConfig.MaxIds.Groups
	item := &dbmodel.Group{
		Name:               req.Name,
		AutoInsertSonarr:   req.AutoInsertSonarr,
		LastInsertSonarrId: 0,
		Tags:               req.Tags,
		GroupTemplateId:    req.GroupTemplateId,
	}

	// 判断是否搞出最大的sonarrId
	if item.AutoInsertSonarr == dbmodel.AutoInsertSonarrNew {
		item.LastInsertSonarrId = medias.MediaService{}.GetMaxSonarrId()
	}
	if item.AutoInsertSonarr == dbmodel.AutoInsertSonarrTags && req.Tags == "" {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "创建分组选择导入Tags 那你得填入tags啊",
		})
		return
	}

	err = groups.GroupsService{}.CreateGroup(item)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "创建分组失败:" + err.Error(),
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

func (this ApiV1Groups) Delete(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	groupIdStr := strings.TrimLeft(c.Param("groupId"), "/")
	if groupIdStr == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "没有找到分组ID",
		})
		return
	}
	groupIds := strings.Split(groupIdStr, ",")
	for _, groupId := range groupIds {
		ok := groups.GroupsService{}.GroupExits(helper.StrToInt(groupId))
		if !ok {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"message": "没有找到分组信息:" + groupId,
			})
			return
		}
		err = groups.GroupsService{}.GroupDelete(helper.StrToInt(groupId))

		// 删除trans文件
		_ = os.Remove("./conf/trans/group_" + groupId + ".xml")
		if err != nil {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"message": "删除失败:" + err.Error(),
			})
			return
		}
		// 删除媒体数据
		err = groups.GroupMediaService{}.DeleteByGroupId(helper.StrToInt(groupId))
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "删除成功",
	})

}

func (this ApiV1Groups) Edit(c *gin.Context) {

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

	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}
	// 获取参数
	var req group.ApiV1GroupsEditReq
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求参数丢失",
		})
		return
	}
	if len(req.Name) <= 3 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "分组名必须大于3个字符",
		})
		return
	}

	// 判断是否切换了类型
	if groupInfo.AutoInsertSonarr != req.AutoInsertSonarr {
		// 重置新的数据
		groupInfo.LastInsertSonarrId = 0
	}

	// 判断是否搞出最大的sonarrId
	if groupInfo.AutoInsertSonarr != dbmodel.AutoInsertSonarrNew && req.AutoInsertSonarr == dbmodel.AutoInsertSonarrNew {
		mediaList := medias.MediaService{}.GetMediaList()

		for _, media := range mediaList {
			if (media.SonarrId) > (groupInfo.LastInsertSonarrId) {
				groupInfo.LastInsertSonarrId = (media.SonarrId)
			}
		}
	}
	groupInfo.Name = req.Name
	groupInfo.Tags = req.Tags
	groupInfo.AutoInsertSonarr = req.AutoInsertSonarr
	groupInfo.GroupTemplateId = req.GroupTemplateId

	err = groups.GroupsService{}.Save(groupInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "保存失败:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "保存成功",
	})
}

// 获取分组下面的媒体信息
func (this ApiV1Groups) GetGroupMedias(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	name := strings.Trim(c.Query("name"), " ")
	groupId := strings.TrimLeft(c.Param("groupId"), "/")
	if groupId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}

	if !(groups.GroupsService{}.GroupExits(helper.StrToInt(groupId))) {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}

	// 返回媒体数据

	var ret []Apiv1GetMediasRes
	groupMedias := groups.GroupMediaService{}.GetGroupMediasAndQuery(helper.StrToInt(groupId), name)

	for _, v := range groupMedias {
		var source []string
		// 查询数据源信息
		for _, sourceId := range v.UseSource {
			if sourceId == "-1" {
				source = append(source, "全部")
				break
			}

			sourceList := sources.SourcesService{}.GetSourcesList(false, sourceId)
			if len(sourceList) > 0 {
				source = append(source, sourceList[0].Name)
			} else {
				source = append(source, sourceId)

			}
		}

		quality := v.Quality
		if quality == "-1" {
			quality = "自动解析"
		}

		language := v.Language
		if language == "-1" {
			language = "自动解析"
		}

		item := Apiv1GetMediasRes{
			Id:              v.Id,
			Title:           v.MediaInfo.CnTitle,
			OriginalTitle:   v.MediaInfo.OriginalTitle,
			ImdbId:          v.MediaInfo.ImdbId,
			TvdbId:          v.MediaInfo.TvdbId,
			TmdbId:          v.MediaInfo.TmdbId,
			SonarrId:        v.SonarrId,
			TitleSlug:       appconf.AppConf.Service.Sonarr.Host + "/series/" + v.MediaInfo.TitleSlug,
			Regex:           v.Regex,
			FilterPushGroup: v.FilterPushGroup,
			Language:        language,
			Quality:         quality,
			UseSource:       v.UseSource,
			UseSourceText:   source,
		}
		ret = append(ret, item)
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    ret,
		"message": "获取成功",
	})
}

// 刷新分组匹配数据
func (this ApiV1Groups) _refresh(groupId string) error {
	groupInfo, err := groups.GroupsService{}.GroupInfo(helper.StrToInt(groupId))
	if err != nil {
		return errors.New("没有找到分组信息")
	}

	// 自动同步分组Sonarr剧集信息
	if groupInfo.AutoInsertSonarr != dbmodel.AutoInsertSonarrNone {
		// 同步Sonarr新数据
		monitor.GroupMonitor{}.SyncSonarrToGroup(*groupInfo)
	}
	err, _ = groups.GroupsService{}.SyncGroupItem(*groupInfo, false)
	return err
}

func (this ApiV1Groups) Refresh(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	groupId := c.Query("id")
	if groupId != "" {
		err = this._refresh(groupId)
		if err != nil {
			//continue
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"message": "刷新失败:" + err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"data":    "",
			"message": "刷新完成",
		})
	} else {

		res := monitor.GroupMonitor{}.SyncGroups(false)
		c.JSON(200, gin.H{
			"code":    0,
			"data":    res,
			"message": "刷新完成",
		})
	}

}

func (d ApiV1Groups) SetSonarrIndex(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	// 判断Sonarr是否已配置成功
	if appconf.AppConf.Service.Sonarr.Apikey == "" || appconf.AppConf.Service.Sonarr.Host == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "Sonarr配置信息未填写",
		})
		return
	}

	groupId := c.PostForm("groupId")
	indexType := c.PostForm("type")
	name := "XArr-Rss"
	groupInfo, err := groups.GroupsService{}.GroupInfo((helper.StrToInt(groupId)))
	if groupInfo != nil {
		name += " - " + groupInfo.Name
	} else {
		groupId = "all"
	}
	if indexType == "" {
		indexType = "rss"
	}

	name += " - " + indexType

	var req sonarr.SonarrApiIndexerPost
	if indexType == "rss" {
		req = sonarr.SonarrApiIndexerPost{
			EnableRss:               true,
			EnableAutomaticSearch:   false,
			EnableInteractiveSearch: false,
			SupportsRss:             true,
			SupportsSearch:          false,
			Protocol:                "torrent",
			Priority:                25,
			DownloadClientId:        0,
			Name:                    name,
			//Name: "TorrentRssIndexer",
			Fields: []struct {
				Name  string      `json:"name"`
				Value interface{} `json:"value,omitempty"`
			}{
				{
					Name:  "baseUrl",
					Value: appconf.AppConf.System.HttpAddr + "/rss/group/group_" + groupId + ".xml",
				},
				{
					Name: "cookie",
				},
				{
					Name:  "allowZeroSize",
					Value: false,
				},
				{
					Name:  "minimumSeeders",
					Value: 1,
				},
				{
					Name: "seedCriteria.seedRatio",
				},
				{
					Name: "seedCriteria.seedTime",
				},
				{
					Name: "seedCriteria.seasonPackSeedTime",
				},
			},
			ImplementationName: "Torrent RSS Feed",
			Implementation:     "TorrentRssIndexer",
			ConfigContract:     "TorrentRssIndexerSettings",
			InfoLink:           "https://wiki.servarr.com/sonarr/supported#torrentrssindexer",
			Tags:               []int{},
		}
	} else {
		if groupId == "all" {
			groupId = "xarr"
		}
		req = sonarr.SonarrApiIndexerPost{
			EnableRss:               true,
			EnableAutomaticSearch:   true,
			EnableInteractiveSearch: true,
			SupportsRss:             true,
			SupportsSearch:          true,
			Protocol:                "torrent",
			Priority:                25,
			DownloadClientId:        0,
			Name:                    name,
			//Name: "TorrentRssIndexer",
			Fields: []struct {
				Name  string      `json:"name"`
				Value interface{} `json:"value,omitempty"`
			}{
				{
					Name:  "baseUrl",
					Value: appconf.AppConf.System.HttpAddr + "/torznab/" + groupId,
				},
				{
					Name:  "apiPath",
					Value: "/api",
				},
				{
					Name:  "allowZeroSize",
					Value: true,
				},
				{
					Name:  "categories",
					Value: []int{5000},
				},
				{
					Name:  "animeCategories",
					Value: []int{5000},
				},
				{
					Name:  "minimumSeeders",
					Value: 1,
				},
				{
					Name:  "animeStandardFormatSearch",
					Value: false,
				},
				{
					Name: "seedCriteria.seedRatio",
				},
				{
					Name: "seedCriteria.seedTime",
				},
				{
					Name: "seedCriteria.seasonPackSeedTime",
				},
			},
			ImplementationName: "Torznab",
			Implementation:     "Torznab",
			ConfigContract:     "TorznabSettings",
			InfoLink:           "https://wiki.servarr.com/sonarr/supported#torznab",
			Tags:               []int{},
		}

	}

	//
	sendData, err := json.Marshal(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "数据出错啦:" + err.Error(),
		})
		return
	}

	// 获取剧集列表
	uri := strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/v3/indexer?apikey=" + appconf.AppConf.Service.Sonarr.Apikey
	err, result := helper.CurlHelper{}.PostString(uri, string(sendData), []helper.CurlHeader{{
		Name:  "Content-Type",
		Value: "application/json",
	}})
	if err != nil {
		logsys.Error("Sonarr信息查询失败:%s", "SonarrId", err.Error())

		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "Sonarr信息查询失败:" + err.Error(),
		})
		return
	}
	logsys.Debug("Sonarr返回:%s", "添加索引", string(result))
	if strings.Contains(string(result), "errorMessage") {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    string(sendData),
			"message": "添加失败,Sonarr返回结果:<br/>" + string(result),
		})
		return
	}

	if strings.Contains(string(result), "Should be unique") {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    string(sendData),
			"message": "添加失败,Sonarr已有相同名称索引!",
		})
		return
	}

	if strings.Contains(string(result), "XArr-Rss") {
		c.JSON(200, gin.H{
			"code":    0,
			"data":    nil,
			"message": "创建成功",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    string(sendData),
			"resp":    string(result),
			"message": "添加失败,可能是rss文件没有匹配内容",
		})
	}

}

func (this ApiV1Groups) GetPushGroup(c *gin.Context) {
	// 判断push group是否存在
	_, err := os.Stat("./conf/push_group.json")
	defaultPush := []string{}
	if err != nil {
		// 使用默认数据
		c.JSON(200, gin.H{
			"code":    0,
			"data":    defaultPush,
			"message": "获取成功",
		})
		return
	}
	// 解析
	data, err := os.ReadFile("./conf/push_group.json")
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    []string{},
			"message": "获取动漫组失败:" + err.Error(),
		})
		return
	}
	list := []string{}
	err = json.Unmarshal(data, &list)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    []string{},
			"message": "获取动漫组失败 配置文件异常:" + err.Error(),
		})
		return
	}

	for k, v := range list {
		list[k] = strings.Replace(v, "【", "[", -1)
		list[k] = strings.Replace(v, "】", "]", -1)
	}
	for _, v := range defaultPush {
		if !array.InArray(v, list) {
			list = append(list, v)
		}
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    list,
		"message": "获取成功",
	})

}
