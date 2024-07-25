package medias

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/sdk/themoviedb"
	"XArr-Rss/util/golimit"
	"XArr-Rss/util/logsys"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type MediaService struct {
}

// 同步所有媒体数据
func SyncMediaRemoteInfo(isref ...bool) {
	rand.Seed(time.Now().UnixNano())
	mediaList := MediaService{}.GetMediaList()

	for _, v := range mediaList {
		if v.TmdbId == "" {
			log.Println("无tmdbid信息的媒体", v.OriginalTitle)
		}
	}

	if appconf.AppConf.Service.Themoviedb.Apikey != "" {
		proxy := appconf.AppConf.Service.GetTMDBProxy()

		sdk := themoviedb.GetTheMoviedbSdk(appconf.AppConf.Service.Themoviedb.Apikey, proxy)

		// 设置协程池
		syn := golimit.NewGoLimit(3)
		for _, v := range mediaList {
			syn.Add()
			go func(mediaInfo dbmodel.Media) {
				defer syn.Done()
				if mediaInfo.ErrTime > time.Now().Unix()-5*60*60 {
					return
				}
				if mediaInfo.CnTitle == "" || mediaInfo.CnTitle == mediaInfo.OriginalTitle || mediaInfo.TmdbId == "" {
					// 根据themoviedb 查询标题
					if mediaInfo.ImdbId != "" {
						err, m := sdk.FindByExternalIds(mediaInfo.ImdbId, "imdb_id")
						if err != nil {
							logsys.Error("查询Themoviedb异常", "Themoviedb", err.Error())
							return
						}
						media := MediaService{}.GetMediaInfo(mediaInfo.SonarrId)
						if media == nil {
							logsys.Error("Sonarr媒体可能被删除了", "Themoviedb")
							return
						}

						if len(m.TvResults) > 0 {
							media.Overview = m.TvResults[0].Overview
							media.TmdbId = strconv.Itoa(m.TvResults[0].Id)
							if m.TvResults[0].Name == mediaInfo.OriginalTitle || (m.TvResults[0].Name == m.TvResults[0].OriginalName && m.TvResults[0].OriginalLanguage != "zh") {
								logsys.Info("判断目标标题可能没有中文,标题:%s,原始标题%s,原始语言%s", m.TvResults[0].Name, m.TvResults[0].OriginalName, m.TvResults[0].OriginalLanguage)
								media.ErrTime = time.Now().Unix()
								err = MediaService{}.Save(media)
								return
							}
							media.CnTitle = m.TvResults[0].Name
							logsys.Info("搜索到中文信息id:%s title:%s 中文标题:%s", "Themoviedb", media.SonarrId, media.OriginalTitle, media.CnTitle)
						} else {
							media.ErrTime = time.Now().Unix()
						}
						err = MediaService{}.Save(media)
						if err != nil {
							logsys.Error("Sonarr媒体保存失败", "SyncMediaRemoteInfo", err.Error())
							return
						}
						mediaInfo = *media

					}

					if mediaInfo.TvdbId != 0 && mediaInfo.TmdbId == "" {
						err, m := sdk.FindByExternalIds(strconv.Itoa(mediaInfo.TvdbId), "tvdb_id")
						if err != nil {
							logsys.Error("查询Themoviedb异常", "Themoviedb", err.Error())
							return
						}

						media := MediaService{}.GetMediaInfo(mediaInfo.SonarrId)
						if media == nil {
							logsys.Error("Sonarr媒体可能被删除了", "Themoviedb")
							return
						}

						if len(m.TvResults) > 0 {
							media.Overview = m.TvResults[0].Overview
							media.CnTitle = m.TvResults[0].Name
							logsys.Info("搜索到中文信息id:%mediaInfo title:%s 中文标题:%s", "Themoviedb", media.SonarrId, media.OriginalTitle, media.CnTitle)
							media.TmdbId = strconv.Itoa(m.TvResults[0].Id)
						} else {
							media.ErrTime = time.Now().Unix()
						}

						err = MediaService{}.Save(media)
						if err != nil {
							logsys.Error("Sonarr媒体保存失败", "SyncMediaRemoteInfo", err.Error())
							return
						}
						mediaInfo = *media

					}
				}

				if len(isref) > 0 && isref[0] == true {
					mediaInfo.Titles = make([]string, 0)
				}

				if len(mediaInfo.Titles) == 0 && mediaInfo.TmdbId != "" {
					// 查询themvoiedb其他的替代标题
					err, titles := sdk.AlternativeTvTitles(mediaInfo.TmdbId)
					if err != nil {
						logsys.Error("查询Themoviedb备用标题异常", "Themoviedb", err.Error())
						return
					}

					media := MediaService{}.GetMediaInfo(mediaInfo.SonarrId)
					if media == nil {
						logsys.Error("Sonarr媒体可能被删除了", "Themoviedb")
						return
					}
					if len(titles.Results) == 0 {
						media.ErrTime = time.Now().Unix()
						err = MediaService{}.Save(media)
						if err != nil {
							logsys.Error("Sonarr媒体保存失败", "SyncMediaRemoteInfo", err.Error())
							return
						}

						return
					}
					if media.Titles == nil {
						media.Titles = make([]string, 0)
					}
					for _, titleResult := range titles.Results {
						has := false
						allowIsoArr := []string{"CN", "ES", "JP", "KR", "MX", "TW", "US"}
						for _, allowIso := range allowIsoArr {
							if titleResult.Iso31661 == allowIso {

								has = true
								break
							}
						}
						if has {
							if titleResult.Iso31661 == "CN" && media.CnTitle == "" {
								media.CnTitle = titleResult.Title
							}
							if !utils.Contains(media.Titles, titleResult.Title) {
								media.Titles = append(media.Titles, titleResult.Title)
							}
						}

					}
					err = MediaService{}.Save(media)
					if err != nil {
						logsys.Error("Sonarr媒体保存失败", "SyncMediaRemoteInfo", err.Error())
						return
					}

				}
			}(v)
			time.Sleep(5 * time.Millisecond)
		}
		//logsys.Error("等待中", "媒体")
		syn.Wait()
	} else {
		logsys.Error("ThemovieDB数据为空,跳过查询中文标题", "媒体")
	}

	//logsys.Info("中文媒体数据处理完成", "媒体")
}

