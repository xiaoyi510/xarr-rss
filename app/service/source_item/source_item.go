package source_item

import (
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/array"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type SourceItemService struct {
}

func (this SourceItemService) GetSourcesItemListCount(sourceId int, query string) int64 {
	// 查询所有数据源信息
	var sources int64
	m := db.MainDb.Model(dbmodel.SourceItem{})
	m = m.Where("source_id = ?", sourceId)
	if query != "" {
		m = m.Where("title like ?", "%"+query+"%")
	}
	find := m.Count(&sources)
	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return sources
	}
	return sources
}

// 判断数据是否存在
func (this SourceItemService) Exits(sourceId int, hashVar string) int64 {
	rest := int64(0)

	count := db.MainDb.Model(dbmodel.SourceItem{}).Where("source_id = ? ", sourceId).Where("hash = ?", hashVar).Count(&rest)
	if count.Error != nil {
		return 0
	}
	return rest

}

func (this SourceItemService) GetSourcesItemListPage(sourceId int, query string, page int, limit int) []dbmodel.SourceItem {
	// 查询所有数据源信息
	sources := []dbmodel.SourceItem{}
	m := db.MainDb.Model(dbmodel.SourceItem{})
	if page > 0 {
		m = m.Offset((page - 1) * limit)
	}
	if limit > 0 {
		m = m.Limit(limit)
	}
	m = m.Where("source_id = ?", sourceId).Order("id desc")
	if query != "" {
		m = m.Where("title like ?", "%"+query+"%")
	}
	find := m.Find(&sources)
	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return sources
	}

	return sources
}
func (this SourceItemService) GetSourcesItemList(sourceId, limit int) []dbmodel.SourceItem {

	ret := this.GetSourcesItemListPage(sourceId, "", 0, limit)

	return ret
}

func (this SourceItemService) FormatToXml(list []dbmodel.SourceItem) (ret model.RssRoot) {

	for _, v := range list {
		ret.Channel.Item = append(ret.Channel.Item, model.RssResultItem{
			Title:         model.CDATA{Text: v.Title},
			OriginalTitle: model.CDATA{},
			OtherTitle:    model.CDATA{},
			CnTitle:       "",
			PubDate:       v.PubDate,
			Season:        0,
			SourceId:      strconv.Itoa(v.SourceId),
			OldMinEpisode: 0,
			OldMaxEpisode: 0,
			MinEpisode:    0,
			MaxEpisode:    0,
			Enclosure: model.RssResultItemEnclosure{
				Type:   v.Enclosure.Type,
				Length: v.Enclosure.Length,
				Url:    v.Enclosure.Url,
			},
			Link: v.Link,
			Guid: model.RssResultItemGuid{
				IsPermaLink: false,
				Text:        v.Guid,
			},
			XArrRssIndexer: model.RssResultItemXArrRssIndexer{},
			OthderId:       model.RssResultItemOtherId{},
		})
	}
	return
}

func (this SourceItemService) GetSourcesItemListByHash(sourceId int, hashs ...string) []dbmodel.SourceItem {
	// 查询所有数据源信息
	sources := []dbmodel.SourceItem{}
	m := db.MainDb.Model(dbmodel.SourceItem{})
	if len(hashs) > 0 {
		if !array.InArray("-1", hashs) {
			m = m.Where("hash in ?", hashs)
		}
	}
	m = m.Where("source_id = ?", sourceId)

	find := m.Find(&sources)
	if find.Error != nil {
		logsys.Error("查询错误: %s", "GetSourcesIds", find.Error.Error())
		return sources
	}
	return sources
}

func (this SourceItemService) CreateItem(data *dbmodel.SourceItem) error {
	return db.MainDb.Model(dbmodel.SourceItem{}).Create(data).Error

}

