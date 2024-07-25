package monitor

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/sdk/transmission"
	"XArr-Rss/app/service/download"
	"XArr-Rss/util/array"
	"XArr-Rss/util/logsys"
	"path"
	"strings"
	"time"
)

type TransmissionMonitor struct {
}

var logListTimeTr int64

// 同步 tr 下载器
func (this *TransmissionMonitor) syncTransmission() {
	// 初始化

	logsys.Info("开始监听Transmission", "Transmission")
	tr := transmission.Transmission{}.GetClient(appconf.AppConf.Service.Transmission.Host, appconf.AppConf.Service.Transmission.Username, appconf.AppConf.Service.Transmission.Password)
	for {
		tr.Host = appconf.AppConf.Service.Transmission.Host
		tr.Username = appconf.AppConf.Service.Transmission.Username
		tr.Password = appconf.AppConf.Service.Transmission.Password

		err := tr.Init()
		if err != nil {
			if err.Error() == "请输入正确的 Transmission 地址" {
			} else {
				logsys.Error("Transmission 初始化错误: %s", "Transmission", err.Error())
			}
			time.Sleep(15 * time.Second)
			// 用户未输入 不监听
			continue
		}

		if tr.Status() == nil {
			// 获取下载列表
			downloadList, _ := download.DownloadService{}.GetDownList(true)
			if len(downloadList) > 0 {
				if logListTimeTr < time.Now().Unix()-60 {
					logListTimeTr = time.Now().Unix()
					logsys.Info("开始检测总共有%d个需要检测 ", "transmission", len(downloadList))
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
					hashList := strings.Split(downloadInfo.Hash, ",")
					err, list := tr.TorrentGet(hashList)
					if err != nil {
						logsys.Error("获取种子错误:%s", "transmission", err.Error())
						continue
					}
					if len(list.Arguments.Torrents) == 0 && downloadInfo.DownClient == "tr" {
						// 已经从qbit下载过 但是任务已消失 则为已完成
						downloadInfo.Status = dbmodel.DownloadListStatusDownSuccess
						err = download.DownloadService{}.SaveDownloadInfo(&downloadInfo)
					}
					//logsys.Info("检测2:%s", "qbit", downloadInfo.Title)
					for _, item := range list.Arguments.Torrents {
						if array.InArray(item.HashString, hashList) {
							if item.HashString == item.Name {
								// 正在下载种子
								logsys.Info("对应hash正在下载种子:%s 进度:%f", "Transmission", item.HashString, item.PercentDone*100)
								continue
							}
							// 获取文件后缀 判断是文件还是目录

							ext := path.Ext(item.Name)

							if item.Name != downloadInfo.Title+ext {
								logsys.Info("找到Hash值相同数据", "Transmission", item.Name, item.HashString)
								err := tr.TorrentRenamePath(item.Id, item.Name, downloadInfo.Title+ext)
								if err != nil {
									logsys.Error("修改名字错误", "Transmission")
									continue
								}
								// 删除监听任务
								downloadInfo.Status = dbmodel.DownloadListStatusRename
								downloadInfo.Process = int(item.PercentDone * 10000)
								downloadInfo.OriginalTitle = item.Name
								if downloadInfo.Process >= 10000 {
									downloadInfo.Status = dbmodel.DownloadListStatusDownSuccess
								}
								downloadInfo.DownClient = "tr"
								err = download.DownloadService{}.SaveDownloadInfo(&downloadInfo)
								if err != nil {
									logsys.Error("移除监听任务失败:%s %s %s", "Transmission", item.Name, downloadInfo.Title+ext, downloadInfo.Hash)
								} else {
									logsys.Info("已修改名称 %s 为:%s", "Transmission", item.Name, downloadInfo.Title+ext)
								}

							} else {
								// 标题一样的删除

								// 删除监听任务
								//err = download.DownloadService{}.DelDownInfo(&downloadInfo)
								downloadInfo.Status = dbmodel.DownloadListStatusRename
								downloadInfo.Process = int(item.PercentDone * 10000)
								if downloadInfo.Process >= 10000 {
									downloadInfo.Status = dbmodel.DownloadListStatusDownSuccess
								}
								downloadInfo.DownClient = "tr"

								err = download.DownloadService{}.SaveDownloadInfo(&downloadInfo)

								if err != nil {
									logsys.Error("移除监听任务失败:%s %s %s", "Transmission", item.Name, downloadInfo.Title+ext, downloadInfo.Hash)
								}
							}
						}
					}

				}

			}

		} else {
			// Transmission 异常
			logsys.Error("Transmission服务异常", "Transmission")

			time.Sleep(1 * time.Minute)
			continue
		}
		time.Sleep(12 * time.Second)
	}

}
