package monitor

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/util/date_util"
	fileUtil "XArr-Rss/util/file_util"
	"XArr-Rss/util/logsys"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// 备份数据库
func DatabaseBackup(confDir string) error {
	dbFile := confDir + "/xarr-rss.db"
	log.Println("数据库位置", dbFile)

	// 复制备份文件
	fileName := confDir + "/backup/xarr-rss-" + date_util.DateNowHourMinuteNumberStr() + ".db"
	err := fileUtil.CopyFile(dbFile, fileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("备份数据库文件失败,原因是没有找到数据库文件,可能是第一次运行")
			return nil
		}

		log.Println("备份数据库文件失败", err)
		return err
	} else {
		log.Println("数据库文件备份完成", fileName)
		return nil
	}
}

func getBackupFiles() []string {
	dir, err := os.ReadDir(variable.ConfDir + "/backup")
	if err != nil {
		logsys.Error("读取配置目录失败:%s", "清理备份文件", err.Error())
		return nil
	}
	var files []string
	for _, dirItem := range dir {
		if !dirItem.IsDir() {
			// 文件
			if strings.Contains(dirItem.Name(), "xarr-rss-") && strings.Contains(dirItem.Name(), ".db") {
				files = append(files, dirItem.Name())
			}
		}
	}
	sort.Strings(files)

	return files
}

// 清理数据库备份数量
func ClearDatabaseBackup() {
	files := getBackupFiles()
	count := appconf.AppConf.System.BackupDatabaseCount
	if count == 0 {
		return
	}
	// 备份超过20次则删除多余数据
	if len(files) > count {
		delFiles := files[:len(files)-count]
		for _, delFile := range delFiles {
			err := os.Remove(variable.ConfDir + "/backup/" + delFile)
			if err != nil {
				logsys.Error("删除文件失败:%s", "清理备份文件", err.Error())
				return
			}
		}
		logsys.Info("清理多余数据库备份成功", "清理备份文件")
	}

}

// 数据库备份监控
func DatabaseMonitor() {
	//
	for {
		// 判断是否到了间隔备份时间
		DatabaseBackupCheck()
		ClearDatabaseBackup()
		time.Sleep(3 * time.Minute)
	}
}

//
func DatabaseBackupCheck() {
	// 获取文件列表
	files := getBackupFiles()
	if len(files) == 0 {
		return
	}

	stat, err := os.Stat(variable.ConfDir + "/backup/" + files[len(files)-1])
	if err != nil {
		return
	}

	// 不开启备份
	hours := appconf.AppConf.System.BackupDatabaseTime
	if hours == 0 {
		return
	}

	// 文件修改时间少于1个小时 定期备份一次
	if stat.ModTime().Unix() < time.Now().Unix()-60*60*hours {
		err := DatabaseBackup(variable.ConfDir)
		if err != nil {
			return
		}
	}

}
