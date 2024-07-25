package appconf

import (
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/util/helper"
	"strconv"
)

// TMDB 代理获取
const TMDB_PROXY_TYPE_NONE = 1   // 不使用
const TMDB_PROXY_TYPE_GLOBAL = 2 // 使用全局代理
const TMDB_PROXY_TYPE_DIY = 3    // 使用自定义代理

type appConfService struct {
	Sonarr struct {
		Host        string `json:"host"`
		Apikey      string `json:"apikey"`
		Version     string `json:"version"`
		RefreshTime int    `json:"refresh_time"`
	} `json:"sonarr"`
	Themoviedb struct {
		Apikey    string `json:"apikey"`
		Proxy     string `json:"proxy"`
		ProxyType int    `json:"proxy_type"`
	} `json:"themoviedb"`
	Transmission ServiceTransmission `json:"transmission"`
	Qbittorrent  ServiceQbittorrent  `json:"qbittorrent"`
}

type ServiceTransmission struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type ServiceQbittorrent struct {
	Host       string `json:"host"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	RenameFile string `json:"rename_file"`
}

// 加载系统配置
func (this *appConfService) Reload() {
	// 读取配置
	configs := options.GetOptions(options.Sonarr_Host,
		options.Sonarr_ApiKey, options.Sonarr_RefreshTime, options.Sonarr_Version,
		options.Themoviedb_ApiKey, options.Themoviedb_Proxy, options.Themoviedb_ProxyType,
		options.Qbittorrent_Host,
		options.Qbittorrent_Username, options.Qbittorrent_Password, options.Qbittorrent_RenameFile,
		options.TransmissionHost,
		options.TransmissionUsername,
		options.TransmissionPassword,
	)
	// 修改配置
	this.Sonarr.Host = configs[options.Sonarr_Host]
	this.Sonarr.Apikey = configs[options.Sonarr_ApiKey]
	this.Sonarr.Version = configs[options.Sonarr_Version]
	this.Sonarr.RefreshTime = helper.StrToInt(configs[options.Sonarr_RefreshTime])

	this.Themoviedb.Apikey = configs[options.Themoviedb_ApiKey]
	this.Themoviedb.Proxy = configs[options.Themoviedb_Proxy]
	this.Themoviedb.ProxyType = helper.StrToInt(configs[options.Themoviedb_ProxyType])

	this.Transmission.Host = configs[options.TransmissionHost]
	this.Transmission.Username = configs[options.TransmissionUsername]
	this.Transmission.Password = configs[options.TransmissionPassword]

	this.Qbittorrent.Host = configs[options.Qbittorrent_Host]
	this.Qbittorrent.Username = configs[options.Qbittorrent_Username]
	this.Qbittorrent.Password = configs[options.Qbittorrent_Password]
	this.Qbittorrent.RenameFile = configs[options.Qbittorrent_RenameFile]

}

// 保存系统配置
func (this *appConfService) Save() error {
	options.SetOption(options.Sonarr_Host, this.Sonarr.Host)
	options.SetOption(options.Sonarr_ApiKey, this.Sonarr.Apikey)
	options.SetOption(options.Sonarr_Version, this.Sonarr.Version)
	options.SetOption(options.Sonarr_RefreshTime, strconv.Itoa(this.Sonarr.RefreshTime))

	options.SetOption(options.Themoviedb_ApiKey, this.Themoviedb.Apikey)
	options.SetOption(options.Themoviedb_Proxy, this.Themoviedb.Proxy)
	options.SetOption(options.Themoviedb_ProxyType, strconv.Itoa(this.Themoviedb.ProxyType))

	options.SetOption(options.TransmissionHost, this.Transmission.Host)
	options.SetOption(options.TransmissionUsername, this.Transmission.Username)
	options.SetOption(options.TransmissionPassword, this.Transmission.Password)

	options.SetOption(options.Qbittorrent_Host, this.Qbittorrent.Host)
	options.SetOption(options.Qbittorrent_Username, this.Qbittorrent.Username)
	options.SetOption(options.Qbittorrent_Password, this.Qbittorrent.Password)
	options.SetOption(options.Qbittorrent_RenameFile, this.Qbittorrent.RenameFile)

	return nil
}

func (this *appConfService) GetTMDBProxy() string {
	switch this.Themoviedb.ProxyType {
	case TMDB_PROXY_TYPE_GLOBAL:
		return AppConf.System.GlobalProxy
	case TMDB_PROXY_TYPE_DIY:
		return this.Themoviedb.Proxy
	case TMDB_PROXY_TYPE_NONE:
		return ""
	}
	return ""
}