func (this SourceItemService) UpsertItem(data dbmodel.SourceItem) *gorm.DB {
	create := db.MainDb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "source_id"}, {Name: "hash"}}, // key colume
		// 也可以用 map [ string ] interface {} { "role" : "user" }
		DoUpdates: clause.AssignmentColumns([]string{
			"title",
			"parse_info",
			"link",
			"enclosure",
		}), // column needed to be updated
	}, clause.Insert{Modifier: "or ignore"}).Create(&data)
	return create

}

// 删除小于多久的信息
func (this SourceItemService) DeleteExpireItem(sourceId, cacheDay int) error {
	// 最低缓存一天
	if cacheDay < 1 {
		cacheDay = 1
	}

	createAt := time.Now().Add(time.Duration(-cacheDay) * time.Hour * 24)
	sql := db.MainDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(dbmodel.SourceItem{}).Where("source_id = ? and created_at <= ?", sourceId, createAt.Format("2006-01-02 15:04:05")).Delete(&dbmodel.SourceItem{})
	})

	del := db.MainDb.Model(dbmodel.SourceItem{}).Where("source_id = ? and created_at <= ?", sourceId, createAt.Format("2006-01-02 15:04:05")).Delete(&dbmodel.SourceItem{})
	if del.Error == nil && del.RowsAffected > 0 {
		logsys.Debug("删除sql:%s", "清理数据缓存", sql)
		logsys.Info("清理数据源缓存 [设置缓存%d天]: %d个", "清理数据缓存", cacheDay, del.RowsAffected)

	}
	db.MainDb.Model(dbmodel.SourceItem{}).Where("pub_date = '' or title = '' or `link` = ''").Delete(&dbmodel.SourceItem{})

	// 删除空时间 空标题数据

	return del.Error
}

// 删除限制数量
func (this SourceItemService) DeleteLimitItem(sourceId, limitCount int) error {
	// 最低缓存一天
	if limitCount < 1 {
		return nil
	}
	sql := db.MainDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(dbmodel.SourceItem{}).Where("source_id=? and id  not in "+"(select id from `source_items` where source_id = ? order by id desc limit "+helper.IntToStr(limitCount)+")", sourceId, sourceId).Delete(&dbmodel.SourceItem{})
	})

	del := db.MainDb.Model(dbmodel.SourceItem{}).Where("source_id=? and id not in "+"(select id from `source_items` where source_id = ? order by id desc limit "+helper.IntToStr(limitCount)+")", sourceId, sourceId).Delete(&dbmodel.SourceItem{})
	if del.Error == nil && del.RowsAffected > 0 {
		logsys.Debug("删除sql:%s", "清理数据缓存", sql)
		logsys.Info("清理数据源缓存 [设置缓存只留下%d]: %d个", "清理数据缓存", limitCount, del.RowsAffected)

	}

	return del.Error
}

func (this SourceItemService) DeleteBySourceId(sourceId int) error {

	return db.MainDb.Model(dbmodel.SourceItem{}).Where("source_id = ?", sourceId).Delete(&dbmodel.SourceItem{}).Error
}

func (this SourceItemService) DeleteByNotInSourceId(notInSourceIds []int) *gorm.DB {
	return db.MainDb.Model(dbmodel.SourceItem{}).Where("source_id not in ?", notInSourceIds).Delete(&dbmodel.SourceItem{})
}

func (this SourceItemService) DeleteErrItem() *gorm.DB {
	//delete from source_items where source_id not in (select id from sources)
	return db.MainDb.Model(dbmodel.SourceItem{}).Where("source_id not in (select id from sources)").Delete(&dbmodel.SourceItem{})
}

func (this SourceItemService) GetSourceItemByHash(sourceId int, hashVar string) *dbmodel.SourceItem {
	rest := &dbmodel.SourceItem{}

	count := db.MainDb.Model(dbmodel.SourceItem{}).Where("source_id = ? ", sourceId).Where("hash = ?", hashVar).Find(rest)
	if count.Error != nil || count.RowsAffected <= 0 {
		return nil
	}
	return rest
}
