package dbmodel

type GroupTemplate struct {
	Model
	Name            string       `json:"name"`                                      // 模板名称
	Language        string       `json:"language"`                                  // 语言 -1 自动搜索
	Quality         string       `json:"quality"`                                   // 质量
	Regex           []GroupRegex `json:"regex" gorm:"serializer:json"`              // 正则匹配规则
	UseSource       []string     `json:"use_source"  gorm:"serializer:json"`        // 使用那些数据源
	FilterPushGroup []string     `json:"filter_push_group"  gorm:"serializer:json"` // 只能使用那些发布组
	EchoTitleAnime  string       `json:"echo_title_anime"`
	EchoTitleTv     string       `json:"echo_title_tv"`
}
