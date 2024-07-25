package qy

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/util/logsys"
	"os"
)

func ConfFileToDb() {
	// 判断confv2.json 是否存在
	_, err := os.Stat(appconf.AppConf.ConfDir + "/confv2.json")
	if os.IsNotExist(err) {
		//logsys.Error("未找到配置文件", "迁移配置")
		return
	}
	logsys.Error("找到配置文件,请手动更新到1.5.x进行数据迁移 2.x以上版本已不支持1.4.x版本进行迁移", "迁移配置")
	return

}
