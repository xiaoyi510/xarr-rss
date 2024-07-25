package monitor

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/app/service/sonarr"
	"XArr-Rss/util/logsys"
	"time"
)

type SonarrMonitor struct {
}

// 下载Sonarr剧集信息到本地
func (this *SonarrMonitor) SyncSonArr() {
	logsys.Info("开始监听Sonarr数据源", "SonarrId")
	for {
		// 获取间隔时间
		if appconf.AppConf.Service.Sonarr.RefreshTime < 1 {
			appconf.AppConf.Service.Sonarr.RefreshTime = 1
		}
		// 延迟刷新时间
		time.Sleep(time.Duration(appconf.AppConf.Service.Sonarr.RefreshTime) * time.Minute)

		// 刷新
		sonarr.SonarrService{}.SyncSonarrToLocal()

	}
}

// 同步Sonarr数据中文信息
func (this *SonarrMonitor) SyncSonarrMediaInfo() {
	logsys.Info("开始监听Sonarr数据中文信息", "Themoviedb")
	time.Sleep(3 * time.Second)
	for {
		medias.SyncMediaRemoteInfo()
		time.Sleep(1 * time.Minute)
	}
	// 同步数据
}

func (this *SonarrMonitor) SyncSonarrMediaEpisode() {
	logsys.Info("开始监听Sonarr数据集数", "Sonarr剧集")
	time.Sleep(3 * time.Second)
	for {
		logsys.Info("开始同步Sonarr剧集信息", "Sonarr剧集")

		// 提交获取更新集信息
		// 查询剧集列表
		list := medias.MediaService{}.GetMediaList()
		for _, item := range list {
			sonarr.SonarrService{}.SyncEpisode(item.SonarrId)
		}

		time.Sleep(3 * time.Minute)
	}
}

// 刷新标签
func (this *SonarrMonitor) SyncSonarrTags() {
	logsys.Info("开始监听Sonarr标签", "SyncSonarrTags")
	time.Sleep(3 * time.Second)
	for {
		logsys.Info("开始监听Sonarr标签", "SyncSonarrTags")

		// 提交获取更新集信息
		// 查询剧集列表
		err := sonarr.SonarrService{}.SyncSonarrTags()
		if err != nil {
			time.Sleep(5 * time.Minute)
			continue
		}

		time.Sleep(3 * time.Minute)
	}
}
