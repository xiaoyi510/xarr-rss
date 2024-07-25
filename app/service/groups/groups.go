package groups

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/match"
	"XArr-Rss/util/golimit"
	"XArr-Rss/util/hash"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/xml"
	"fmt"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type GroupsService struct {
}

// 获取分组下的媒体数据
func (this GroupsService) GetGroupList(ids ...int) []dbmodel.Group {
	group := []dbmodel.Group{}
	find := db.MainDb.Model(dbmodel.Group{})
	if len(ids) > 0 {
		if len(ids) == 1 {
			if ids[0] > 0 {
				find.Where("id = ?", ids[0])
			}
		} else {
			find.Where("id in ?", ids)
		}
	}
	find.Find(&group)
	if find.Error != nil {
		logsys.Error("查询分组数据错误:%s", "GetGroupList", find.Error.Error())

	}
	return group
}

// 从数据源同步分组匹配数据
func (this GroupsService) SyncGroupItem(groupInfo dbmodel.Group, useCache bool) (error, *model.RssRoot) {
	groupsCache := &model.RssRoot{Version: "2.0"}
	groupsCache.Channel.Link = appconf.AppConf.System.HttpAddr + "/rss/group_" + strconv.Itoa(int(groupInfo.Id)) + ".xml"
	groupsCache.Channel.Description = "XArr-Rss 分组订阅 " + groupInfo.Name
	groupsCache.Channel.Title = "XArr-Rss-" + groupInfo.Name

	// 生成等待组
	var synMedia = golimit.NewGoLimit(2)
	// 查询group下面的媒体列表
	groupMedias := GroupMediaService{}.GetGroupMedias(groupInfo.Id)
	for _, groupMedia := range groupMedias {
		synMedia.Add()
		// 解析分组数据
		go func(groupMedia *dbmodel.GroupMedia) {
			groupMediaParse := match.ParseGroupMediaInfo(groupMedia, useCache)
			if groupMediaParse != nil {

				// 归纳数据对应数据源数据
				var sourceItemsGroup = make(map[string][]model.RssResultItem)
				for _, v := range groupMediaParse.Channel.Item {
					sourceItemsGroup[v.SourceId] = append(sourceItemsGroup[v.SourceId], v)
				}
				for sourceId, sourceItems := range sourceItemsGroup {
					groupsMediaCache2 := &model.RssRoot{Version: "2.0"}
					groupsMediaCache2.Channel.Link = appconf.AppConf.System.HttpAddr + "/rss/group_" + strconv.Itoa(int(groupMedia.GroupId)) + "/" + strconv.Itoa(groupMedia.Id) + "/" + sourceId + ".xml"
					groupsMediaCache2.Channel.Description = "XArr-Rss 分组媒体数据源订阅 " + groupMedia.MediaInfo.CnTitle
					groupsMediaCache2.Channel.Title = "XArr-Rss-" + groupMedia.MediaInfo.CnTitle
					groupsMediaCache2.Channel.Item = sourceItems
					err := this.SaveGroupMediaSourceXml(groupsMediaCache2, strconv.Itoa(groupMedia.GroupId), strconv.Itoa(groupMedia.Id), sourceId)
					if err != nil {
						logsys.Error("保存数据源匹配数据错误:%s", "保存数据源", err.Error())
					} else {
						logsys.Info("保存数据源匹配数据完成-分组/媒体/数据源:%s", "保存数据源", "/rss/group_"+strconv.Itoa(int(groupMedia.GroupId))+"/"+strconv.Itoa(groupMedia.Id)+"/"+sourceId+".xml")
					}
				}
				/////////////////////////////////

				groupsCache.Channel.Item = append(groupsCache.Channel.Item, groupMediaParse.Channel.Item...)
				// 保存到分组媒体单独的数据
				err := this.SaveGroupMediaXml(groupMediaParse, groupInfo.Id, groupMedia.Id)
				if err != nil {
					logsys.Error("保存媒体匹配结果失败:%s", "保存数据源", err.Error())
				} else {
					logsys.Info("保存数据源匹配数据完成-分组/媒体:%s", "保存数据源", "/rss/group_"+strconv.Itoa(int(groupMedia.GroupId))+"/"+strconv.Itoa(groupMedia.Id)+".xml")
				}

			}

			synMedia.Done()
		}(groupMedia)
	}
	synMedia.Wait()

	// 保存分组信息
	err := this.SaveGroupXml(groupsCache, strconv.Itoa(groupInfo.Id))
	if err != nil {
		return err, nil
	} else {
		logsys.Info("保存数据源匹配数据完成-分组:%s", "保存数据源", "/rss/group_"+strconv.Itoa(int(groupInfo.Id))+".xml")
	}

	return nil, groupsCache

}

