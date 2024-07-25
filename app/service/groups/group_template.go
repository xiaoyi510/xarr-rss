package groups

import (
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/json_util"
	"XArr-Rss/util/logsys"
	"database/sql"
	"errors"
)

type GroupTemplateService struct {
}

func (this GroupTemplateService) GetList() []*dbmodel.GroupTemplate {
	items := []*dbmodel.GroupTemplate{}

	find := db.MainDb.Model(&dbmodel.GroupTemplate{}).Order("id desc").Find(&items)
	if find.Error != nil {
		if find.Error != sql.ErrNoRows {
			logsys.Error("查询分组模板错误:%s", find.Error.Error())
		}
		return nil
	}
	return items
}

func (this GroupTemplateService) FindOne(id int32) *dbmodel.GroupTemplate {
	var item *dbmodel.GroupTemplate

	find := db.MainDb.Model(dbmodel.GroupTemplate{}).Where("id = ?", id).First(&item)
	if find.Error != nil {
		if find.Error != sql.ErrNoRows {
			logsys.Error("查询分组模板错误:%s", find.Error.Error())
		}
		return nil
	}
	return item
}
func (this GroupTemplateService) FindALl(id ...int32) map[int]*dbmodel.GroupTemplate {
	ret := make(map[int]*dbmodel.GroupTemplate)
	items := []*dbmodel.GroupTemplate{}

	find := db.MainDb.Model(dbmodel.GroupTemplate{}).Where("id in ?", id).Find(&items)
	if find.Error != nil {
		if find.Error != sql.ErrNoRows {
			logsys.Error("查询分组模板错误:%s", find.Error.Error())
		}
		return nil
	}
	for _, ite := range items {
		ret[ite.Id] = ite
	}

	return ret
}

func (this GroupTemplateService) Delete(id int32) error {

	rest := db.MainDb.Model(dbmodel.GroupTemplate{}).Where("id = ?", id).Delete(&dbmodel.GroupTemplate{})

	return rest.Error
}

func (this GroupTemplateService) Add(item *dbmodel.GroupTemplate) error {
	create := db.MainDb.Model(dbmodel.GroupTemplate{}).Create(item)
	return create.Error
}

func (this GroupTemplateService) Edit(item *dbmodel.GroupTemplate) error {
	create := db.MainDb.Save(item)
	return create.Error
}

func (this GroupTemplateService) BatchUse(id int32, groupMediaIds []int32) error {
	info := this.FindOne(id)
	if info == nil {
		return errors.New("没有找到模板信息")
	}

	setInfo := db.MainDb.Model(&dbmodel.GroupMedia{}).Where("id in (?)", groupMediaIds).Updates(map[string]interface{}{
		"language":          info.Language,
		"quality":           info.Quality,
		"use_source":        json_util.JsonMarshalToString(info.UseSource),
		"regex":             json_util.JsonMarshalToString(info.Regex),
		"filter_push_group": json_util.JsonMarshalToString(info.FilterPushGroup),
		"echo_title_anime":  info.EchoTitleAnime,
		"echo_title_tv":     info.EchoTitleTv,
	})
	return setInfo.Error
}
