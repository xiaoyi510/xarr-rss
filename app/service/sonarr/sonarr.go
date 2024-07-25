package sonarr

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type SonarrService struct {
}

// 测试连接状态
func (this SonarrService) TestConnect() (error, *model.SonarrSystemStatusRes) {
	// 判断Sonarr是否已配置成功
	if appconf.AppConf.Service.Sonarr.Apikey == "" || appconf.AppConf.Service.Sonarr.Host == "" {
		return logsys.Error("Sonarr信息未填写", "SonarrId"), nil
	}

	// 获取剧集列表
	uri := strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/v3/system/status?apikey=" + appconf.AppConf.Service.Sonarr.Apikey
	err, result, _ := helper.CurlHelper{}.GetUri(uri, nil, nil, true)
	if err != nil {
		return logsys.Error("Sonarr查询失败:%s", "SonarrId", err.Error()), nil
	}

	resp := &model.SonarrSystemStatusRes{}
	err = json.Unmarshal(result, resp)
	if err != nil {
		return logsys.Error("Sonarr查询异常:%s", "SonarrId", err.Error()), nil
	}

	return nil, resp
}

// 同步SonarrTags
func (this SonarrService) SyncSonarrTags() error {
	// 判断Sonarr是否已配置成功
	if appconf.AppConf.Service.Sonarr.Apikey == "" || appconf.AppConf.Service.Sonarr.Host == "" {
		return logsys.Error("Sonarr信息未填写", "SyncSonarrTags")
	}

	// 获取剧集列表
	uri := strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/tag?apikey=" + appconf.AppConf.Service.Sonarr.Apikey
	if appconf.AppConf.Service.Sonarr.Version == "4" {
		uri = strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/v3/tag?apikey=" + appconf.AppConf.Service.Sonarr.Apikey
	}

	err, result, _ := helper.CurlHelper{}.GetUri(uri, nil, nil, true)
	if err != nil {
		return logsys.Error("Sonarr信息查询失败:%s", "SyncSonarrTags", err.Error())
	}
	// 解析返回数据
	var resp model.SonarrTags
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return logsys.Error("Sonarr信息查询异常:%s", "SyncSonarrTags", err.Error())
	}
	tagList := make(map[int]string)
	for _, tag := range resp {
		tagList[tag.Id] = tag.Label
	}
	variable.SonarrTags = tagList

	logsys.Info("同步Sonarr标签完成", "SyncSonarrTags")
	return nil
}

// 同步sonarr 剧集到本地
func (this SonarrService) SyncSonarrToLocal(isRef ...bool) error {
	// 判断Sonarr是否已配置成功
	if appconf.AppConf.Service.Sonarr.Apikey == "" || appconf.AppConf.Service.Sonarr.Host == "" {
		return logsys.Error("Sonarr信息未填写", "SonarrId")
	}

	// 获取剧集列表
	uri := strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/series?apikey=" + appconf.AppConf.Service.Sonarr.Apikey
	if appconf.AppConf.Service.Sonarr.Version == "4" {
		uri = strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/v3/series?apikey=" + appconf.AppConf.Service.Sonarr.Apikey
	}
	err, result, _ := helper.CurlHelper{}.GetUri(uri, nil, nil, true)
	if err != nil {
		return logsys.Error("Sonarr信息查询失败:%s", "SyncSonarrToLocal", err.Error())
	}

	// 解析返回数据
	var resp []model.SonarrSeries

	err = json.Unmarshal(result, &resp)
	if err != nil {
		return logsys.Error("Sonarr信息查询异常:%s", "SonarrId", err.Error())
	}

	// 记录没有sonarr id 的记录
	mediasKeys := make(map[int]int)
	mediaList := medias.MediaService{}.GetMediaList()
	for _, v := range mediaList {
		mediasKeys[v.SonarrId] = 0
	}

	// 处理为config数据
	for _, v := range resp {
		// 判断媒体是否存在
		atitles := []string{v.Title}
		for _, v2 := range v.AlternateTitles {
			atitles = append(atitles, helper.StrReplace(v2.Title, []string{"-"}, []string{" "}))
		}
		if appconf.AppConf.Service.Sonarr.Version == "4" {
			v.SeasonCount = v.Statistics.SeasonCount
			v.TotalEpisodeCount = v.Statistics.TotalEpisodeCount
			v.EpisodeCount = v.Statistics.EpisodeCount
		}

		mediasKeys[v.Id] = 1

		monitored := 0
		if v.Monitored {
			monitored = 1
		}
		item := &dbmodel.Media{
			SonarrId:          v.Id,
			CnTitle:           "",
			OriginalTitle:     v.Title,
			Overview:          v.Overview,
			AlternateTitles:   atitles,
			SeasonCount:       v.SeasonCount,
			TotalEpisodeCount: v.TotalEpisodeCount,
			EpisodeCount:      v.EpisodeCount,
			Monitored:         monitored,
			Seasons:           v.Seasons,
			Year:              v.Year,
			TvdbId:            v.TvdbId,
			TvRageId:          v.TvRageId,
			TvMazeId:          strconv.Itoa(v.TvMazeId),
			ImdbId:            v.ImdbId,
			TitleSlug:         v.TitleSlug,
			SeriesType:        v.SeriesType,
			Tags:              this.GetTagsName(v.Tags),
		}

		// 获取已存数据
		mediaInfo := medias.MediaService{}.GetMediaInfoByList(mediaList, v.Id)
		if mediaInfo != nil && len(isRef) == 0 {
			// 中文数据
			if mediaInfo.Overview != "" {
				item.Overview = mediaInfo.Overview
			}
			item.CnTitle = mediaInfo.CnTitle
			item.TmdbId = mediaInfo.TmdbId
			item.ErrTime = mediaInfo.ErrTime
			item.Titles = mediaInfo.Titles
		}
		if len(isRef) > 0 && isRef[0] {
			// 需要重新加载
			item.Overview = v.Overview
			item.CnTitle = ""
			item.TmdbId = ""
			item.ErrTime = 0
			item.Titles = make([]string, 0)
		}

		// 更新数据
		err = medias.MediaService{}.Upsert(item)
		if err != nil {
			logsys.Error("保存Soarr媒体失败:%s", "SyncSonarrToLocal", err.Error())
		}

	}

	// 删除没有使用到的sonarr信息
	for k, v := range mediasKeys {
		if v == 0 {
			this.DeleteSonarrMedia(k)
		}
	}

	//logsys.Info("Sonarr同步完成", "SonarrId")

	// 刷新group 里面的影片信息
	return nil
}
func (this SonarrService) GetTagsName(tags []int) []string {
	var ret []string
	for _, tag := range tags {
		tagLabel, ok := variable.SonarrTags[tag]
		if ok {
			ret = append(ret, tagLabel)
		}
	}
	return ret
}

