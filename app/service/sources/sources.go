package sources

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/cache"
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/source_item"
	"XArr-Rss/util/array"
	"XArr-Rss/util/golimit"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"strconv"
	"time"
)

type SourcesService struct {
}

// 获取数据源ID列表
func (this SourcesService) GetSourcesIds() []string {
	get, found := cache.GocacheClient.Get(cache.CACHE_KEY_SOURCEIDS)
	if found {
		logsys.Debug("使用缓存数据", "GetSourcesIds")
		return get.([]string)
	}
	// 查询所有数据源信息
	sources := []dbmodel.Sources{}
	find := db.MainDb.Model(dbmodel.Sources{}).Select("id").Find(&sources)

	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return []string{}
	}
	ret := []string{}
	for _, v := range sources {
		ret = append(ret, strconv.Itoa(int(v.Model.Id)))
	}
	cache.GocacheClient.Set(cache.CACHE_KEY_SOURCEIDS, ret, time.Minute*5)
	return ret
}

func (this SourcesService) GetSourcesList(searchEnableStatus bool, ids ...string) []*dbmodel.Sources {
	// 查询所有数据源信息
	sources := []*dbmodel.Sources{}
	m := db.MainDb.Model(dbmodel.Sources{})
	if len(ids) > 0 {
		if !array.InArray("-1", ids) {
			m = m.Where("id in ?", ids)
		}
	}

	if searchEnableStatus {
		m = m.Where("status = 1")
	}

	find := m.Find(&sources)
	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return sources
	}
	return sources
}

func (this SourcesService) GetSourcesListAutoSearch() []*dbmodel.Sources {
	// 查询所有数据源信息
	sources := []*dbmodel.Sources{}
	m := db.MainDb.Model(dbmodel.Sources{})

	m = m.Where("status = 1")
	m = m.Where("auto_search != ''")

	find := m.Find(&sources)
	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return sources
	}
	return sources
}

func (this SourcesService) GetSourcesCount() int64 {
	// 查询所有数据源信息
	sources := int64(0)
	m := db.MainDb.Model(dbmodel.Sources{})

	find := m.Count(&sources)
	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return sources
	}
	return sources
}

func (this SourcesService) GetSourcesInfo(id string) (*dbmodel.Sources, error) {
	// 查询所有数据源信息
	sources := &dbmodel.Sources{}
	find := db.MainDb.Model(dbmodel.Sources{}).Where("id = ?", id).First(sources)
	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return sources, find.Error
	}
	return sources, nil

}

func (this SourcesService) CreateSource(data *dbmodel.Sources) error {
	// 添加数据源都清理下缓存
	cache.GocacheClient.Delete(cache.CACHE_KEY_SOURCEIDS)

	return db.MainDb.Model(dbmodel.Sources{}).Create(data).Error
}

func (this SourcesService) Delete(info *dbmodel.Sources) error {
	cache.GocacheClient.Delete(cache.CACHE_KEY_SOURCEIDS)

	return db.MainDb.Model(dbmodel.Sources{}).Delete(info).Error

}

func (this SourcesService) Save(info *dbmodel.Sources) error {
	cache.GocacheClient.Delete(cache.CACHE_KEY_SOURCEIDS)

	save := db.MainDb.Save(info)
	return save.Error
}

// 获取数据源中的数据项
func GetSourceItems(useSourceId []string, useCache bool) []dbmodel.SourceItem {
	ret := []dbmodel.SourceItem{}

	for _, sourceId := range useSourceId {
		if sourceId == "-1" {
			useSourceId = SourcesService{}.GetSourcesIds()
			break
		}
	}
	// result, found := cache.GocacheClient.Get(cache.CACHE_KEY_GetSourcesItemList)
	//	if found {
	//		return result.([]*dbmodel.SourceItem)
	//	}
	//cache.GocacheClient.Set(cache.CACHE_KEY_GetSourcesItemList, ret, time.Minute*5)

	for _, sourceId := range useSourceId {
		if useCache {
			cacheList, found := cache.GocacheClient.Get(cache.CACHE_KEY_GetSourceItems + "_" + sourceId)
			if found {
				ret = append(ret, cacheList.([]dbmodel.SourceItem)...)
				logsys.Debug("使用缓存数据:%s", "缓存", cache.CACHE_KEY_GetSourceItems+"_"+sourceId)
				continue
			}
		}

		list := source_item.SourceItemService{}.GetSourcesItemList(helper.StrToInt(sourceId), 0)
		if useCache {
			cache.GocacheClient.Set(cache.CACHE_KEY_GetSourceItems+"_"+sourceId, list, time.Minute*6)
		}
		ret = append(ret, list...)
	}
	return ret
}

