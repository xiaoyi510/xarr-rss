package server

import (
	"XArr-Rss/app/dao/monitor"
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/service/groups"
	"XArr-Rss/server/http"
	"XArr-Rss/util/logsys"
	"bytes"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func Run(confDir, version string, public *embed.FS, menuJson []byte) {

	// 备份数据库
	monitor.DatabaseBackup(confDir)

	// 初始化配置
	InitConf(confDir, version, menuJson)
	// 开始同步数据
	//dao.MonitorDao{}.Sync()
	monitor.CenterMonitor{}.Run()

	// 运行http
	http.Run(public)
	// 检测退出新号
	go waitClose()
}

// 初始化配置
func InitConf(confDir, version string, menuJson []byte) {
	if version != "" {
		variable.ServerVersion = version
	}
	fileInfo, err := os.Stat(os.Args[0])
	if err == nil {
		variable.ProgramFileTime = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	} else {
		logsys.Error("获取文件错误:%s", "系统", err.Error())
	}
	logsys.Info("当前软件版本: %s 编译时间: %s", "系统", variable.ServerVersion, variable.ProgramFileTime)
	appconf.AppConf.ConfDir = confDir

	// 初始化菜单项
	err = json.Unmarshal(menuJson, &appconf.AppConf.Menu)
	if err != nil {
		logsys.Error("菜单信息异常:", err.Error())
	}
	_, err = os.Stat(confDir + "/debug.txt")
	if err == nil {
		variable.IsDebug = true
		logsys.IsDebug = true
	}
	variable.ConfDir = confDir

	// 连接数据库
	db.InitDb(confDir)
	options.InitOptions()

	// 升级 v1.5 配置
	//qy.ConfFileToDb()

	groups.GroupsService{}.InitGroups()
	// 载入系统配置
	appconf.ReloadConf()
	//conf_old.InitConf(confDir)

	logsys.Info("载入配置完成", "系统")

}

// 检测Ctrl+C退出
func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func waitClose() {
	//>> 监听关闭
	sigChan := make(chan os.Signal)

	//signal.Notify(sigChan)

	signal.Notify(sigChan, os.Interrupt, os.Kill)

	defer func() {
		fmt.Println("关闭服务中...")
		if variable.MainDb != nil {
			variable.MainDb.Close()
		}
		os.Exit(0)
	}()
	fmt.Println("监听任务中")
	//time.Sleep(time.Hour)
	for {
		sig := <-sigChan
		log.Println("信号:")
		log.Println(sig)
		//>> 需要判断sig信号

		break
	}
}
