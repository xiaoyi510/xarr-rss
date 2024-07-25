package v1

import (
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/apiv1/toznab"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/groups"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

type Apiv1ToznabDao struct {
}

func (this Apiv1ToznabDao) Api(c *gin.Context) {
	//if !variable.ServerState.IsVip {
	//	logsys.Error("Toznab功能只有赞助会员可以使用", "Toznab")
	//	return
	//}
	log.Println("ttt")
	c.Request.URL.RawQuery = strings.ReplaceAll(c.Request.URL.RawQuery, "+", "%2b")

	groupId := c.Param("groupId")
	if groupId == "xarr" {
		groupId = "all"
	}
	t := c.Query("t")
	logsys.Debug("索引请求urlraw:%s", "Toznab搜索", c.Request.RequestURI)

	switch t {
	case "caps":
		this.Caps(c)
	case "search":
		this.Search(c, groupId, t)
	case "tvsearch":
		this.Search(c, groupId, t)
	default:
		logsys.Error("请求方式错误:"+t, "Toznab")
		c.String(200, "服务运行中")
	}
}

// 搜索内容
func (this Apiv1ToznabDao) Search(c *gin.Context, groupId, t string) {
	query := c.Query("q")
	tvdbid := c.Query("tvdbid")
	ep := c.Query("ep")
	logsys.Debug("索引进行搜索 q:%s", "Toznab搜索", query)

	res := toznab.Apiv1ToznabSearchRes{
		Text:    "",
		Version: "1.0",
		Atom:    "http://www.w3.org/2005/Atom",
		Torznab: "http://torznab.com/schemas/2015/feed",
		Channel: toznab.Apiv1ToznabSearchChannel{
			Text: "",
			Link: toznab.Apiv1ToznabSearchChannelLink{
				Text: "",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Title: toznab.Apiv1ToznabSearchChannelTitle{
				Text: "XArr-Rss",
			},
			Item: ApiTorznabNewDao{}.GetDefaultItem(),
		},
	}

	this.getGroupMediaXml(groupId, query)

	// 查询group信息
	_, err := groups.GroupsService{}.GroupInfo(helper.StrToInt(groupId))
	if err != nil {
		logsys.Error("没有找到分组信息:%s", "toznab", err.Error())
		c.XML(200, res)
		return
	}

	var groupXml model.RssRoot

	// 提取groupXml
	_, err = os.Stat("./conf/trans/group_" + groupId + ".xml")
	if err != nil {
		logsys.Error("获取./conf/trans/group_"+groupId+".xml 异常:%s", "torznab", err.Error())
		c.XML(200, res)
		return
	}
	// 解析groupxml
	groupXmlData, err := os.ReadFile("./conf/trans/group_" + groupId + ".xml")
	if err != nil {
		logsys.Error("读取./conf/trans/group_"+groupId+".xml 异常:%s", "torznab", err.Error())
		c.XML(200, res)
		return
	}
	err = xml.Unmarshal(groupXmlData, &groupXml)
	if err != nil {
		logsys.Error("解析./conf/trans/group_"+groupId+".xml 异常:%s", "torznab", err.Error())
		c.XML(200, res)
		return
	}

	for _, v := range groupXml.Channel.Item {
		if v.Title.Text == "XArr-Rss 默认占位使用" {
			continue
		}
		if query != "" {
			qArr := strings.Split(query, "+")
			if len(qArr) == 2 {
				// 搜索标题
				qItem := qArr[0]
				qItem = strings.ToLower(helper.ReplaceRegString(qItem))
				if !strings.Contains(strings.ToLower(helper.ReplaceRegString(v.Title.Text)), qItem) && !strings.Contains(strings.ToLower(helper.ReplaceRegString(v.OriginalTitle.Text)), qItem) && !strings.Contains(strings.ToLower(helper.ReplaceRegString(v.OtherTitle.Text)), qItem) {
					continue
				}
				if v.OldMinEpisode <= helper.StrToInt(qArr[1]) && helper.StrToInt(qArr[1]) <= v.OldMaxEpisode {
					// 允许搜索
				} else if v.MinEpisode <= helper.StrToInt(qArr[1]) && helper.StrToInt(qArr[1]) <= v.MaxEpisode {
					// 允许搜索
				} else {
					// 不允许搜索
					continue
				}

			}

		}
		//if season != "" {
		//	if v.Season != helper.StrToInt(season) {
		//		continue
		//	}
		//}
		if tvdbid != "" {
			if v.OthderId.TvdbId != tvdbid {
				continue
			}
		}
		if ep != "" && v.MinEpisode != -1 {
			if v.MaxEpisode != -1 {
				// 是连续的
				if !(v.MinEpisode <= helper.StrToInt(ep) && helper.StrToInt(ep) <= v.MaxEpisode) {
					continue
				}
			} else {
				if v.MinEpisode != helper.StrToInt(ep) {
					continue
				}
			}
		}

		res.Channel.Item = append(res.Channel.Item, toznab.Apiv1ToznabSearchChannelItem{
			Text: "",
			Title: toznab.Apiv1ToznabSearchChannelItemTitle{
				Text: v.Title.Text,
			},
			OriginalTitle: toznab.Apiv1ToznabSearchChannelItemTitle{
				Text: v.OriginalTitle.Text,
			},
			OtherTitle: toznab.Apiv1ToznabSearchChannelItemTitle{
				Text: v.OtherTitle.Text,
			},
			Guid: toznab.Apiv1ToznabSearchChannelItemGuid{
				Text: v.Guid.Text,
			},
			XArrRssindexer: toznab.Apiv1ToznabSearchChannelItemXArrRssindexer{
				Text: v.XArrRssIndexer.Text,
				ID:   v.XArrRssIndexer.ID,
			},
			Comments: toznab.Apiv1ToznabSearchChannelItemComments{
				Text: v.Link,
			},
			PubDate: toznab.Apiv1ToznabSearchChannelItemPubDate{
				Text: v.PubDate,
			},
			Size: toznab.Apiv1ToznabSearchChannelItemSize{
				Text: v.Enclosure.Length,
			},
			Link: toznab.Apiv1ToznabSearchChannelItemLink{
				Text: v.Enclosure.Url,
			},
			Category: toznab.Apiv1ToznabSearchChannelItemCategory{
				Text: "5000",
			},
			Enclosure: toznab.Apiv1ToznabSearchChannelItemEnclosure{
				Text:   "",
				URL:    v.Enclosure.Url,
				Length: v.Enclosure.Length,
				Type:   v.Enclosure.Type,
			},
			Attr: []toznab.Apiv1ToznabSearchChannelItemAttr{
				{
					Text:  "",
					Name:  "category",
					Value: "5000",
				},
				{
					Text:  "",
					Name:  "tag",
					Value: "freeleech",
				},
				{
					Text:  "",
					Name:  "rageid",
					Value: "0",
				},
				{
					Text:  "",
					Name:  "tvdbid",
					Value: v.OthderId.TvdbId,
				},
				{
					Text:  "",
					Name:  "imdb",
					Value: v.OthderId.ImdbId,
				},
				{
					Text:  "",
					Name:  "tmdbid",
					Value: v.OthderId.TmdbId,
				},
				{
					Text:  "",
					Name:  "traktid",
					Value: "0",
				},
				{
					Text:  "",
					Name:  "seeders",
					Value: "20",
				},
				{
					Text:  "",
					Name:  "grabs",
					Value: "21",
				},
				{
					Text:  "",
					Name:  "peers",
					Value: "22",
				},
				{
					Text:  "",
					Name:  "downloadvolumefactor",
					Value: "0",
				},
				{
					Text:  "",
					Name:  "uploadvolumefactor",
					Value: "1",
				},
			},
		})
	}

	c.XML(200, res)

}

func (this Apiv1ToznabDao) Caps(c *gin.Context) {
	res := toznab.Apiv1ToznabCapsRes{
		XMLName: xml.Name{},
		Text:    "",
		Server: toznab.Apiv1ToznabCapsServer{
			Text:  "",
			Title: "XArr-Rss",
		},
		Limits: toznab.Apiv1ToznabCapsLimits{
			Text:    "",
			Max:     "100",
			Default: "100",
		},
		Searching: toznab.Apiv1ToznabCapsSearching{
			Text: "",
			Search: toznab.Apiv1ToznabCapsSearchingSearch{
				Text:            "",
				Available:       "no",
				SupportedParams: "q,ep,season,episode,tvdbid,imdbid,rid,title",
			},
			TvSearch: toznab.Apiv1ToznabCapsSearchingTvSearch{
				Text:            "",
				Available:       "yes",
				SupportedParams: "q,ep,season,episode,tvdbid,imdbid,rid,title",
			},
			MovieSearch: toznab.Apiv1ToznabCapsSearchingMovieSearch{
				Text:            "",
				Available:       "no",
				SupportedParams: "q",
			},
			MusicSearch: toznab.Apiv1ToznabCapsSearchingMusicSearch{
				//Text:            "",
				Available:       "no",
				SupportedParams: "q",
			},
			AudioSearch: toznab.Apiv1ToznabCapsSearchingAudioSearch{
				Text:            "",
				Available:       "no",
				SupportedParams: "q",
			},
			BookSearch: toznab.Apiv1ToznabCapsSearchingBookSearch{
				Text:            "",
				Available:       "no",
				SupportedParams: "q",
			},
		},
		Categories: toznab.Apiv1ToznabCapsCategories{
			Text: "",
			Category: []toznab.Apiv1ToznabCapsCategoriesCategory{
				{
					Text: "",
					ID:   "5000",
					Name: "TV",
					Subcat: []toznab.Apiv1ToznabSubCat{
						{
							ID:   "5070",
							Name: "TV/Anime",
						},
					},
				}, {
					Text: "",
					ID:   "100001",
					Name: "Anime",
				},
			},
		},
		Tags: toznab.Apiv1ToznabCapsTags{},
	}
	c.XML(200, res)
}

func (this Apiv1ToznabDao) proxyTorznab(groupInfo *dbmodel.Group, cat, q string) []model.RssResultItem {
	// 1. 查询绑定的数据源信息
	// http://10.10.102.105:9117/api/v2.0/indexers/mikan/results/torznab/api?t=search&cat=5000,100001&extended=1&apikey=lvvplrg6dry8ongl5g7a75uh8x190uwa&offset=0&limit=100&q=A%20Will%20Eternal+55
	// http://192.168.10.10:9696/16/api?t=search&cat=5000,100001&extended=1&apikey=lvvplrg6dry8ongl5g7a75uh8x190uwa&offset=0&limit=100&q=A%20Will%20Eternal+55

	//for _, bindSourceId := range groupInfo.BindSourceIds {
	//	sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(helper.IntToStr(bindSourceId))
	//	if err != nil || sourceInfo == nil {
	//		return nil
	//	}
	//	uri := ""
	//	switch sourceInfo.ProxySiteType {
	//	case dbmodel.ProxySiteTypeJacket:
	//		// 使用jacket 请求
	//		uri = sourceInfo.ProxySite + "/api?t=search&cat=" + cat + "&extended=1&apikey=" + sourceInfo.ProxySiteApiKey + "&offset=0&limit=100&q=" + q
	//		curl := helper.GetCurlHttpHelperDefault()
	//		proxy := ""
	//		if sourceInfo.UseProxy == 1 {
	//			proxy = options.GetOption(options.OptionsGlobalProxy)
	//		}
	//		_, result := curl.GetProxyResult(uri, proxy)
	//		logsys.Debug("代理搜索返回结果:%s", "代理网站", string(result))
	//	}
	//	if uri == "" {
	//		continue
	//	}
	//}
	return nil
}

func (this Apiv1ToznabDao) getGroupSourceItems(groupInfo *dbmodel.Group) (ret []model.RssResultItem) {
	//wg := golimit.GoLimit{}
	//wg.SetMax(3)
	//for _, bindSourceId := range groupInfo.BindSourceIds {
	//	sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(helper.IntToStr(bindSourceId))
	//	if err != nil || sourceInfo == nil {
	//		return nil
	//	}
	//	wg.Add()
	//	go func(sourceInfoParam *dbmodel.Sources) {
	//		defer wg.Done()
	//		//uri := ""
	//		switch sourceInfoParam.ProxySiteType {
	//		case dbmodel.ProxySiteTypeJacket:
	//			// 代理jacket
	//
	//		case dbmodel.ProxySiteTypeProwlarr:
	//			// 代理prowlarr
	//
	//			break
	//		case dbmodel.ProxySiteTypeDefaultRss:
	//			// 读取xml  数据源_分组_媒体
	//
	//			break
	//		}
	//	}(sourceInfo)
	//}
	//wg.Wait()

	return nil
}

func (this Apiv1ToznabDao) getGroupMediaXml(groupId string, query string) {
	queryArr := strings.Split(query, "+")
	if len(queryArr) == 2 {
		// 找到是那个sonarr
	}

}