// 获取数据源中的数据项
func GetSourceItemsCore(useSourceId []string) []dbmodel.SourceItem {
	ret := []dbmodel.SourceItem{}

	for _, sourceId := range useSourceId {
		if sourceId == "-1" {
			useSourceId = SourcesService{}.GetSourcesIds()
			break
		}
	}

	for _, sourceId := range useSourceId {
		list := source_item.SourceItemService{}.GetSourcesItemList(helper.StrToInt(sourceId), 0)
		ret = append(ret, list...)
	}
	return ret
}

func ReParseTitles(reParse bool, sourceIds ...string) {
	logsys.Debug("开始重新生成匹配缓存", "生成数据源匹配缓存")
	// 为空则刷新全部
	if len(sourceIds) == 0 {
		sourceIds = SourcesService{}.GetSourcesIds()
	}
	for _, sourceId := range sourceIds {
		ReParseSourceTitle(sourceId, reParse)
	}
}

func ReParseSourceTitle(sourceId string, reParse bool) {

	list := source_item.SourceItemService{}.GetSourcesItemList(helper.StrToInt(sourceId), 0)
	waitLimit := golimit.NewGoLimit(10)
	for _, sourceItem2 := range list {
		waitLimit.Add()

		//sourceItem := sourceItem
		go func(sourceItem dbmodel.SourceItem) {
			defer waitLimit.Done()

			// 避免已经匹配过的再来匹配一次
			if sourceItem.ParseInfo != nil && reParse == false {
				return
			}
			parse := torrent_title_parse.TorrentTitleParse{}
			parseInfo := parse.Parse(sourceItem.Title)
			if parseInfo.MinEpisode > 0 {
				sourceItem.ParseInfo = ConvertToSourceParseInfo(parseInfo)
			}
			err := source_item.SourceItemService{}.UpsertItem(sourceItem)
			if err.Error != nil {
				logsys.Error("更新数据源匹配内容失败:%s", "生成数据源匹配缓存", err.Error.Error())
				return
			}
		}(sourceItem2)

	}

	waitLimit.Wait()
}

func ConvertToSourceParseInfo(result *torrent_title_parse.MatchResult) *dbmodel.SourceItemParseInfo {
	if result.MinEpisode <= 0 && result.MaxEpisode <= 0 && result.Season <= 0 {
		return nil
	}
	return &dbmodel.SourceItemParseInfo{
		AnalyzeTitle:      result.AnalyzeTitle,
		AudioEncode:       result.AudioEncode,
		VideoEncode:       result.VideoEncode,
		MediaType:         result.MediaType,
		MinEpisode:        result.MinEpisode,
		MaxEpisode:        result.MaxEpisode,
		Version:           result.Version,
		Season:            result.Season,
		ReleaseGroup:      result.ReleaseGroup,
		QualityResolution: result.QualityResolution,
		Language:          result.Language,
		ProductionCompany: result.ProductionCompany,
		Subtitles:         result.Subtitles,
	}
}

// 同步数据源下载到本地 同步第一页无搜索内容的数据
func (this SourcesService) SyncSource(source *dbmodel.Sources, isReload bool) *model.RssRoot {
	if source.LastRefreshTime > time.Now().Unix() && isReload == false {
		return nil
	}

	// 不需要从这里同步
	if source.AutoSearch != "" && isReload == false {
		return nil
	}

	logsys.Info("开始获取数据源信息[%s]:%s", "数据源", source.Name, source.Url)
	// 判断使用什么方式进行同步

	rssResult, _, err := this.SearchSourceItems(source, nil)
	if err != nil {
		return nil
	}

	// 最低5分钟一次
	if source.RefreshTime <= 5*60 {
		source.RefreshTime = 5 * 60
	}
	source.LastRefreshTime = time.Now().Unix() + source.RefreshTime

	// 保存刷新时间
	this.Save(source)

	return rssResult
}