func (this GroupsService) SaveGroupXml(groupsCache *model.RssRoot, groupId string) error {

	// 设置默认分组数据
	if len(groupsCache.Channel.Item) == 0 {
		groupsCache.Channel.Item = append(groupsCache.Channel.Item, model.RssResultItem{
			Title: model.CDATA{
				Text: "XArr-Rss 默认占位使用",
			},
			PubDate: "2022-04-21T05:38:00",
			Enclosure: model.RssResultItemEnclosure{
				Type:   "application/x-bittorrent",
				Length: "329672288",
				Url:    "https://xarr.52nyg.com",
			},
			Link: "https://xarr.52nyg.com",
			Guid: model.RssResultItemGuid{
				IsPermaLink: false,
				Text:        "https://xarr.52nyg.com",
			},
			XArrRssIndexer: model.RssResultItemXArrRssIndexer{
				Text: groupId,
				ID:   groupId,
			},
		})
	}

	// 写出到group 翻译后的数据
	marshal, err := xml.MarshalIndent(groupsCache, "", "")
	if err != nil {
		return logsys.Error("分组数据转换失败:%s", "同步分组数据", err.Error())
	}

	err = os.WriteFile(appconf.AppConf.ConfDir+"/trans/group_"+groupId+".xml", marshal, 0666)
	if err != nil {
		return logsys.Error("写出分组数据失败:%s", "同步分组数据", err.Error())
	}

	return nil

}

func (this GroupsService) SaveGroupMediaSourceXml(groupsCache *model.RssRoot, groupId, mediaId, sourceId string) error {
	if groupId == "" || mediaId == "" || sourceId == "" {
		return logsys.Error("没有完整的信息进行保存缓存", "保存分组媒体数据")
	}
	// 设置默认分组数据
	if len(groupsCache.Channel.Item) == 0 {
		groupsCache.Channel.Item = append(groupsCache.Channel.Item, model.RssResultItem{
			Title: model.CDATA{
				Text: "XArr-Rss 默认占位使用",
			},
			PubDate: "2022-04-21T05:38:00",
			Enclosure: model.RssResultItemEnclosure{
				Type:   "application/x-bittorrent",
				Length: "329672288",
				Url:    "https://xarr.52nyg.com",
			},
			Link: "https://xarr.52nyg.com",
			Guid: model.RssResultItemGuid{
				IsPermaLink: false,
				Text:        "https://xarr.52nyg.com",
			},
			XArrRssIndexer: model.RssResultItemXArrRssIndexer{
				Text: groupId,
				ID:   groupId,
			},
		})
	}
	// 写出到group 翻译后的数据
	marshal, err := xml.MarshalIndent(groupsCache, "", "")
	if err != nil {
		return logsys.Error("分组数据转换失败:%s", "保存分组媒体数据", err.Error())
	}

	dirName := fmt.Sprintf("/trans/group_%s/%s", groupId, mediaId)
	// 创建目录
	_, err = os.Stat(appconf.AppConf.ConfDir + dirName)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(appconf.AppConf.ConfDir+dirName, 0777)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = os.WriteFile(appconf.AppConf.ConfDir+dirName+"/"+sourceId+".xml", marshal, 0750)
	if err != nil {
		return logsys.Error("写出分组数据失败:%s", "同步分组数据", err.Error())
	}
	return nil
}

func (this GroupsService) SaveGroupMediaXml(groupsCache *model.RssRoot, groupId, mediaId int) error {
	if groupId <= 0 || mediaId <= 0 {
		return logsys.Error("没有完整的信息进行保存缓存", "保存分组媒体数据")
	}
	// 设置默认分组数据
	if len(groupsCache.Channel.Item) == 0 {
		groupsCache.Channel.Item = append(groupsCache.Channel.Item, model.RssResultItem{
			Title: model.CDATA{
				Text: "XArr-Rss 默认占位使用",
			},
			PubDate: "2022-04-21T05:38:00",
			Enclosure: model.RssResultItemEnclosure{
				Type:   "application/x-bittorrent",
				Length: "329672288",
				Url:    "https://xarr.52nyg.com",
			},
			Link: "https://xarr.52nyg.com",
			Guid: model.RssResultItemGuid{
				IsPermaLink: false,
				Text:        "https://xarr.52nyg.com",
			},
			XArrRssIndexer: model.RssResultItemXArrRssIndexer{
				Text: strconv.Itoa(groupId),
				ID:   strconv.Itoa(groupId),
			},
		})
	}
	// 写出到group 翻译后的数据
	marshal, err := xml.MarshalIndent(groupsCache, "", "")
	if err != nil {
		return logsys.Error("分组数据转换失败:%s", "保存分组媒体数据", err.Error())
	}

	dirName := fmt.Sprintf("/trans/group_%d", groupId)
	// 创建目录
	_, err = os.Stat(appconf.AppConf.ConfDir + dirName)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(appconf.AppConf.ConfDir+dirName, 0777)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = os.WriteFile(appconf.AppConf.ConfDir+dirName+"/"+strconv.Itoa(mediaId)+".xml", marshal, 0777)
	if err != nil {
		return logsys.Error("写出分组数据失败:%s", "同步分组数据", err.Error())
	}
	return nil
}

// 判断分组是否存在
func (this GroupsService) GroupExits(groupId int) bool {
	var count int64

	db.MainDb.Model(dbmodel.Group{}).Where("id = ?", groupId).Count(&count)
	return count > 0
}

