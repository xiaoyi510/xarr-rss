package appconf

import (
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/util/helper"
	"strconv"
)

type appConfSystem struct {
	WordsRule                   string `json:"words_rule"`
	HttpAddr                    string `json:"http_addr"`
	LogLevel                    string `json:"log_level"`
	LicenseKey                  string `json:"license_key"`
	Port                        int    `json:"http_listen_port"`
	GlobalProxy                 string `json:"global_proxy"`
	UserAgent                   string `json:"user_agent"`
	EchoTitleAnime              string `json:"echo_title_anime"`
	EchoTitleTv                 string `json:"echo_title_tv"`
	BackupDatabaseCount         int    `json:"backup_database_count"`
	BackupDatabaseTime          int64  `json:"backup_database_time"`
	SonarrUnmonitoredRmGmedia   string `json:"sonarr_unmonitored_rm_gmedia"`
	SonarrUnmonitoredNaddGmedia string `json:"sonarr_unmonitored_nadd_gmedia"`
}

// 加载系统配置
func (receiver *appConfSystem) Reload() {
	// 读取配置
	configs := options.GetOptions(
		options.Options_HttpAddr,
		options.Options_LogLevel,
		options.Options_LicenseKey,
		options.OptionsHTTPLitenPort,
		options.OptionsGlobalProxy,
		options.OptionsUserAgent,
		options.EchoTitleAnime,
		options.EchoTitleTv,
		options.BackupDatabaseCount,
		options.BackupDatabaseTime,
		options.SonarrUnMonitoredNAddGroupMedia,
		options.SonarrUnMonitoredRemoveGroupMedia,
		options.WordsRule,
	)
	// 修改配置
	AppConf.System.Port = helper.StrToInt(configs[options.OptionsHTTPLitenPort])
	if AppConf.System.Port <= 0 {
		AppConf.System.Port = 8086
	}
	AppConf.System.HttpAddr = configs[options.Options_HttpAddr]
	AppConf.System.WordsRule = configs[options.WordsRule]
	AppConf.System.LogLevel = configs[options.Options_LogLevel]
	AppConf.System.LicenseKey = configs[options.Options_LicenseKey]
	AppConf.System.GlobalProxy = configs[options.OptionsGlobalProxy]
	AppConf.System.UserAgent = configs[options.OptionsUserAgent]
	AppConf.System.EchoTitleAnime = configs[options.EchoTitleAnime]
	AppConf.System.EchoTitleTv = configs[options.EchoTitleTv]
	AppConf.System.SonarrUnmonitoredNaddGmedia = configs[options.SonarrUnMonitoredNAddGroupMedia]
	AppConf.System.SonarrUnmonitoredRmGmedia = configs[options.SonarrUnMonitoredRemoveGroupMedia]
	AppConf.System.BackupDatabaseTime = helper.StrToInt64(configs[options.BackupDatabaseTime])
	AppConf.System.BackupDatabaseCount = helper.StrToInt(configs[options.BackupDatabaseCount])

}

// 保存系统配置
func (receiver *appConfSystem) Save() error {
	options.SetOption(options.OptionsHTTPLitenPort, strconv.Itoa(AppConf.System.Port))
	options.SetOption(options.BackupDatabaseTime, strconv.Itoa(int(AppConf.System.BackupDatabaseTime)))
	options.SetOption(options.BackupDatabaseCount, strconv.Itoa(AppConf.System.BackupDatabaseCount))
	options.SetOption(options.Options_HttpAddr, AppConf.System.HttpAddr)
	options.SetOption(options.Options_LogLevel, AppConf.System.LogLevel)
	options.SetOption(options.Options_LicenseKey, AppConf.System.LicenseKey)
	options.SetOption(options.OptionsGlobalProxy, AppConf.System.GlobalProxy)
	options.SetOption(options.OptionsUserAgent, AppConf.System.UserAgent)
	options.SetOption(options.EchoTitleAnime, AppConf.System.EchoTitleAnime)
	options.SetOption(options.EchoTitleTv, AppConf.System.EchoTitleTv)
	options.SetOption(options.WordsRule, AppConf.System.WordsRule)
	options.SetOption(options.SonarrUnMonitoredNAddGroupMedia, AppConf.System.SonarrUnmonitoredNaddGmedia)
	options.SetOption(options.SonarrUnMonitoredRemoveGroupMedia, AppConf.System.SonarrUnmonitoredRmGmedia)
	return nil
}
