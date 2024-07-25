package main

import (
	"XArr-Rss/server"
	"XArr-Rss/util/logsys"
	"embed"
	"os"
)

var Version string

func main() {

	logsys.Info("软件定制请联系QQ: 319555754 QQ群:996973766 ", "系统")
	// 初始化配置文件
	resourceInit()
	// 初始化配置
	server.Run("./conf", Version, &publicDir, menuJson)
}

//go:embed web/data/menu.json
var menuJson []byte

//go:embed web
var publicDir embed.FS

func resourceInit() bool {

	createDir("./conf")
	createDir("./conf/cache")
	createDir("./conf/backup")
	createDir("./conf/trans")
	createDir("./conf/torrents")
	createDir("./conf/data")
	createDir("./images")

	return true
}

func createDir(dirname string) {

	// 判断目录是否存在
	dirInfo, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		// 创建目录
		err := os.Mkdir(dirname, 0777)
		if err != nil {
			logsys.Panicln("创建目录失败,请检查权限", "系统")
		}
		return
	} else if err != nil {
		panic(err)
	}

	if dirInfo.IsDir() == false {
		if err != nil {
			logsys.Panicln(dirname+"  已有文件存在,请检查", "系统")

		}
	}
}
