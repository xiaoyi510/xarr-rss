package dbmodel

type GroupMedia struct {
	Model
	GroupId         int          `json:"group_id" gorm:"uniqueIndex:media_unique"`
	SonarrId        int          `json:"sonarr_id" gorm:"uniqueIndex:media_unique"` // SonarrId
	Language        string       `json:"language"`                                  // 语言 -1 自动搜索
	Quality         string       `json:"quality"`                                   // 质量
	Regex           []GroupRegex `json:"regex" gorm:"serializer:json"`              // 正则匹配规则
	UseSource       []string     `json:"use_source"  gorm:"serializer:json"`        // 使用那些数据源
	FilterPushGroup []string     `json:"filter_push_group"  gorm:"serializer:json"` // 只能使用那些发布组
	MediaInfo       Media        `gorm:"foreignKey:sonarr_id;references:sonarr_id"`
	EchoTitleAnime  string       `json:"echo_title_anime"`
	EchoTitleTv     string       `json:"echo_title_tv"`
	FromGroupTemp   int          `json:"from_group_temp"` // 从哪个分组模板来的
}

const (
	REG_TYPE_DEFAULT = 1 // 默认方式
	REG_TYPE_REGEXP  = 2 // 正则表达式
)

type GroupRegex struct {
	MatchType string `json:"match_type"` // 匹配类型 为auto 则为自动处理
	Reg       string `json:"reg"`        // 正则
	RegType   int    `json:"reg_type"`   // 正则类型 当匹配类型为auto 时 可以为正则表达式 或 ,| 方式  1 ,| 检查 2 正则 默认1
	Season    int    `json:"season"`     // 对应季度
	Offset    int    `json:"offset"`     // 集数偏移
}

////////////////////////////////写法2
//type GroupRegexs []GroupRegex

//
//func (loc GroupRegexs) GormDataType() string {
//	return "json"
//}
//func (loc GroupRegexs) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
//	marshal, err := json.Marshal(loc)
//	if err != nil {
//		panic("GG")
//		return clause.Expr{}
//	}
//
//	return clause.Expr{
//		SQL:  "?",
//		Vars: []interface{}{string(marshal)},
//	}
//}
//func (loc *GroupRegexs) Scan(v interface{}) error {
//	// Scan a value into struct from database driver
//	err := json.Unmarshal([]byte(v.(string)), loc)
//	if err != nil {
//		return err
//	}
//	return nil
//}
