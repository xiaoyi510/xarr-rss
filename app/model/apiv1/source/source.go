package source

import "XArr-Rss/app/model/dbmodel"

type Apiv1SourceAddReq struct {
	Name string `json:"name"`
	//Site         string               `json:"site"`
	Url             string               `json:"url"`
	UseProxy        int                  `json:"use_proxy"`
	RefreshTime     int64                `json:"refresh_time"`
	MaxReadCount    int                  `json:"max_read_count"`
	CacheDay        int                  `json:"cache_day"` // 缓存天数
	Regex           *dbmodel.SourceRegex `json:"regex"`
	ProxySiteType   int                  `json:"proxy_site_type" `   // 代理站点类型 1 默认rss数据源 2 jacket  3 prowlarr 4 关键字
	ProxySiteApiKey string               `json:"proxy_site_api_key"` // 代理站点网站 Apikey
	MaxCount        int                  `json:"max_count"`          // 最大缓存数量
	Status          int                  `json:"status"`             // 状态
	AutoSearch      string               `json:"auto_search"`        // 自动扩展搜索
	DownloadPasskey string               `json:"download_passkey"`   // 下载地址自动增加passkey

}

type Apiv1SourceEditReq struct {
	Name string `json:"name"`
	//Site         string               `json:"site"`
	Url             string               `json:"url"`
	UseProxy        int                  `json:"use_proxy"`
	RefreshTime     int64                `json:"refresh_time"`
	MaxReadCount    int                  `json:"max_read_count"`
	CacheDay        int                  `json:"cache_day"` // 缓存天数
	Regex           *dbmodel.SourceRegex `json:"regex"`
	ProxySiteType   int                  `json:"proxy_site_type"`    // 代理站点类型 1 默认rss数据源 2 jacket  3 prowlarr 4 关键字
	ProxySiteApiKey string               `json:"proxy_site_api_key"` // 代理站点网站 Apikey
	MaxCount        int                  `json:"max_count"`          // 最大缓存数量
	Status          int                  `json:"status"`             // 状态
	AutoSearch      string               `json:"auto_search"`        // 自动扩展搜索
	DownloadPasskey string               `json:"download_passkey"`   // 下载地址自动增加passkey

}