// 获取分组信息
func (this GroupsService) GroupInfo(groupId int) (*dbmodel.Group, error) {
	var info dbmodel.Group
	//Model: dbmodel.Model{Id: groupId},
	first := db.MainDb.Model(dbmodel.Group{}).Where("id = ?", groupId).First(&info)
	return &info, first.Error
}

// 删除分组
func (this GroupsService) GroupDelete(groupId int) error {

	rest := db.MainDb.Model(dbmodel.Group{}).Where("id = ?", groupId).Delete(&dbmodel.Group{})

	return rest.Error
}

// 将所有的group_xx.xml 组合为groups
func (this GroupsService) SyncGroups() *model.RssRoot {
	groupAll := &model.RssRoot{Version: "2.0"}
	groupAll.Channel.Link = appconf.AppConf.System.HttpAddr + "/rss/group/group_all.xml"
	groupAll.Channel.Description = "XArr-Rss"
	groupAll.Channel.Title = "XArr-Rss"
	groupAll.Channel.Item = append(groupAll.Channel.Item, model.RssResultItem{
		Title: model.CDATA{
			Text: "XArr-Rss 默认占位使用",
		},
		PubDate: "2022-04-21T05:38:00",
		Enclosure: model.RssResultItemEnclosure{
			Type:   "application/x-bittorrent",
			Length: "996973766",
			Url:    "https://xarr.52nyg.com",
		},
		Link: "https://xarr.52nyg.com",
		Guid: model.RssResultItemGuid{
			IsPermaLink: false,
			Text:        hash.Md5{}.HashString("https://xarr.52nyg.com" + time.Now().GoString()),
		},
	})

	// 枚举目录下面的所有group_xx.xml
	files, _ := ioutil.ReadDir("./conf/trans")
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		// 找到groupId
		name := f.Name()
		groupId := strings.Replace(name, "group_", "", -1)
		groupId = strings.Replace(groupId, ".xml", "", -1)
		if groupId == "all" {
			continue
		}

		if !this.GroupExits(helper.StrToInt(groupId)) {
			logsys.Error("分组可能被删除,删除分组文件", "同步分组")
			os.Remove("./conf/trans/" + name)
			continue
		}

		// 读取文件内容
		result, err := os.ReadFile("./conf/trans/" + name)
		if err != nil {
			logsys.Error("读取分组信息失败:%s", "同步分组", err.Error())
			continue
		}

		groupCache := &model.RssRoot{}
		err = xml.Unmarshal(result, groupCache)
		if err != nil {
			logsys.Error("解析分组信息失败:%s", "同步分组", err.Error())
			continue
		}
		if groupCache != nil && len(groupCache.Channel.Item) > 0 {
			for _, v2 := range groupCache.Channel.Item {
				var has = false
				// 检测是否已拥有
				for _, ss := range groupAll.Channel.Item {
					if ss.Guid.Text == v2.Guid.Text {
						has = true
						continue
					}
				}

				if v2.Title.Text != "XArr-Rss 默认占位使用" && has != true {
					groupAll.Channel.Item = append(groupAll.Channel.Item, v2)
				}
			}
		}

	}

	// 写出到group 翻译后的数据
	marshal2, err := xml.MarshalIndent(groupAll, "", "")
	if err != nil {
		logsys.Error("分组数据转换失败:%s", "同步分组数据", err.Error())
		return nil
	}
	err = os.WriteFile(appconf.AppConf.ConfDir+"/trans/group_all.xml", marshal2, 0666)
	if err != nil {
		logsys.Error("写出分组数据失败:%s", "同步分组数据", err.Error())
		return nil
	}
	//logsys.Info("分组数据同步完成", "同步分组数据")

	return groupAll
}

func (this GroupsService) GetGroupCount() int64 {
	var count int64

	db.MainDb.Model(dbmodel.Group{}).Count(&count)
	return count
}

func (this GroupsService) CreateGroup(data *dbmodel.Group) error {
	create := db.MainDb.Model(dbmodel.Group{}).Create(data)
	return create.Error
}

func (this GroupsService) Save(group *dbmodel.Group) error {
	save := db.MainDb.Save(group)
	return save.Error
}

// 初始化分组
func (this GroupsService) InitGroups() {
	var count int64
	db.MainDb.Model(dbmodel.Group{}).Count(&count)
	if count == 0 {

		//db.MainDb.Clauses(clause.Insert{Modifier: "or ignore"})

		create := db.MainDb.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "name"}}, // key colume
			// 也可以用 map [ string ] interface {} { "role" : "user" }
			DoUpdates: clause.AssignmentColumns([]string{}), // column needed to be updated
		}, clause.Insert{Modifier: "or ignore"}).Create(&dbmodel.Group{
			Name:               "默认分组",
			AutoInsertSonarr:   dbmodel.AutoInsertSonarrNone,
			LastInsertSonarrId: 0,
		})
		if create.Error != nil {
			logsys.Error("配置默认数据错误:%s", "InitGroups", create.Error.Error())
		}
	}

}

func (this GroupsService) GetTotalCount() int {
	var ret int64
	db.MainDb.Model(&dbmodel.Group{}).Count(&ret)
	return int(ret)
}