// 同步sonarr里面的剧集信息
func (this SonarrService) SyncEpisode(sonarrId int) {
	//logsys.Debug("开始同步剧集[%d]", "SonarrId", sonarrId)
	// 获取剧集列表
	uri := strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/episode?apikey=" + appconf.AppConf.Service.Sonarr.Apikey + "&seriesId=" + strconv.Itoa(sonarrId)

	if appconf.AppConf.Service.Sonarr.Version == "4" {
		uri = strings.TrimRight(appconf.AppConf.Service.Sonarr.Host, "/") + "/api/v3/episode?apikey=" + appconf.AppConf.Service.Sonarr.Apikey + "&seriesId=" + strconv.Itoa(sonarrId)
	}
	err, result, _ := helper.CurlHelper{}.GetUri(uri, nil, nil, true)
	if err != nil {
		logsys.Error("SonarrId 剧集信息查询失败:%s", "同步Sonarr剧集", err.Error())
		return
	}

	var resp []model.SonarrEpisode

	err = json.Unmarshal(result, &resp)
	if err != nil {
		logsys.Error("SonarrId 剧集信息解析异常:%s", "同步Sonarr剧集", err.Error())
		return
	}

	/////////////////////获取数据库中的集列表
	episodeList := medias.MediaEpisodeService{}.GetEpisodes(sonarrId)
	// 组合剧集信息
	episodeMapList := map[string]dbmodel.MediaEpisodeList{}
	for _, v := range episodeList {
		episodeMapList[strconv.Itoa(v.SeasonNumber)+"_"+strconv.Itoa(v.SeasonNumber)+"_"+strconv.Itoa(v.EpisodeNumber)] = v
	}
	//if variable.IsDebug {
	//	logsys.Debug("开始保存", "同步Sonarr剧集")
	//}
	for _, v := range resp {
		hasFile := 0
		if v.HasFile {
			hasFile = 1
		}
		val := &dbmodel.MediaEpisodeList{
			SonarrId:              sonarrId,
			SeasonNumber:          v.SeasonNumber,
			EpisodeNumber:         v.EpisodeNumber,
			EpisodeTitle:          v.Title,
			AirDate:               v.AirDate,
			AirDateUtc:            v.AirDate,
			HasFile:               hasFile,
			EpisodeId:             v.Id,
			AbsoluteEpisodeNumber: v.AbsoluteEpisodeNumber,
		}

		upRest := medias.MediaEpisodeService{}.Upsert(val)
		//
		//if variable.IsDebug {
		//	logsys.Debug("保存 SonarrId:%d-S%d-E%d 保存结果:%d", "同步Sonarr剧集", val.SonarrId, val.SeasonNumber, val.EpisodeNumber, upRest.RowsAffected)
		//}
		if upRest.Error != nil {
			logsys.Error("保存剧集数据失败:%s", "同步剧集", err.Error())
		}
	}
}

func (this SonarrService) GetSonarrTitleById(sonarrId int) string {
	mediaInfo := medias.MediaService{}.GetMediaInfo(sonarrId)
	if mediaInfo != nil {
		if mediaInfo.CnTitle != "" {
			return mediaInfo.CnTitle
		}

		return mediaInfo.OriginalTitle
	}
	return ""
}

func (this SonarrService) GetSonarrMediaById(sonarrId int) (error, dbmodel.Media) {
	mediaInfo := medias.MediaService{}.GetMediaInfo(sonarrId)
	if mediaInfo == nil {
		return errors.New("没有找到"), dbmodel.Media{}
	}
	return nil, *mediaInfo
}

func (this SonarrService) DeleteSonarrMedia(sonarrId int) {
	// 删除所有 SonarrId相同的
	db.MainDb.Model(dbmodel.GroupMedia{}).Where("sonarr_id = ?", sonarrId).Delete(&dbmodel.GroupMedia{})
	db.MainDb.Model(dbmodel.Media{}).Where("sonarr_id = ?", sonarrId).Delete(&dbmodel.Media{})

}

func (this SonarrService) GetDownload() {

}
