package groups

import (
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model/apiv1/group"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/logsys"
	"gorm.io/gorm/clause"
	"time"
)

type GroupMediaService struct {
}

func (this GroupMediaService) AddGroupMedias(groupId int, req group.Apiv1GroupMediaAdd) error {
	item := &dbmodel.GroupMedia{
		SonarrId:        req.SonarrId,
		Regex:           req.Regex,
		Language:        req.Language,
		Quality:         req.Quality,
		UseSource:       req.UseSource,
		FilterPushGroup: req.FilterPushGroup,
		GroupId:         groupId,
		EchoTitleAnime:  req.EchoTitleAnime,
		EchoTitleTv:     req.EchoTitleTv,
	}

	return db.MainDb.Create(item).Error
}

func (this GroupMediaService) Upsert(data *dbmodel.GroupMedia) error {
	create := db.MainDb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "group_id"}, {Name: "sonarr_id"}}, // key colume
		// 也可以用 map [ string ] interface {} { "role" : "user" }
		//DoUpdates: clause.AssignmentColumns([]string{}), // column needed to be updated
	}, clause.Insert{Modifier: "or ignore"}).Create(data)
	return create.Error
}

// 自动添加分组媒体
func (this GroupMediaService) AutoAddGroupMedia(groupId int, groupTemplateId int32, sonarrId int) error {
	regexs := []dbmodel.GroupRegex{}

	regexs = append(regexs, dbmodel.GroupRegex{
		MatchType: "auto",
		Reg:       "", // 默认就是中文+英文名
		RegType:   dbmodel.REG_TYPE_DEFAULT,
		Season:    -1,
		Offset:    0,
	})
	info := &dbmodel.GroupMedia{
		GroupId:         groupId,
		SonarrId:        sonarrId,
		Language:        "-1",
		Quality:         "-1",
		Regex:           regexs,
		UseSource:       []string{"-1"},
		FilterPushGroup: nil,
		EchoTitleAnime:  "",
		EchoTitleTv:     "",
	}

	if groupTemplateId > 0 {
		// 查询分组模板信息
		groupTemplateInfo := GroupTemplateService{}.FindOne(groupTemplateId)
		if groupTemplateInfo != nil {
			info.FromGroupTemp = int(groupTemplateId)
			if groupTemplateInfo.Language != "" {
				info.Language = groupTemplateInfo.Language
			}
			if groupTemplateInfo.Quality != "" {
				info.Quality = groupTemplateInfo.Quality
			}

			if len(groupTemplateInfo.Regex) > 0 {
				info.Regex = groupTemplateInfo.Regex
			}
			if len(groupTemplateInfo.UseSource) > 0 {
				info.UseSource = groupTemplateInfo.UseSource
			}
			if len(groupTemplateInfo.FilterPushGroup) > 0 {
				info.FilterPushGroup = groupTemplateInfo.FilterPushGroup
			}
			if groupTemplateInfo.EchoTitleAnime != "" {
				info.EchoTitleAnime = groupTemplateInfo.EchoTitleAnime
			}

			if groupTemplateInfo.EchoTitleTv != "" {
				info.EchoTitleTv = groupTemplateInfo.EchoTitleTv
			}

		}
	}

	return this.Upsert(info)
}

func (this GroupMediaService) CreateGroupMedias(item *dbmodel.GroupMedia) error {
	return db.MainDb.Create(item).Error
}

// 获取分组下的媒体数据
func (this GroupMediaService) GetGroupMedias(groupId int) []*dbmodel.GroupMedia {
	groupMedias := []*dbmodel.GroupMedia{}
	find := db.MainDb.Table("group_media").Joins("MediaInfo").Where("group_id = ?", groupId).Find(&groupMedias)
	if find.Error != nil {
		logsys.Error("查询分组媒体数据错误:%s", "GetGroupMedias", find.Error.Error())
	}
	return groupMedias
}

