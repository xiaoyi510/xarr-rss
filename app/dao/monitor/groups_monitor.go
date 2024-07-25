package monitor

import (
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/groups"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/util/array"
	"XArr-Rss/util/golimit"
	"XArr-Rss/util/logsys"
	"XArr-Rss/util/regexp_ext"
	"strings"
	"time"
)

type GroupMonitor struct {
}

// 循环查询
func (this *GroupMonitor) syncGroupsRange() {
	logsys.Info("开始同步分组数据", "同步分组数据")
	for {
		this.SyncGroups(true)
		time.Sleep(10 * time.Minute)
	}
}

// 同步所有分组信息
func (this GroupMonitor) SyncGroups(useCache bool) *model.RssRoot {
	// 查询group

	groupList := groups.GroupsService{}.GetGroupList()
	syn := golimit.NewGoLimit(2)

	for _, v := range groupList {
		syn.Add()
		go func(v dbmodel.Group) {
			// 自动同步分组Sonarr剧集信息
			if v.AutoInsertSonarr != dbmodel.AutoInsertSonarrNone {
				// 同步Sonarr新数据
				this.SyncSonarrToGroup(v)
			}
			// 生成分组数据
			err, _ := groups.GroupsService{}.SyncGroupItem(v, useCache)
			if err != nil {
				//continue
			}
			syn.Done()
		}(v)
	}
	syn.Wait()
	return groups.GroupsService{}.SyncGroups()
}

// 从Sonarr同步分组媒体数据
func (this GroupMonitor) SyncSonarrToGroup(group dbmodel.Group) {
	// 判断VIP才可以使用
	if !variable.ServerState.IsVip {
		return
	}
	if group.AutoInsertSonarr == dbmodel.AutoInsertSonarrNone {
		return
	}

	logsys.Debug("同步Sonarr媒体数据到分组中", "同步Sonarr媒体")
	// 获取分组下的所有媒体数据
	mediaList := medias.MediaService{}.GetMediaList()
	if len(mediaList) == 0 {
		// 没有媒体
		return
	}
	// 判断是否没有监听过  如果没有 则设置最后一个sonarrId为监控开始点
	if group.LastInsertSonarrId <= 0 {
		group.LastInsertSonarrId = mediaList[len(mediaList)-1].SonarrId
		//return
		// 如果是只需要新数据 第一次则返回
		if group.AutoInsertSonarr == dbmodel.AutoInsertSonarrNew {
			return
		}
	}

	// 查询分组下的媒体信息
	groupMediaList := groups.GroupMediaService{}.GetGroupMedias(group.Id)
	removeGroupMedia := options.GetOption(options.SonarrUnMonitoredRemoveGroupMedia)
	nAddGroupMedia := options.GetOption(options.SonarrUnMonitoredNAddGroupMedia)

	for _, media := range mediaList {
		// 判断是否sonarr中被禁止监听
		if media.Monitored == 0 {
			if removeGroupMedia == "1" {
				// 删除分组下的媒体
				groups.GroupMediaService{}.DeleteMediaBySonarrId(media.SonarrId)
			}

			// 不自动添加到分组
			if nAddGroupMedia == "1" {
				continue
			}
		}

		if (group.AutoInsertSonarr == dbmodel.AutoInsertSonarrNew && media.SonarrId > group.LastInsertSonarrId) ||
			(group.AutoInsertSonarr == dbmodel.AutoInsertSonarrTags && media.Tags != nil && len(media.Tags) > 0 && group.Tags != "" &&
				this.matchTag(media.Tags, group.Tags)) ||
			group.AutoInsertSonarr == dbmodel.AutoInsertSonarrAuto {
			inGroupMedia := false

			// 判断分组中没有此媒体的信息
			for _, item := range groupMediaList {
				if item.SonarrId == media.SonarrId {
					inGroupMedia = true
					break
				}
			}
			if !inGroupMedia {
				// 没有影片数据
				// 增加分组影片数据
				err := groups.GroupMediaService{}.AutoAddGroupMedia(group.Id, group.GroupTemplateId, media.SonarrId)
				if err != nil {
					// 自动添加影片数据异常
					logsys.Error("Sonarr自动添加影片数据异常:%s", "同步分组数据", err.Error())
					continue
				}
				logsys.Info("同步Sonarr媒体到分组[%s]:%s", "同步分组数据", group.Name, media.CnTitle)

			}
		}

	}
}

func (this GroupMonitor) matchTag(str []string, regexStr string) bool {
	if len(regexStr) == 0 {
		return true
	}
	if len(regexStr) > 4 && regexStr[0:4] == "reg:" {
		return regexp_ext.MatchStringArr(str, regexStr[0:4]) != nil
	} else {
		return array.ArrayHasInArrayAnd(str, strings.Split(regexStr, ","))
	}
}
