package sources

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/apiv1/toznab"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/matchv2"
	"XArr-Rss/app/service/source_item"
	"XArr-Rss/util/date_util"
	"XArr-Rss/util/hash"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/xml"
	"errors"
	"github.com/spf13/cast"
	"net/url"
	"regexp"
	"strings"
)

// 此方法不涉及缓存时间
func (this SourcesService) SearchSourceItems(source *dbmodel.Sources, queryData *url.Values) (*model.RssRoot, []dbmodel.SourceItem, error) {
	// 生成uri

	proxy := ""
	if source.UseProxy == 1 {
		proxy = options.GetOption(options.OptionsGlobalProxy)
	}
	rssResult := &model.RssRoot{}
	err := errors.New("")
	err = nil
	if source.ProxySiteType == dbmodel.ProxySiteTypeJacket || source.ProxySiteType == dbmodel.ProxySiteTypeProwlarr {
		// 生成查询数据
		if queryData == nil {
			queryData = &url.Values{}
		}
		if !queryData.Has("apikey") {
			queryData.Set("apikey", source.ProxySiteApiKey)
		}

		if !queryData.Has("cat") {
			queryData.Set("cat", "5000,5070,100001")
		}
		if !queryData.Has("t") {
			queryData.Set("t", "tvsearch")
		}

		curlHelper := &helper.CurlHelper{}
		logsys.Info("torznab请求:Url:%s", "数据源", source.Url+"?"+queryData.Encode())
		uriT := strings.TrimRight(source.Url, "/")

		cerr, data, _ := curlHelper.Init(nil).SetProxy(proxy).Get(uriT, queryData, true)
		if cerr != nil {
			return nil, nil, nil
		}
		// 解析为rssResult
		//rssResult :=

		rssResult, err = ParseTorznabResult(data)
	} else {
		// 直接请求 后面关键字可能需要这个
		searchTitle := ""
		if queryData != nil {
			searchTitle = queryData.Get("auto_search_title")
		}
		rssResult, err = GetRssUrlResult(source.Url, proxy, source.Name, source.AutoSearch, searchTitle)
	}
	if err != nil {
		// 获取数据错误
		return nil, nil, err
	}
	if rssResult.Channel.Item == nil {
		return nil, nil, errors.New("没有获取到数据")
	}
	// 倒序
	Reverse(&rssResult.Channel.Item)
	////////////////////////// 过滤筛选数据
	// 限制数据源是否包含名称
	if source.Regex.MustHave != "" {
		reg, err := regexp.Compile(source.Regex.MustHave)
		if err != nil {
			logsys.Error("匹配标题必须存在规则失败(%s)：%s", "数据源", source.Name, err.Error())
		} else {
			var newList []model.RssResultItem
			for _, v2 := range rssResult.Channel.Item {
				if reg.MatchString(v2.Title.Text) {
					newList = append(newList, v2)
				}
			}
			rssResult.Channel.Item = newList
		}
	}
	// 判断是否必须不能包含
	if source.Regex.MustDotHave != "" {
		reg, err := regexp.Compile(source.Regex.MustDotHave)
		if err != nil {
			logsys.Error("匹配标题必须存在规则失败(%s)：%s", "数据源", source.Name, err.Error())
		} else {
			var newList []model.RssResultItem
			for _, v2 := range rssResult.Channel.Item {
				if !reg.MatchString(v2.Title.Text) {
					newList = append(newList, v2)
				}
			}
			rssResult.Channel.Item = newList
		}
	}

	// 限制单次最高获取数量
	if len(rssResult.Channel.Item) > source.MaxReadCount && source.MaxReadCount > 0 {
		// 截取 如果总数大于截取数量
		rssResult.Channel.Item = rssResult.Channel.Item[len(rssResult.Channel.Item)-source.MaxReadCount:]

	}

	//rssResult.Channel
	//hashVars := []string{}
	//for _, sourceItem := range rssResult.Channel.Item {
	//	// 查询hash 是否在history中
	//	hashVars = append(hashVars, hash.Md5{}.HashString(sourceItem.Guid.Text+sourceItem.Title.Text+sourceItem.Enclosure.Length+cast.ToString(source.Id)))
	//}

	insertCount := 0
	sc := source_item.SourceItemService{}
	souceItems := []dbmodel.SourceItem{}
	for _, sourceItem := range rssResult.Channel.Item {

		var hashVar = ""
		if source.ProxySiteType == dbmodel.ProxySiteTypeJacket || source.ProxySiteType == dbmodel.ProxySiteTypeProwlarr {
			sourceItem.Guid.Text = hash.Md5{}.HashString(sourceItem.Guid.Text + sourceItem.Title.Text + sourceItem.Enclosure.Length)

			hashVar = hash.Md5{}.HashString(sourceItem.Guid.Text)

		} else {
			// jacket 处理跟 rss不相同
			sourceItem.Guid.Text = hash.Md5{}.HashString(sourceItem.Guid.Text + sourceItem.Title.Text + sourceItem.Enclosure.Length + sourceItem.Link + sourceItem.PubDate + cast.ToString(source.Id) + source.Url)

			// 查询hash 是否在history中
			hashVar = hash.Md5{}.HashString(sourceItem.Guid.Text + sourceItem.Title.Text + sourceItem.Enclosure.Length + sourceItem.Link)
		}

		item := sc.GetSourceItemByHash(source.Id, hashVar)
		if item == nil {
			// 不在数据库则插入数据库
			data := dbmodel.SourceItem{
				SourceId: source.Id,
				Hash:     hashVar,
				PubDate:  sourceItem.PubDate,
				Enclosure: dbmodel.SourceItemEnclosure{
					Type:   sourceItem.Enclosure.Type,
					Length: sourceItem.Enclosure.Length,
					Url:    sourceItem.Enclosure.Url,
				},
				Content: "",
				Link:    sourceItem.Link,
				Guid:    sourceItem.Guid.Text,
				Title:   sourceItem.Title.Text,
			}

			// 获取解析结果
			parse := torrent_title_parse.TorrentTitleParse{}
			parseResult := parse.Parse(sourceItem.Title.Text)
			sourceItem.Enclosure.Length = matchv2.MatchV2{}.ParseSourceItemLength(sourceItem.Enclosure.Length)

			// 转换解析信息
			data.ParseInfo = ConvertToSourceParseInfo(parseResult)

			runResult := source_item.SourceItemService{}.UpsertItem(data)
			if runResult.Error != nil {
				logsys.Error("保存数据源内容失败:%s", "同步数据源", runResult.Error.Error())
				continue
			} else if runResult.RowsAffected > 0 {
				logsys.Debug("解析标题信息:%s", "测试", sourceItem.Title.Text)
			}

			souceItems = append(souceItems, data)
			insertCount++
		} else {
			souceItems = append(souceItems, *item)
		}
	}

	// 3. 删除数据源旧数据
	err = source_item.SourceItemService{}.DeleteExpireItem(source.Id, source.CacheDay)
	if err != nil {
		logsys.Error("删除过期缓存数据错误:%s", "数据源", err.Error())
	}
	// 4. 删除多余的数据
	err = source_item.SourceItemService{}.DeleteLimitItem(source.Id, source.MaxCount)
	if err != nil {
		logsys.Error("删除过多数据错误:%s", "数据源", err.Error())
	}
	//////////////////
	logsys.Info("已保存缓存信息 ["+source.Name+"]新增数据:%d", "数据源", insertCount)

	return rssResult, souceItems, nil

}

