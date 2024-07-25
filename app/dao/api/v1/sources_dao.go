package v1

import (
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/apiv1/source"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/source_item"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"errors"
	"github.com/gin-gonic/gin"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ApiV1Sources struct {
}

func (s ApiV1Sources) Get(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceId := strings.TrimLeft(c.Param("sourceId"), "/")
	name := c.Query("name")
	url := c.Query("url")

	var ret []dbmodel.Sources
	sourceList := sources.SourcesService{}.GetSourcesList(false)
	for _, v := range sourceList {
		if name != "" && !strings.Contains(v.Name, name) {
			// 判断搜索name
			continue
		}
		if url != "" && !strings.Contains(v.Url, url) {
			// 判断搜索name
			continue
		}

		if sourceId != "" {
			if strconv.Itoa(int(v.Id)) == sourceId {
				ret = append(ret, *v)
			}
		} else {
			ret = append(ret, *v)
		}
	}
	if len(ret) == 0 {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    ret,
			"count":   len(ret),
			"message": "没有查询到数据",
		})
		return
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return (ret[i].Id) > (ret[j].Id)
	})

	for k, v := range ret {
		ret[k].RefreshTime = v.RefreshTime / 60
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    ret,
		"count":   len(ret),
		"message": "success",
	})
}

func (s ApiV1Sources) GetMedias(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceId := strings.TrimLeft(c.Param("sourceId"), "/")

	var ret *dbmodel.Sources
	ret, err = sources.SourcesService{}.GetSourcesInfo(sourceId)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "没有查询到数据源",
		})
		return
	}

	list := source_item.SourceItemService{}.GetSourcesItemList(int(ret.Id), 100)
	result := source_item.SourceItemService{}.FormatToXml(list)

	c.JSON(200, gin.H{
		"code":    0,
		"data":    result,
		"message": "success",
	})
}

func (s ApiV1Sources) Add(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceLen := sources.SourcesService{}.GetSourcesCount()

	if sourceLen > 5 && !variable.ServerState.IsVip {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "非赞助会员用户只能创建5条数据源哦!",
		})
		return
	}
	var req source.Apiv1SourceAddReq
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "数据解析错误:" + err.Error(),
		})
		return
	}

	if strings.Index(req.Url, "http") != 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "订阅地址填写错误,请以http[s]://开头",
		})
		return
	}

	if len(req.Name) < 3 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "分组名必须大于2个字符",
		})
		return
	}
	if req.RefreshTime < 5 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "刷新间隔时间不得低于5分钟",
		})
		return
	}
	proxy := ""
	if req.UseProxy == 1 {
		proxy = options.GetOption(options.OptionsGlobalProxy)
		if proxy == "" {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"success": false,
				"message": "兄弟,你把 [系统设置-代理] 填写保存后 再来设置开启代理吧",
			})
			return
		}
	}
	if req.ProxySiteType != dbmodel.ProxySiteTypeDefaultRss {
		if req.ProxySiteApiKey == "" {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"success": false,
				"message": "代理请求必须包含apikey",
			})
			return
		}
	}
	uri := req.Url
	if req.ProxySiteType != dbmodel.ProxySiteTypeDefaultRss {
		uri += "?t=caps&apikey=" + req.ProxySiteApiKey
	}
	//if req.ProxySiteType != dbmodel.ProxySiteTypeJacket {
	//	uri += "/api?t=caps&apikey=" + req.ProxySiteApiKey
	//}
	//if req.ProxySiteType != dbmodel.ProxySiteTypeProwlarr {
	//	uri += "?t=caps&apikey=" + req.ProxySiteApiKey
	//}

	err = checkSourceUrl(uri, proxy, req.Name, req.AutoSearch)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "订阅数据源获取错误,请检查订阅地址,错误信息:" + err.Error(),
		})
		return
	}

	err, errRegStr := checkSourceReg(req.Regex.MustHave, req.Regex.MustDotHave)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "错误的正则规则：" + errRegStr + " 错误：" + err.Error(),
		})
		return
	}
	//>> 新增Source
	data := &dbmodel.Sources{
		Name: req.Name,
		//Site:            req.Site,
		Url:             req.Url,
		RefreshTime:     req.RefreshTime * 60,
		LastRefreshTime: 0,
		MaxReadCount:    req.MaxReadCount,
		Regex:           req.Regex,
		CacheDay:        req.CacheDay,
		UseProxy:        req.UseProxy,
		ProxySiteType:   req.ProxySiteType,
		ProxySiteApiKey: req.ProxySiteApiKey,
		MaxCount:        req.MaxCount,
		Status:          req.Status,
		AutoSearch:      req.AutoSearch,
		DownloadPasskey: req.DownloadPasskey,
	}
	err = sources.SourcesService{}.CreateSource(data)
	if err != nil && strings.Index(err.Error(), "UNIQUE constraint failed") > 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "新增数据源错误：数据源名称重复",
		})
		return
	}
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "新增数据源错误：" + err.Error(),
		})
		return
	}
	// 提交刷新
	go sources.SourcesService{}.SyncSource(data, true)

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"success": true,
		"message": "保存成功",
	})
}

