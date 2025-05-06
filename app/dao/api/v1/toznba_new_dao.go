package v1

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/apiv1/toznab"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/groups"
	"XArr-Rss/app/service/match"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/util/date_util"
	"XArr-Rss/util/hash"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"XArr-Rss/util/uri_util"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ApiTorznabNewDao struct {
}

func (this ApiTorznabNewDao) Api(c *gin.Context) {
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
	case "search", "tvsearch":
		this.Search(c, groupId)
	default:
		logsys.Error("请求方式错误:"+t, "Toznab")
		c.String(200, "服务运行中 v2")
	}

}

func (this ApiTorznabNewDao) Caps(c *gin.Context) {
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
				Available:       "yes",
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

// 搜索 search tvsearch
func (this ApiTorznabNewDao) Search(c *gin.Context, groupId string) {
	//searchType := c.Query("t")
	query := c.Query("q")
	//title := c.Query("title")

	//cat := c.Query("cat")

	// 找到sonarr媒体信息
	sonarrMediaInfo, episodeNum := this.querySonarrMediaInfoByQuery(c, query)
	//log.Println(sonarrMediaInfo)
	// 如果找到媒体信息 则只匹配对应媒体生成的数据 否则匹配全部内容
	if sonarrMediaInfo != nil {
		if groupId == "all" {

			groupList := groups.GroupsService{}.GetGroupList()
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
					Item: this.GetDefaultItem(),
				},
			}

			for _, groupInfo := range groupList {
				// 找到分组媒体数据
				groupMedia, _ := groups.GroupMediaService{}.GetGroupMediaInfoByGroupAndSonarrId(groupInfo.Id, sonarrMediaInfo.SonarrId)
				if groupMedia != nil {
					ress := this._searchGroupMedia(c, groupMedia, episodeNum)

					res.Channel.Item = append(res.Channel.Item, ress.Channel.Item...)

				}
			}

			c.XML(200, res)

		} else {
			// 找到分组媒体数据
			groupMedia, _ := groups.GroupMediaService{}.GetGroupMediaInfoByGroupAndSonarrId(helper.StrToInt(groupId), sonarrMediaInfo.SonarrId)
			if groupMedia != nil {
				res := this._searchGroupMedia(c, groupMedia, episodeNum)
				c.XML(200, res)
				return
			}
		}

	} else {
		// 找到分组媒体数据
		res := this._searchGroupMediaCache(c, groupId, episodeNum)
		c.XML(200, res)
		return

	}

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
			Item: this.GetDefaultItem(),
		},
	}
	c.XML(200, res)

}

// 查询媒体信息
func (this ApiTorznabNewDao) querySonarrMediaInfoByQuery(c *gin.Context, query string) (*dbmodel.Media, int) {
	tvdbid := c.Query("tvdbid")
	imdbid := c.Query("imdbid")
	rid := c.Query("rid")
	///////// 解析query 进行搜索
	queryArr := strings.Split(query, "+")
	queryLen := len(queryArr)

	episodeNum := 0
	if queryLen > 0 {
		episodeNum = cast.ToInt(strings.TrimLeft(queryArr[queryLen-1], "0"))
	}

	// 采用tvdbId 进行搜索
	media := medias.MediaService{}.GetMediaInfoByOtherId(tvdbid, imdbid, rid)
	if media != nil {
		return media, episodeNum
	}
	if query == "" {
		return nil, episodeNum
	}

	// 获取左侧的标题
	query = strings.Join(queryArr[:queryLen-1], "+")

	// 将标题进行搜索
	media = medias.MediaService{}.GetMediaInfoByQuery(query)

	return media, episodeNum
}

