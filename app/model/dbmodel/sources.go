package dbmodel

const (
	ProxySiteTypeDefaultRss = iota + 1
	ProxySiteTypeJacket
	ProxySiteTypeProwlarr
	ProxySiteTypeKeywords
)

type Sources struct {
	Model
	Name            string       `json:"name" gorm:"unique:true"`           // 名称
	Url             string       `json:"url"`                               // 订阅地址
	UseProxy        int          `json:"use_proxy"`                         // 代理 1 使用 0 不用
	RefreshTime     int64        `json:"refresh_time"`                      // 刷新间隔时间
	LastRefreshTime int64        `json:"last_refresh_time"`                 // 最后一次刷新时间
	MaxReadCount    int          `json:"max_read_count"`                    // 获取最大数量
	CacheDay        int          `json:"cache_day"`                         // 缓存天数 超过天数的数据会自动移出
	Regex           *SourceRegex `json:"regex"  gorm:"serializer:json"`     // 正则匹配内容
	Query           string       `json:"query"`                             // 请求参数内容 query={keyword}   ---- v1.5.5
	ProxySiteType   int          `json:"proxy_site_type" gorm:"default:1;"` // 代理站点类型 1 默认rss数据源 2 jacket  3 prowlarr 4 关键字
	ProxySiteApiKey string       `json:"proxy_site_api_key"`                // 代理站点网站 Apikey
	MaxCount        int          `json:"max_count"`                         // 最大缓存数量
	Status          int          `json:"status" gorm:"default:1;"`          // 状态
	AutoSearch      string       `json:"auto_search" `                      // 自动检索模块
	DownloadPasskey string       `json:"download_passkey" `                 // 自动检索模块
}

type SourceRegex struct {
	MustHave    string `json:"must_have"`
	MustDotHave string `json:"must_dot_have"`
}
