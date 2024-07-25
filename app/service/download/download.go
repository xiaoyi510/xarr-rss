package download

import (
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/logsys"
	"regexp"
	"strings"
	"time"
)

type DownloadService struct {
}

func (this DownloadService) GetDownInfo(hash string) (ret *dbmodel.DownloadList) {
	ret = &dbmodel.DownloadList{}
	first := db.MainDb.Model(&dbmodel.DownloadList{}).Where("hash = ?", hash).First(ret)
	if first.Error != nil {
		return nil
	}
	return ret
}
func (this DownloadService) GetDownList(notDownSuccess bool) ([]dbmodel.DownloadList, error) {
	var ret []dbmodel.DownloadList
	if notDownSuccess {
		first := db.MainDb.Model(&dbmodel.DownloadList{}).Where("status in ?", []int{dbmodel.DownloadListStatusWait, dbmodel.DownloadListStatusRename}).Find(&ret)
		return ret, first.Error
	}
	first := db.MainDb.Model(&dbmodel.DownloadList{}).Find(&ret)
	return ret, first.Error
}

func (this DownloadService) DelDownInfo(data *dbmodel.DownloadList) error {
	return db.MainDb.Model(&dbmodel.DownloadList{}).Delete(data).Error
}

func (this DownloadService) AddDownloadInfo(data *dbmodel.DownloadList) {
	create := db.MainDb.Create(data)
	if create.Error != nil {
		logsys.Error("保存监控任务失败:%s", "AddDownloadInfo", create.Error.Error())
	} else {
		logsys.Info("创建监听任务成功", "AddDownloadInfo")
	}
}

func (this DownloadService) SaveDownloadInfo(data *dbmodel.DownloadList) error {
	create := db.MainDb.Save(data)
	if create.Error != nil {
		logsys.Error("保存监控任务失败:%s", "SaveDownloadInfo", create.Error.Error())
		return create.Error
	} else {
		//logsys.Info("保存监听任务成功", "SaveDownloadInfo")
		return nil
	}
}

// 获取今日新增量
func (this DownloadService) GetTodayNew() int {
	var ret int64
	db.MainDb.Model(&dbmodel.DownloadList{}).Where("created_at > datetime(?)", time.Now()).Count(&ret)
	return int(ret)
}

// 获取今日新增且种子已完成量
func (this DownloadService) GetTodayNewDownSuccess() int {
	var ret int64
	db.MainDb.Model(&dbmodel.DownloadList{}).Where("created_at > datetime(?) and status = ?", time.Now(), dbmodel.DownloadListStatusDownSuccess).Count(&ret)
	return int(ret)
}

// 获取总监控数量
func (this DownloadService) GetTotalCount() int {
	var ret int64
	db.MainDb.Model(&dbmodel.DownloadList{}).Count(&ret)
	return int(ret)
}

func (this DownloadService) GetUrlTorrentHash(uri string) string {
	reg := regexp.MustCompile(`(.?)([a-z\dA-Z]{40})(.?)`)
	a := reg.FindAllStringSubmatch(uri, -1)
	for _, v := range a {
		str2 := "abcdefghijklmnopqrstuvwxyz0123456789"
		if (!strings.Contains(str2, strings.ToLower(v[1])) || v[1] == "") && (!strings.Contains(str2, strings.ToLower(v[3])) || v[3] == "") {
			return v[2]
		}
	}
	return ""
}