func ParseTorznabResult(result []byte) (ret *model.RssRoot, err error) {
	ret = &model.RssRoot{}
	// 解析代理之后获取的数据源信息
	res := &toznab.Apiv1ToznabSearchRes{}
	err = xml.Unmarshal(result, res)
	if err != nil {
		logsys.Error("解析代理返回结果错误:%s 原始内容:%s", "代理请求", err.Error(), string(result))
		return nil, err
	}

	// 格式化到rssRoot里面
	for _, item := range res.Channel.Item {
		pubDate := item.PubDate.Text
		if pubDate == "" {
			pubDate = date_util.TimeNowStr()
		}

		// 找到合适的种子地址
		torrentUri := item.Enclosure.URL
		// 第一个匹配 guid
		if strings.Index(item.Guid.Text, "https://") == 0 || strings.Index(item.Guid.Text, "http://") == 0 {
			// 判断是否是种子地址
			if strings.Contains(item.Guid.Text, ".torrent") {
				// 使用这个东西
				torrentUri = item.Guid.Text
			}
		} else if strings.Index(item.Guid.Text, "magnet:?") == 0 {
			torrentUri = item.Guid.Text
		}

		ret.Channel.Item = append(ret.Channel.Item, model.RssResultItem{
			Title: model.CDATA{
				Text: item.Title.Text,
			},
			PubDate: pubDate,
			Enclosure: model.RssResultItemEnclosure{
				Type:   item.Enclosure.Type,
				Length: item.Enclosure.Length,
				Url:    torrentUri, //item.Enclosure.URL,
			},
			Link: item.Comments.Text,
			Guid: model.RssResultItemGuid{
				IsPermaLink: false,
				Text:        item.Guid.Text,
			},
		})
	}

	return ret, nil
}

// 数组倒序函数
func Reverse(arr *[]model.RssResultItem) {
	var temp model.RssResultItem
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}
