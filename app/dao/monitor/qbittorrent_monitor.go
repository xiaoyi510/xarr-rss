package monitor

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/model/qbit/torrents"
	"XArr-Rss/app/sdk/qbit"
	"XArr-Rss/app/sdk/qbit/client"
	"XArr-Rss/app/service/download"
	"XArr-Rss/util/array"
	"XArr-Rss/util/logsys"
	"path/filepath"
	"strings"
	"time"
)

type QbittorrentMonitor struct {
}

var logListTime int64

// 监听Qbit下载情况
func (this *QbittorrentMonitor) SyncQbit() {
	logsys.Info("开始监听Qbit", "qbit")

	qbit.Qbit.Client = &client.QbitClient{}
	qbit.Qbit.Client.Init(appconf.AppConf.Service.Qbittorrent.Host)
	for {
		if appconf.AppConf.Service.Qbittorrent.Host == "" {
			time.Sleep(15 * time.Second)
			continue
		}
		// 解决qbit客户端保存未生效问题
		qbit.Qbit.Client.SetConf(appconf.AppConf.Service.Qbittorrent.Host)

		if checkQbit() {
			//logsys.Debug("检测登录", "qbit")
			if qbit.Qbit.CheckLogin(appconf.AppConf.Service.Qbittorrent.Username, appconf.AppConf.Service.Qbittorrent.Password) {
				downloadList, _ := download.DownloadService{}.GetDownList(true)
				if len(downloadList) > 0 && logListTime < time.Now().Unix()-60 {
					logListTime = time.Now().Unix()
					logsys.Info("开始检测总共有%d个需要检测 ", "qbit", len(downloadList))
				}
				// 登录成功
				for _, downloadInfo := range downloadList {
					// 判断种子是否超时1小时没客户端下载 就设置为成功
					if downloadInfo.CreatedAt.Unix() <= time.Now().Unix()-60*60 && downloadInfo.Status == dbmodel.DownloadListStatusWait && downloadInfo.DownClient == "" {
						downloadInfo.Status = dbmodel.DownloadListStatusDownSuccess
						download.DownloadService{}.SaveDownloadInfo(&downloadInfo)
						continue
					}

					// 超过七天没完成的任务设置已完成
					if downloadInfo.CreatedAt.Unix() <= time.Now().Unix()-60*60*24*7 && downloadInfo.Status == dbmodel.DownloadListStatusWait {
						downloadInfo.Status = dbmodel.DownloadListStatusDownSuccess
						download.DownloadService{}.SaveDownloadInfo(&downloadInfo)
						continue
					}

					// 查询hash 对应的种子信息
					hashList := strings.Split(strings.ToLower(downloadInfo.Hash), ",")
					err, list := qbit.Qbit.GetTorrents().Info(hashList, torrents.ApiTorrentsInfoReq{})
					if err != nil {
						logsys.Error("获取种子[%s]错误:%s", "qbit", downloadInfo.Hash, err.Error())
						continue
					}
					//log.Println("测试种子", len(*list), *list, downloadInfo.DownClient, downloadInfo.CreatedAt, downloadInfo.Hash)
					if len(*list) == 0 && downloadInfo.DownClient == "qbit" && downloadInfo.CreatedAt.Unix() < time.Now().Unix()-60*2 {
						// 已经从qbit下载过 但是任务已消失 则为已完成
						downloadInfo.Status = dbmodel.DownloadListStatusDownSuccess
						err = download.DownloadService{}.SaveDownloadInfo(&downloadInfo)
						logsys.Error("种子[%s]未能从下载器中获取信息", "qbit", downloadInfo.Hash)

						continue
					}

					//logsys.Info("检测2:%s", "qbit", downloadInfo.Title)
					for _, item := range *list {
						if array.InArray(strings.ToLower(item.Hash), hashList) {
							if item.Hash == item.Name {
								// 正在下载种子
								logsys.Info("对应hash正在下载种子:%s 进度:%f", "qbit", item.Hash, item.Progress)
								continue
							}

							if item.Name != downloadInfo.Title {
								logsys.Info("找到Hash值相同数据,但未修改标题", "qbit", item.Name, item.Hash)
								err, _ := qbit.Qbit.GetTorrents().Rename(torrents.ApiTorrentRenameReq{
									Hash: item.Hash,
									Name: downloadInfo.Title,
								})

								if err != nil {
									logsys.Error("修改名字错误", "qbit")
									continue
								}
								// 删除监听任务
								downloadInfo.Status = dbmodel.DownloadListStatusRename
								downloadInfo.Process = int(item.Progress * 10000)
								downloadInfo.OriginalTitle = item.Name
								downloadInfo.DownClient = "qbit"

								err = download.DownloadService{}.SaveDownloadInfo(&downloadInfo)
								if err != nil {
									logsys.Error("移除监听任务失败:%s %s %s", "qbit", item.Name, downloadInfo.Title, downloadInfo.Hash)
								} else {
									logsys.Info("已修改名称 %s 为:%s", "qbit", item.Name, downloadInfo.Title)
								}

								if appconf.AppConf.Service.Qbittorrent.RenameFile == "1" {
									// 需要处理文件名
									logsys.Info("准备修改下载文件的文件名", "qbit")

									_, fileList := qbit.Qbit.GetTorrents().Files(torrents.ApiTorrentFilesReq{
										Hash: item.Hash,
									})
									if fileList != nil && len(*fileList) == 1 {
										logsys.Info("匹配到只有一个文件,可以修改文件名", "qbit")

										// 获取文件名后缀
										fileExt := filepath.Ext((*fileList)[0].Name)

										err, _ := qbit.Qbit.GetTorrents().RenameFile(torrents.ApiTorrentRenameFileReq{
											Hash:    item.Hash,
											OldPath: (*fileList)[0].Name,
											NewPath: strings.ReplaceAll(strings.ReplaceAll((*fileList)[0].Name, item.Name, downloadInfo.Title), fileExt, "") + fileExt,
										})
										if err != nil {
											logsys.Error("修改文件名失败:%s", "qbit", err.Error())
										}
									} else {
										logsys.Error("匹配文件异常,不修改文件名", "qbit")

									}

								}

							} else {
								// 标题一样的删除

								// 删除监听任务
								//err = download.DownloadService{}.DelDownInfo(&downloadInfo)
								downloadInfo.Status = dbmodel.DownloadListStatusRename
								downloadInfo.Process = int(item.Progress * 10000)
								downloadInfo.DownClient = "qbit"
								if downloadInfo.Process >= 10000 {
									downloadInfo.Status = dbmodel.DownloadListStatusDownSuccess
								}

								err = download.DownloadService{}.SaveDownloadInfo(&downloadInfo)

								if err != nil {
									logsys.Error("移除监听任务失败:%s %s %s", "qbit", item.Name, downloadInfo.Title, downloadInfo.Hash)
								}
							}
						}
					}

				}
			} else {
				logsys.Debug("Qbit登录异常", "qbit")

			}
		} else {
			//logsys.Debug("检测Qbit出错", "qbit")
			time.Sleep(1 * time.Minute)
		}
		//logsys.Debug("等待12秒继续运行", "qbit")
		time.Sleep(12 * time.Second)
	}

}

func checkQbit() bool {
	// 判断Qb是否输入了Key
	if appconf.AppConf.Service.Qbittorrent.Host == "" {
		logsys.Error("请输入Qbit环境信息", "qbit")
		return false
	}
	//if appconf.AppConf.Service.Qbittorrent.Username == "" {
	//	logsys.Error("请输入Qbit环境信息", "qbit")
	//	return false
	//
	//}
	//if appconf.AppConf.Service.Qbittorrent.Password == "" {
	//	logsys.Error("请输入Qbit环境信息", "qbit")
	//	return false
	//}
	return true

}
