package dbmodel

type Options struct {
	Model
	Key   string `json:"key,omitempty" gorm:"uniqueIndex"`
	Value string `json:"value,omitempty"`

	//HttpAddr   string `json:"http_addr"`   // 外网映射地址
	//LogLevel   string `json:"log_level"`   // 日志等级
	//LicenseKey string `json:"license_key"` // 授权信息
}