// 获取分组下的媒体数据
func (this GroupMediaService) GetGroupMediaList() []*dbmodel.GroupMedia {
	groupMedias := []*dbmodel.GroupMedia{}
	find := db.MainDb.Table("group_media").Find(&groupMedias)
	if find.Error != nil {
		logsys.Error("查询分组媒体数据错误:%s", "GetGroupMedias", find.Error.Error())
	}
	return groupMedias
}

// 根据关键字筛选
func (this GroupMediaService) GetGroupMediasAndQuery(groupId int, query string) []dbmodel.GroupMedia {
	groupMedias := []dbmodel.GroupMedia{}
	find := db.MainDb.Table("group_media").Joins("MediaInfo").
		Where("group_id = ?", groupId)
	if query != "" {
		find.Where("cn_title like ? or original_title like ? ", "%"+query+"%", "%"+query+"%")
	}
	find.
		Order("id desc").
		Find(&groupMedias)
	if find.Error != nil {
		logsys.Error("查询分组媒体数据错误:%s", "GetGroupMedias", find.Error.Error())
	}
	return groupMedias
}

func (this GroupMediaService) GroupMediaExits(groupId int, sonarrId int) bool {
	var count int64
	find := db.MainDb.Model(dbmodel.GroupMedia{}).Where("group_id = ? and sonarr_id =? ", groupId, sonarrId).Count(&count)
	if find.Error != nil {
		logsys.Error("查询分组媒体数据错误:%s", "GetGroupMedias", find.Error.Error())

	}
	return count > 0
}

func (this GroupMediaService) GroupMediaIdExits(sonarr_id int) bool {
	var count int64
	find := db.MainDb.Model(dbmodel.GroupMedia{}).Where("sonarr_id =? ", sonarr_id).Count(&count)
	if find.Error != nil {
		logsys.Error("查询分组媒体数据错误:%s", "GetGroupMedias", find.Error.Error())

	}
	return count > 0
}

func (this GroupMediaService) GetGroupMediaInfo(id int) (*dbmodel.GroupMedia, error) {
	var groupMedia dbmodel.GroupMedia
	find := db.MainDb.Model(dbmodel.GroupMedia{}).Joins("MediaInfo").Where("id = ?", id).First(&groupMedia)
	if find.Error != nil {
		return nil, find.Error
	}
	return &groupMedia, nil
}

func (this GroupMediaService) GetGroupMediaInfoByGroupAndSonarrId(groupId int, sonarrId int) (*dbmodel.GroupMedia, error) {
	var groupMedia dbmodel.GroupMedia
	find := db.MainDb.Model(dbmodel.GroupMedia{}).Joins("MediaInfo").Where("group_id = ? and group_media.sonarr_id = ?", groupId, sonarrId).First(&groupMedia)
	if find.Error != nil {
		return nil, find.Error
	}
	return &groupMedia, nil
}

func (this GroupMediaService) Save(data *dbmodel.GroupMedia) error {
	save := db.MainDb.Save(data)
	return save.Error
}

func (this GroupMediaService) DeleteMedia(data *dbmodel.GroupMedia) error {
	tx := db.MainDb.Delete(data)
	return tx.Error
}

func (this GroupMediaService) DeleteMediaBySonarrId(sonarrId int) error {
	tx := db.MainDb.Model(&dbmodel.GroupMedia{}).Where("sonarr_id = ?", sonarrId).Delete(&dbmodel.GroupMedia{})
	return tx.Error
}

// 获取总数量
func (this GroupMediaService) GetTotalCount() int {
	var ret int64
	db.MainDb.Model(&dbmodel.GroupMedia{}).Count(&ret)
	return int(ret)
}

// 获取今日新增
func (this GroupMediaService) GetTodayNewCount() interface{} {
	var ret int64
	db.MainDb.Model(&dbmodel.GroupMedia{}).Where("created_at > datetime(?)", time.Now()).Count(&ret)
	return int(ret)
}

// 根据分组ID删除媒体数据
func (this GroupMediaService) DeleteByGroupId(groupId int) error {
	rest := db.MainDb.Model(dbmodel.GroupMedia{}).Where("group_id = ?", groupId).Delete(&dbmodel.GroupMedia{})
	return rest.Error
}