// 搜索分组媒体下的数据项 1. 获取数据源/匹配结果  2. 根据匹配结果 判断搜索季是否合适
func (this ApiTorznabNewDao) _searchGroupMedia(c *gin.Context, groupMedia *dbmodel.GroupMedia, episodeNum int) toznab.Apiv1ToznabSearchRes {
	// 查询分组媒体中使用的数据源列表
	sourceList := sources.SourcesService{}.GetSourcesList(true, groupMedia.UseSource...)

	searchResult := &model.RssRoot{}

	query := uri_util.UriToUrlValues(c.Request.RequestURI)
	for _, sourceInfo := range sourceList {
		if sourceInfo.ProxySiteType == dbmodel.ProxySiteTypeJacket || sourceInfo.ProxySiteType == dbmodel.ProxySiteTypeProwlarr {
			// 代理搜索
			tmpItems := this._searchProxy(query, groupMedia, sourceInfo)
			if tmpItems != nil {
				searchResult.Channel.Item = append(searchResult.Channel.Item, tmpItems...)
			}
		} else if sourceInfo.ProxySiteType == dbmodel.ProxySiteTypeDefaultRss {
			// 读取已经生成的 分组数据源匹配结果
			// 读取缓存
			// 目录 group_1/mediaid_sourceId.xml
			tmpRoot := this._readGroupMediaCache(groupMedia, sourceInfo)
			if tmpRoot != nil && tmpRoot.Channel.Item != nil {
				searchResult.Channel.Item = append(searchResult.Channel.Item, tmpRoot.Channel.Item...)
			}
		}
	}

	// 判断是否需要搜索指定集数
	if episodeNum > 0 {
		for index := len(searchResult.Channel.Item) - 1; index >= 0; index-- {
			item := searchResult.Channel.Item[index]
			if item.MinEpisode == item.MaxEpisode || item.MaxEpisode == 0 {
				if item.MinEpisode == episodeNum {
					continue
				}
			}
			if item.MinEpisode < item.MaxEpisode {
				if item.MinEpisode <= episodeNum && episodeNum <= item.MaxEpisode {
					continue
				}
			}
			if item.OldMinEpisode == item.OldMaxEpisode || item.OldMaxEpisode == 0 {
				if item.OldMinEpisode == episodeNum {
					continue
				}
			}
			if item.OldMinEpisode < item.OldMaxEpisode {
				if item.OldMinEpisode <= episodeNum && episodeNum <= item.OldMaxEpisode {
					continue
				}
			}

			searchResult.Channel.Item = append(searchResult.Channel.Item[:index], searchResult.Channel.Item[index+1:]...)

		}
	}

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
			Item: this.GetDefaultItem(),
		},
	}

	// 处理匹配

	if searchResult != nil && searchResult.Channel.Item != nil {
		searchResult.Channel.Item = this.matchItems(c.Query("q"), c.Query("title"), c.Query("ep"), c.Query("season"), c.Query("tvdbid"), searchResult.Channel.Item)

		for _, v := range searchResult.Channel.Item {
			pubDate := v.PubDate
			if pubDate == "" {
				pubDate = date_util.TimeNowStr()
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
					Text: pubDate,
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
	}

	return res
	//c.XML(200, res)
}

func (this ApiTorznabNewDao) _searchGroupMediaCache(c *gin.Context, groupId string, episodeNum int) toznab.Apiv1ToznabSearchRes {
	searchResult := &model.RssRoot{}

	tmpRoot := this._readGroupCache(groupId)
	if tmpRoot != nil && tmpRoot.Channel.Item != nil {
		searchResult.Channel.Item = append(searchResult.Channel.Item, tmpRoot.Channel.Item...)
	}
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
			Item: this.GetDefaultItem(),
		},
	}
	if searchResult != nil && searchResult.Channel.Item != nil {
		searchResult.Channel.Item = this.matchItems(c.Query("q"), c.Query("title"), c.Query("ep"), c.Query("season"), c.Query("tvdbid"), searchResult.Channel.Item)

		for _, v := range searchResult.Channel.Item {
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

		// 判断是否需要搜索指定集数
		if episodeNum > 0 {
			for index := len(searchResult.Channel.Item) - 1; index >= 0; index-- {
				item := searchResult.Channel.Item[index]
				if item.MinEpisode == item.MaxEpisode || item.MaxEpisode == 0 {
					if item.MinEpisode == episodeNum {
						continue
					}
				}
				if item.MinEpisode < item.MaxEpisode {
					if item.MinEpisode <= episodeNum && episodeNum <= item.MaxEpisode {
						continue
					}
				}

				if item.OldMinEpisode == item.OldMaxEpisode || item.OldMaxEpisode == 0 {
					if item.OldMinEpisode == episodeNum {
						continue
					}
				}
				if item.OldMinEpisode < item.OldMaxEpisode {
					if item.OldMinEpisode <= episodeNum && episodeNum <= item.OldMaxEpisode {
						continue
					}
				}
				searchResult.Channel.Item = append(searchResult.Channel.Item[:index], searchResult.Channel.Item[index+1:]...)

			}
		}
	}

	return res
}

