package monitor

import (
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/service/sonarr"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/util/logsys"
	"time"
)

type CenterMonitor struct {
}

// 监控中心
func (this CenterMonitor) Run() {
	logsys.Debug("开始启动监控服务", "CenterMonitor")
	this.BeforeRun()
	this.monitorSonarr()
	this.monitorSource()
	this.monitorGroup()
	this.monitorDownload()
	this.monitorPushGroup()
	this.monitorOther()
	this.monitorBackup()
	this.monitorSourceAutoSearch()
	logsys.Debug("启动监控服务完成", "CenterMonitor")

}

func (this CenterMonitor) BeforeRun() {
	cha := make(chan int)
	go func() {
		sonarr.SonarrService{}.SyncSonarrTags()
		sonarr.SonarrService{}.SyncSonarrToLocal()
		cha <- 1
	}()
	select {
	case <-cha:
		return

	case <-time.After(time.Second * 5):
		logsys.Debug("等待Sonarr初始化超时", "CenterMonitor")
		return

	}
}

func (this CenterMonitor) monitorSonarr() {
	// Sonarr 相关
	sonarrMonitor := SonarrMonitor{}
	go sonarrMonitor.SyncSonArr()
	go sonarrMonitor.SyncSonarrMediaInfo()
	go sonarrMonitor.SyncSonarrMediaEpisode()
	go sonarrMonitor.SyncSonarrTags()
}

// 监听分组
func (this CenterMonitor) monitorGroup() {
	group := GroupMonitor{}
	go group.syncGroupsRange()
}

// 监听下载器
func (this CenterMonitor) monitorDownload() {
	qb := QbittorrentMonitor{}
	go qb.SyncQbit()
	tr := TransmissionMonitor{}
	go tr.syncTransmission()

}

// 发布组监听器
func (this CenterMonitor) monitorPushGroup() {
	pushGroup := PushGroupMonitor{}
	go pushGroup.SyncPushGroup()
}

func (this CenterMonitor) monitorOther() {
	// 刷新数据源匹配内容
	go sources.ReParseTitles(false)

	// 检测http登录信息
	go variable.CheckTokenExpires()
}

// 数据源监听器
func (this CenterMonitor) monitorSource() {
	source := SourcesMonitor{}
	go source.SyncDownSources()
}

func (this CenterMonitor) monitorBackup() {
	go DatabaseMonitor()
}

func (this CenterMonitor) monitorSourceAutoSearch() {
	source := SourcesMonitor{}
	go source.SourceAutoSearch()
}
