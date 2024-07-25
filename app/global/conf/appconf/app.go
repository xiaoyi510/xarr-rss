package appconf

import "XArr-Rss/app/model/webmenu"

var AppConf = &appConf{}

// 系统配置中心
type appConf struct {
	ConfDir string
	System  appConfSystem
	Service appConfService
	Menu    []*webmenu.WebMenu
}

// 载入系统配置
func ReloadConf() {
	AppConf.System.Reload()
	AppConf.Service.Reload()
}