var searchTotal = make(map[string]int64)

func (this ApiTorznabNewDao) _searchProxy(query *url.Values, groupMedia *dbmodel.GroupMedia, sourceInfo *dbmodel.Sources) []model.RssResultItem {
	query.Del("tvdbid")
	query.Del("imdbid")

	// 计算搜索频率
	hashQuery := hash.Md5{}.HashString(query.Get("q"))
	if _, ok := searchTotal[hashQuery]; !ok {
		searchTotal[hashQuery] = 0
	}

	// 增加搜索频率
	searchTotal[hashQuery]++
	if searchTotal[hashQuery] > 50000000 {
		searchTotal[hashQuery] = 1
	}

	if groupMedia.MediaInfo.CnTitle != "" && searchTotal[hashQuery]%2 != 0 {
		query.Set("q", groupMedia.MediaInfo.CnTitle)
	}
	_, sourcecItems, err := sources.SourcesService{}.SearchSourceItems(sourceInfo, query)
	if err != nil {
		return nil
	}
	// 这里要进行数据匹配
	groupsMediaCache := match.ParseGroupMediaSourceItems(groupMedia, sourcecItems)
	if groupsMediaCache == nil {
		return nil
	}
	// 返回匹配结果
	return groupsMediaCache.Channel.Item

}

func (this ApiTorznabNewDao) matchItems(query, queryTitle, episode, seasonStr, tvdbid string, items []model.RssResultItem) (ret []model.RssResultItem) {
	//query := c.Query("q")
	qArr := strings.Split(query, "+")
	//queryTitle := c.Query("title")
	//episode := c.Query("ep")
	if len(qArr) == 2 {
		queryTitle = qArr[0]
		episode = qArr[1]
	}

	queryTitle = strings.ToLower(helper.ReplaceRegString(queryTitle))

	season := helper.StrToInt(seasonStr)

	for _, v := range items {
		if v.Title.Text == "XArr-Rss 默认占位使用" {
			continue
		}

		if queryTitle != "" {
			// 搜索标题
			if !strings.Contains(strings.ToLower(helper.ReplaceRegString(v.Title.Text)), queryTitle) && !strings.Contains(strings.ToLower(helper.ReplaceRegString(v.OriginalTitle.Text)), queryTitle) && !strings.Contains(strings.ToLower(helper.ReplaceRegString(v.OtherTitle.Text)), queryTitle) {
				continue
			}
		}

		if episode != "" {
			if v.OldMinEpisode <= helper.StrToInt(episode) && helper.StrToInt(episode) <= v.OldMaxEpisode {
				// 允许搜索
			} else if v.MinEpisode <= helper.StrToInt(episode) && helper.StrToInt(episode) <= v.MaxEpisode {
				// 允许搜索
			} else {
				// 如果有季的话 就不判断 0了
				if v.Season >= 0 {

				} else {
					// 不允许搜索
					continue
				}

			}
		}

		if season >= 0 {
			if v.Season != season && v.Season >= 0 {
				continue
			}
		}
		if tvdbid != "" {
			if v.OthderId.TvdbId != tvdbid {
				continue
			}
		}
		//if episode != "" && v.MinEpisode != -1 {
		//	if v.MaxEpisode != -1 {
		//		// 是连续的
		//		if !(v.MinEpisode <= helper.StrToInt(episode) && helper.StrToInt(episode) <= v.MaxEpisode) {
		//			continue
		//		}
		//	} else {
		//		if v.MinEpisode != helper.StrToInt(episode) {
		//			continue
		//		}
		//	}
		//}

		ret = append(ret, v)
	}
	return
}

