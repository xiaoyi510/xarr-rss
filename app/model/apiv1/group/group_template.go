package group

import "XArr-Rss/app/model/dbmodel"

type Apiv1GroupTemplateAdd struct {
	Name            string               `json:"name"`               // 模板名称
	Language        string               `json:"language"`           // 语言 -1 自动搜索
	Quality         string               `json:"quality"`            // 质量
	Regex           []dbmodel.GroupRegex `json:"regex" `             // 正则匹配规则
	UseSource       []string             `json:"use_source"  `       // 使用那些数据源
	FilterPushGroup []string             `json:"filter_push_group" ` // 只能使用那些发布组
	EchoTitleAnime  string               `json:"echo_title_anime"`
	EchoTitleTv     string               `json:"echo_title_tv"`
}

type Apiv1GroupTemplateInfo struct {
	Id int32 `json:"id"`
}

type Apiv1GroupTemplateEdit struct {
	Id              int32                `json:"id"`
	Name            string               `json:"name"`               // 模板名称
	Language        string               `json:"language"`           // 语言 -1 自动搜索
	Quality         string               `json:"quality"`            // 质量
	Regex           []dbmodel.GroupRegex `json:"regex" `             // 正则匹配规则
	UseSource       []string             `json:"use_source"  `       // 使用那些数据源
	FilterPushGroup []string             `json:"filter_push_group" ` // 只能使用那些发布组
	EchoTitleAnime  string               `json:"echo_title_anime"`
	EchoTitleTv     string               `json:"echo_title_tv"`
}

type Apiv1GroupTemplateDelete struct {
	Id int32 `json:"id"`
}

type Apiv1GroupTemplateBatchUse struct {
	Id            int32   `json:"id"`
	GroupMediaIds []int32 `json:"group_media_ids"`
}
