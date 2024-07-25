package monitor

import (
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/groups"
	medias2 "XArr-Rss/app/service/medias"
	"XArr-Rss/app/service/source_item"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/util/golimit"
	"XArr-Rss/util/logsys"
	"net/url"
	"time"
)

type SourcesMonitor struct {
}

// 更新数据源到本地
func (this *SourcesMonitor) SyncDownSources() {
	logsys.Debug("开始监听数据源", "监控-数据源")
	wg := golimit.NewGoLimit(2)

	for {
		logsys.Debug("获取数据源列表", "监控-数据源")

		// 获取所有数据源
		sourceList := sources.SourcesService{}.GetSourcesList(true)

		// 开始刷新数据源 限制最大同步数量 避免CPU爆炸
		for _, sourceInfo := range sourceList {
			if sourceInfo.AutoSearch != "" {
				continue
			}
			wg.Add()
			go func(v *dbmodel.Sources) {
				defer wg.Done()
				// 开始同步
				sources.SourcesService{}.SyncSource(v, false)

			}(sourceInfo)
		}
		logsys.Debug("开始等数据源同步完成", "监控-数据源")
		wg.Wait()
		logsys.Debug("数据源全部同步完成", "监控-数据源")

		// 删除不在数据源中的数据项
		del := source_item.SourceItemService{}.DeleteErrItem()
		if del.Error == nil && del.RowsAffected > 0 {
			logsys.Info("删除多余的数据源信息:%d个", "监控-数据源", del.RowsAffected)
		} else if del.Error != nil {
			logsys.Error("删除多余的数据源失败:%s", "监控-数据源", del.Error.Error())
		}

		// 数据源全部刷新完成后等5分钟
		logsys.Debug("开始等待下一波", "监控-数据源")

		time.Sleep(15 * time.Second)
	}
}

func (this *SourcesMonitor) SourceAutoSearch() {
	logsys.Debug("开始监控自动检索项目", "监控-自动检索")

	for true {
		// 查询所有数据源
		sourceList := sources.SourcesService{}.GetSourcesListAutoSearch()
		if !variable.ServerState.IsVip {
			time.Sleep(5 * time.Second)
			logsys.Debug("开通VIP后方可使用", "监控-自动检索")
			continue
		}

		// 查询所有sonar媒体
		medias := medias2.MediaService{}.GetMediaList()
		for _, source := range sourceList {
			if source.LastRefreshTime > time.Now().Unix() {
				continue
			}

			//  监控最新列表
			sources.SourcesService{}.SearchSourceItems(source, &url.Values{})

			for _, media := range medias {
				// 判断是否在分组中
				exGroup := groups.GroupMediaService{}.GroupMediaIdExits(media.SonarrId)
				if !exGroup {
					continue
				}

				// 没有在监控中的不搜索
				if media.Monitored != 1 {
					continue
				}

				if media.CnTitle != "" {
					switch source.AutoSearch {
					case "mikan":
						qry := &url.Values{}
						qry.Set("auto_search_title", media.CnTitle)
						sources.SourcesService{}.SearchSourceItems(source, qry)
					}
				}
				// 每个搜索间隔5秒钟
				time.Sleep(5 * time.Second)
			}

			// 最低5分钟一次
			if source.RefreshTime <= 5*60 {
				source.RefreshTime = 5 * 60
			}
			source.LastRefreshTime = time.Now().Unix() + source.RefreshTime
			// 保存刷新时间
			sources.SourcesService{}.Save(source)
		}

		time.Sleep(2 * time.Minute)
	}

}