// 读取分组缓存
func (this ApiTorznabNewDao) _readGroupMediaCache(groupMedia *dbmodel.GroupMedia, info *dbmodel.Sources) *model.RssRoot {
	groupXml := &model.RssRoot{}

	// 提取groupXml
	groupId := helper.IntToStr(groupMedia.GroupId)
	filePath := appconf.AppConf.ConfDir + "/trans/group_" + groupId + "/" + strconv.Itoa(groupMedia.Id) + "/" + strconv.Itoa(info.Id) + ".xml"

	_, err := os.Stat(filePath)
	if err != nil {
		logsys.Error("获取缓存文件[%s] 异常:%s", "torznab", filePath, err.Error())
		return nil
	}
	// 解析groupxml
	groupXmlData, err := os.ReadFile(filePath)
	if err != nil {
		logsys.Error("读取缓存文件[%s] 异常:%s", "torznab", filePath, err.Error())
		return nil
	}
	err = xml.Unmarshal(groupXmlData, groupXml)
	if err != nil {
		logsys.Error("解析缓存文件[%s] 异常:%s", "torznab", filePath, err.Error())
		return nil
	}
	return groupXml
}

// 读取分组缓存
func (this ApiTorznabNewDao) _readGroupCache(groupId string) *model.RssRoot {
	groupXml := &model.RssRoot{}

	// 提取groupXml
	filePath := appconf.AppConf.ConfDir + "/trans/group_" + groupId + ".xml"

	_, err := os.Stat(filePath)
	if err != nil {
		logsys.Error("获取缓存文件[%s] 异常:%s", "torznab", filePath, err.Error())
		return nil
	}
	// 解析groupxml
	groupXmlData, err := os.ReadFile(filePath)
	if err != nil {
		logsys.Error("读取缓存文件[%s] 异常:%s", "torznab", filePath, err.Error())
		return nil
	}
	err = xml.Unmarshal(groupXmlData, groupXml)
	if err != nil {
		logsys.Error("解析缓存文件[%s] 异常:%s", "torznab", filePath, err.Error())
		return nil
	}
	return groupXml
}

func (d ApiTorznabNewDao) GetDefaultItem() []toznab.Apiv1ToznabSearchChannelItem {
	return []toznab.Apiv1ToznabSearchChannelItem{
		{
			Title: toznab.Apiv1ToznabSearchChannelItemTitle{
				Text: "XArr占位 QQ群:996973766",
			},
			OriginalTitle: toznab.Apiv1ToznabSearchChannelItemTitle{
				Text: "XArr占位 QQ群:996973766",
			},
			OtherTitle: toznab.Apiv1ToznabSearchChannelItemTitle{
				Text: "",
			},
			Guid: toznab.Apiv1ToznabSearchChannelItemGuid{
				Text: "qq群:996973766",
			},
			XArrRssindexer: toznab.Apiv1ToznabSearchChannelItemXArrRssindexer{},
			Comments:       toznab.Apiv1ToznabSearchChannelItemComments{},
			PubDate: toznab.Apiv1ToznabSearchChannelItemPubDate{
				Text: "2022-04-21T05:38:00",
			},
			Size: toznab.Apiv1ToznabSearchChannelItemSize{
				Text: "996973766",
			},
			Link: toznab.Apiv1ToznabSearchChannelItemLink{
				Text: "https://xarr.52nyg.com",
			},
			Category: toznab.Apiv1ToznabSearchChannelItemCategory{},
			Enclosure: toznab.Apiv1ToznabSearchChannelItemEnclosure{
				URL:    "https://xarr.52nyg.com",
				Length: "996973766",
				Type:   "application/x-bittorrent",
				//Type: "application/x-nzb",
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
					Value: "",
				},
				{
					Text:  "",
					Name:  "imdb",
					Value: "",
				},
				{
					Text:  "",
					Name:  "tmdbid",
					Value: "",
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
		},
	}
}
