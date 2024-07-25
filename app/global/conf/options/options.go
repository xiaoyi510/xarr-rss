package options

import (
	"XArr-Rss/app/global/conf/db"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/logsys"
	"gorm.io/gorm/clause"
)

const (
	Options_HttpAddr     = "http_addr"
	Options_LogLevel     = "log_level"
	Options_LicenseKey   = "license_key"
	OptionsHTTPLitenPort = "http_listen_port"
	OptionsGlobalProxy   = "global_proxy"
	OptionsUserAgent     = "user_agent"

	Sonarr_Host        = "sonarr.host"
	Sonarr_ApiKey      = "sonarr.apikey"
	Sonarr_RefreshTime = "sonarr.refresh_time"
	Sonarr_Version     = "sonarr.version"

	Themoviedb_ApiKey    = "themoviedb.api_key"
	Themoviedb_Proxy     = "themoviedb.proxy"
	Themoviedb_ProxyType = "themoviedb.proxy_type"

	// Qb
	Qbittorrent_Host       = "qbittorrent.host"
	Qbittorrent_Cookie     = "qbittorrent.cookie"
	Qbittorrent_Username   = "qbittorrent.username"
	Qbittorrent_Password   = "qbittorrent.password"
	Qbittorrent_RenameFile = "qbittorrent.rename_file"

	// tr
	TransmissionHost     = "transmission.host"
	TransmissionUsername = "transmission.username"
	TransmissionPassword = "transmission.password"

	// 输出标题
	EchoTitleAnime = "echo_title.anime"
	EchoTitleTv    = "echo_title.tv"

	// 数据库备份数量和时间配置
	BackupDatabaseCount = "backup_database_count"
	BackupDatabaseTime  = "backup_database_time"

	SonarrUnMonitoredRemoveGroupMedia = "sonarr_unmonitored_rm_gmedia"
	SonarrUnMonitoredNAddGroupMedia   = "sonarr_unmonitored_nadd_gmedia"
	WordsRule                         = "words_rule"
)

// 初始化配置项
func InitOptions() {
	defaultValue := map[string]string{
		Options_HttpAddr:                  "http://127.0.0.1:8086",
		OptionsHTTPLitenPort:              "8086",
		Options_LogLevel:                  "info",
		Options_LicenseKey:                "",
		OptionsGlobalProxy:                "",
		Sonarr_Host:                       "http://127.0.0.1:8989/",
		Sonarr_ApiKey:                     "",
		Sonarr_RefreshTime:                "5",
		Themoviedb_ApiKey:                 "",
		Themoviedb_Proxy:                  "",
		Themoviedb_ProxyType:              "1",
		Qbittorrent_Host:                  "",
		TransmissionHost:                  "",
		TransmissionUsername:              "",
		TransmissionPassword:              "",
		Qbittorrent_Cookie:                "",
		Qbittorrent_Username:              "",
		Qbittorrent_Password:              "",
		Qbittorrent_RenameFile:            "2",
		EchoTitleAnime:                    "",
		EchoTitleTv:                       "",
		BackupDatabaseCount:               "",
		BackupDatabaseTime:                "",
		SonarrUnMonitoredRemoveGroupMedia: "",
		SonarrUnMonitoredNAddGroupMedia:   "",
		WordsRule:                         "#全局替换词\n\n\n# 某剧集替换词\n\n\n",
	}

	for k, v := range defaultValue {
		create := db.MainDb.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}}, // key colume
			// 也可以用 map [ string ] interface {} { "role" : "user" }
			DoUpdates: clause.AssignmentColumns([]string{}), // column needed to be updated
		}, clause.Insert{Modifier: "or ignore"}).Create(&dbmodel.Options{
			Key:   k,
			Value: v,
		})
		if create.Error != nil {
			logsys.Error("配置默认数据错误:%s", "InitOptions", create.Error.Error())
		}
	}

}

// 获取参数列表
func GetOptions(key ...string) map[string]string {
	op := []dbmodel.Options{}
	ret := make(map[string]string)
	// 设置默认值
	for _, k := range key {
		ret[k] = ""
	}

	first := db.MainDb.Model(dbmodel.Options{}).Where("key in ?", key).Find(&op)
	if first.Error != nil {
		return map[string]string{}
	}
	for _, v := range op {
		ret[v.Key] = v.Value
	}
	return ret
}

// 获取设置
func GetOption(key string) string {
	op := dbmodel.Options{}
	first := db.MainDb.Model(&dbmodel.Options{}).Where("key = ?", key).First(&op)
	if first.Error != nil {
		return ""
	}
	return op.Value
}

func SetOption(key, value string) error {
	// 先
	// 在 `id` 冲突时将列更新为新值
	create := db.MainDb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}}, // key colume
		// 也可以用 map [ string ] interface {} { "role" : "user" }
		DoUpdates: clause.AssignmentColumns([]string{"value"}), // column needed to be updated
	}, clause.Insert{Modifier: "or ignore"}).Create(&dbmodel.Options{
		Key:   key,
		Value: value,
	})

	//op := dbmodel.Options{
	//	Key:   key,
	//	Value: value,
	//}

	return create.Error
}