func (receiver MediaService) GetMaxSonarrId() int {
	a := dbmodel.Media{}
	db.MainDb.Model(dbmodel.Media{}).Select("sonarr_id").Order("sonarr_id desc").First(&a)
	return (a.SonarrId)
}

func (receiver MediaService) GetMediaList() []dbmodel.Media {
	a := []dbmodel.Media{}
	db.MainDb.Model(dbmodel.Media{}).Order("sonarr_id desc").Find(&a)
	return a
}

func (receiver MediaService) GetMediaListForTitle(name string) []dbmodel.Media {
	a := []dbmodel.Media{}
	db.MainDb.Model(dbmodel.Media{}).Where("cn_title like ? or original_title like ?", "%"+name+"%", "%"+name+"%").Order("sonarr_id desc").Find(&a)
	return a
}

func (receiver MediaService) GetMediaCount() int64 {
	var a int64
	db.MainDb.Model(dbmodel.Media{}).Count(&a)
	return a
}

func (receiver MediaService) GetMediaInfo(sonarrId int) *dbmodel.Media {
	a := &dbmodel.Media{}
	first := db.MainDb.Model(dbmodel.Media{}).Where("sonarr_id = ?", sonarrId).First(&a)
	if first.RowsAffected <= 0 {
		return nil
	}
	return a
}

func (receiver MediaService) GetMediaInfoByList(list []dbmodel.Media, sonarrId int) *dbmodel.Media {
	for _, v := range list {
		if v.SonarrId == sonarrId {
			return &v
		}
	}
	return nil
}

func (receiver MediaService) Save(media *dbmodel.Media) error {
	return db.MainDb.Save(media).Error
}

func (receiver MediaService) Upsert(media *dbmodel.Media) error {
	create := db.MainDb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sonarr_id"}}, // key colume
		// 也可以用 map [ string ] interface {} { "role" : "user" }
		DoUpdates: clause.AssignmentColumns([]string{
			"original_title",
			"alternate_titles",
			"season_count",
			"total_episode_count",
			"episode_count",
			"seasons",
			"year",
			"tvdb_id",
			"tvrage_id",
			"tvmaze_id",
			"imdb_id",
			"title_slug",
			"series_type",
			"tags",
			"titles",
			"cn_title",
			"tmdb_id",
			"overview",
			"monitored",
		}), // column needed to be updated
	}, clause.Insert{Modifier: "or ignore"}).Create(media)
	return create.Error

}

// 根据扩展ID进行搜索
func (receiver MediaService) GetMediaInfoByOtherId(tvdbid string, imdbid string, rid string) *dbmodel.Media {
	if imdbid == "" && rid == "" && tvdbid == "" {
		return nil
	}
	query := db.MainDb.Model(dbmodel.Media{})
	if imdbid != "" {
		query.Where("imdb_id = ?", imdbid)
	}
	if tvdbid != "" {
		query.Where("tvdb_id = ?", tvdbid)
	}
	if rid != "" {
		query.Where("tvrage_id = ?", rid)
	}

	ret := &dbmodel.Media{}
	query.
		Order("sonarr_id desc").
		Find(ret)
	return ret
}

func (receiver MediaService) GetMediaInfoByQuery(queryKey string) *dbmodel.Media {
	if queryKey == "" {
		return nil
	}
	query := db.MainDb.Model(dbmodel.Media{})
	if queryKey != "" {
		query = query.Where("alternate_titles like ?", `%`+queryKey+`%`)
	}

	ret := &dbmodel.Media{}
	res := query.
		Order("sonarr_id desc").
		Find(ret)
	if res.Error != nil {
		logsys.Error("查询数据失败:%s", "查询媒体数据", res.Error.Error())
		return nil
	}
	if res.RowsAffected <= 0 {
		return nil
	}
	return ret
}