func checkSourceReg(mustHas string, mustDotHave string) (error, string) {
	if mustHas != "" {
		_, err := regexp.Compile(mustHas)
		if err != nil {
			return err, mustHas
		}
	}

	if mustDotHave != "" {
		_, err := regexp.Compile(mustDotHave)
		if err != nil {
			return err, mustDotHave
		}
	}

	return nil, ""
}

func (s ApiV1Sources) Delete(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceId := strings.TrimLeft(c.Param("sourceId"), "/")
	if sourceId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return

	}
	sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(sourceId)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,操作个der",
		})
		return
	}
	err = source_item.SourceItemService{}.DeleteBySourceId(sourceInfo.Id)
	if err != nil {
		logsys.Error("删除数据失败:%s", "删除数据源", err.Error())
	}
	err = sources.SourcesService{}.Delete(sourceInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "删除数据源失败:" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "删除成功",
	})
}

func (s ApiV1Sources) Edit(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceId := strings.TrimLeft(c.Param("sourceId"), "/")
	if sourceId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return

	}

	sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(sourceId)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return
	}
	// 获取参数
	var req source.Apiv1SourceEditReq
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求参数丢失:" + err.Error(),
		})
		return
	}
	if len(req.Name) < 3 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "分组名必须大于2个字符",
		})
		return
	}

	if req.RefreshTime < 5 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "刷新间隔时间不得低于5分钟",
		})
		return
	}

	if strings.Index(req.Url, "http") != 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "订阅地址填写错误,请以http[s]://开头",
		})
		return
	}
	proxy := ""
	if req.UseProxy == 1 {
		proxy = options.GetOption(options.OptionsGlobalProxy)
		if proxy == "" {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"success": false,
				"message": "兄弟,你把 [系统设置-代理] 填写保存后 再来设置开启代理吧",
			})
			return
		}
	}

	if req.ProxySiteType != dbmodel.ProxySiteTypeDefaultRss {
		if req.ProxySiteApiKey == "" {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"success": false,
				"message": "代理请求必须包含apikey",
			})
			return
		}
	}
	uri := req.Url
	if req.ProxySiteType != dbmodel.ProxySiteTypeDefaultRss {
		uri += "?t=caps&apikey=" + req.ProxySiteApiKey
	}
	err = checkSourceUrl(uri, proxy, req.Name, req.AutoSearch)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "订阅数据源获取错误,请检查订阅地址:" + err.Error(),
		})
		return
	}

	err, errRegStr := checkSourceReg(req.Regex.MustHave, req.Regex.MustDotHave)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"success": false,
			"message": "错误的正则规则：" + errRegStr + " 错误：" + err.Error(),
		})
		return
	}
	sourceInfo.Name = req.Name
	//sourceInfo.Site = req.Site
	sourceInfo.UseProxy = req.UseProxy
	sourceInfo.Url = req.Url
	sourceInfo.CacheDay = req.CacheDay
	sourceInfo.RefreshTime = req.RefreshTime * 60
	sourceInfo.MaxReadCount = req.MaxReadCount
	sourceInfo.Regex = req.Regex
	sourceInfo.ProxySiteType = req.ProxySiteType
	sourceInfo.ProxySiteApiKey = req.ProxySiteApiKey
	sourceInfo.MaxCount = req.MaxCount
	sourceInfo.Status = req.Status
	sourceInfo.AutoSearch = req.AutoSearch
	sourceInfo.DownloadPasskey = req.DownloadPasskey

	err = sources.SourcesService{}.Save(sourceInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "保存失败:" + err.Error(),
		})
		return
	}

	go sources.SourcesService{}.SyncSource(sourceInfo, true)

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "保存成功",
	})
}

