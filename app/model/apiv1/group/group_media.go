package group

import (
	"XArr-Rss/app/model/dbmodel"
)

type Apiv1GroupMediaAdd struct {
	SonarrId        int                  `json:"sonarr_id"`         // SonarrId
	Regex           []dbmodel.GroupRegex `json:"regex"`             // 正则匹配规则
	Language        string               `json:"language"`          //  语言
	Quality         string               `json:"quality"`           // 质量
	UseSource       []string             `json:"use_source"`        // 使用那些数据源
	FilterPushGroup []string             `json:"filter_push_group"` // 过滤发布组
	EchoTitleAnime  string               `json:"echo_title_anime"`
	EchoTitleTv     string               `json:"echo_title_tv"`
}
