package medias

import (
	"XArr-Rss/app/global/cache"
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model/dbmodel"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type MediaEpisodeService struct {
}

func (this MediaEpisodeService) GetEpisodes(sonarrId int) []dbmodel.MediaEpisodeList {
	var a []dbmodel.MediaEpisodeList
	db.MainDb.Model(dbmodel.MediaEpisodeList{}).Where("sonarr_id = ?", sonarrId).Order("episode_number desc").Find(&a)
	return a
}

func (this MediaEpisodeService) GetSeasonEpisodes(sonarrId int, season int) []dbmodel.MediaEpisodeList {
	var a []dbmodel.MediaEpisodeList
	db.MainDb.Model(dbmodel.MediaEpisodeList{}).Where("sonarr_id = ? and season = ?", sonarrId, season).Order("season_number desc ,episode_number desc").Find(&a)
	return a
}

func (this MediaEpisodeService) GetSeasonEpisode(sonarrId int, season int, episode int) *dbmodel.MediaEpisodeList {
	cacheData, found := cache.GocacheClient.Get(cache.CACHE_KEY_GetSeasonEpisode + fmt.Sprintf("%d-%d-%d", sonarrId, season, episode))
	if found {
		return cacheData.(*dbmodel.MediaEpisodeList)
	}
	var a *dbmodel.MediaEpisodeList
	first := db.MainDb.Model(dbmodel.MediaEpisodeList{}).
		Where("sonarr_id = ? and season_number = ? and episode_number = ? ", sonarrId, season, episode).
		Order("episode_number desc").
		First(&a)
	if first.Error != nil || first.RowsAffected <= 0 {
		return nil
	}
	cache.GocacheClient.Set(cache.CACHE_KEY_GetSeasonEpisode+fmt.Sprintf("%d-%d-%d", sonarrId, season, episode), a, time.Minute*2)

	return a
}

func (this MediaEpisodeService) GetSeasonAbEpisode(sonarrId int, season int, episode int) *dbmodel.MediaEpisodeList {
	cacheData, found := cache.GocacheClient.Get(cache.CACHE_KEY_GetSeasonAbEpisode + fmt.Sprintf("%d-%d-%d", sonarrId, season, episode))
	if found {
		return cacheData.(*dbmodel.MediaEpisodeList)
	}

	result := &dbmodel.MediaEpisodeList{}
	first := db.MainDb.Model(dbmodel.MediaEpisodeList{}).
		Where("sonarr_id = ? and season_number = ? and absolute_episode_number = ? ", sonarrId, season, episode).
		Order("absolute_episode_number desc").
		First(result)
	if first.Error != nil || first.RowsAffected <= 0 {
		return nil
	}
	cache.GocacheClient.Set(cache.CACHE_KEY_GetSeasonAbEpisode+fmt.Sprintf("%d-%d-%d", sonarrId, season, episode), result, time.Minute*2)

	return result
}

func (receiver MediaEpisodeService) SaveEpisode(media *dbmodel.MediaEpisodeList) error {
	return db.MainDb.Save(media).Error
}

func (this MediaEpisodeService) CreateEpisode(data *dbmodel.MediaEpisodeList) error {
	create := db.MainDb.Model(dbmodel.MediaEpisodeList{}).Create(data)
	return create.Error
}

func (receiver MediaEpisodeService) Upsert(media *dbmodel.MediaEpisodeList) *gorm.DB {
	create := db.MainDb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sonarr_id"}, {Name: "season_number"}, {Name: "episode_number"}}, // key colume
		// 也可以用 map [ string ] interface {} { "role" : "user" }
		DoUpdates: clause.AssignmentColumns([]string{
			"episode_title",
			"air_date",
			"air_date_utc",
			"episode_id",
			"has_file",
			"absolute_episode_number",
		}), // column needed to be updated
	}, clause.Insert{Modifier: "or ignore"}).Create(media)
	return create

}