func (s ApiV1Sources) Refresh(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceId := strings.TrimLeft(c.Param("sourceId"), "/")
	if sourceId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return

	}
	sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(sourceId)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return
	}

	res := sources.SourcesService{}.SyncSource(sourceInfo, true)
	if res == nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "刷新失败",
		})

		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    res,
		"message": "刷新完成",
	})

}

func (s ApiV1Sources) RefreshParse(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceId := strings.TrimLeft(c.Param("sourceId"), "/")
	if sourceId == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return

	}
	sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(sourceId)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "什么瘠薄都没有,编辑个der",
		})
		return
	}

	sources.ReParseSourceTitle(strconv.Itoa(sourceInfo.Id), true)

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "刷新完成",
	})

}

// 刷新数据源匹配
func (this ApiV1Sources) ReParseInfo(c *gin.Context) {
	sourceId := c.Query("source_id")
	sources.ReParseTitles(true, sourceId)

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "刷新完成",
	})
}

func (s ApiV1Sources) GetMediasJson(c *gin.Context) {
	name := c.Query("name")
	page := helper.StrToInt(c.Query("page"))
	limit := helper.StrToInt(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	sourceId := strings.TrimLeft(c.Param("sourceId"), "/")

	sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(sourceId)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "没有查询到数据源",
		})
		return
	}

	// 查询下面的媒体

	list := source_item.SourceItemService{}.GetSourcesItemListPage(int(sourceInfo.Id), name, page, limit)
	//result := source_item.SourceItemService{}.FormatToXml(list)

	type retData struct {
		Title        string `json:"title,omitempty"`
		Link         string `json:"link,omitempty"`
		PubDate      string `json:"pub_date,omitempty"`
		EnclosureUrl string `json:"enclosure_url,omitempty"`
		IsParse      bool   `json:"is_parse,omitempty"`
	}
	var retList []retData

	for _, v := range list {
		retList = append(retList, retData{
			Title:        v.Title,
			Link:         v.Link,
			PubDate:      v.PubDate,
			EnclosureUrl: v.Enclosure.Url,
			IsParse:      v.ParseInfo != nil && v.ParseInfo.MinEpisode > 0,
		})
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    retList,
		"count":   source_item.SourceItemService{}.GetSourcesItemListCount(int(sourceInfo.Id), name),
		"message": "success",
	})
}

func checkSourceUrl(url, proxy, name, autoSearch string) error {
	rssResult, err := sources.GetRssUrlResult(url, proxy, name, autoSearch, "")
	if err != nil {
		return logsys.Error("%s", "数据源", err.Error())
	}
	if rssResult == nil {
		return logsys.Error("获取数据源出错,请检查连接是否正常", "数据源")
	}

	// 判断是否为会员
	if variable.ServerState.IsVip == false {
		// 验证url地址是否能够匹配
		for _, v := range rssResult.Channel.Item {
			//v.Link
			// 取出hash
			compile, err := regexp.Compile(`[retData-z0-9A-Z]{40}`)
			if err != nil {
				return errors.New("匹配规则出错")
			}
			if compile.MatchString(v.Link) {
				continue
			} else {
				// hash未找到
				return errors.New("抱歉,当前为普通用户,不支持连接中没有种子hash值的链接,赞助会员支持磁力和任何种子")
			}

		}
	}

	return nil
}

func (s ApiV1Sources) ImportJacketIndex(c *gin.Context) {

}
