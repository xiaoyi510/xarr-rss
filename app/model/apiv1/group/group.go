package group

type ApiV1GroupsAddReq struct {
	Name             string `json:"name"`
	AutoInsertSonarr int    `json:"auto_insert_sonarr"` // 自动同步Sonarr数据 0 手动 1 同步最新 2 同步所有
	Tags             string `json:"tags"`
	GroupTemplateId  int32  `json:"group_template_id"` // 创建时使用的模板ID

}

type ApiV1GroupsEditReq struct {
	Name             string `json:"name"`
	AutoInsertSonarr int    `json:"auto_insert_sonarr"` // 自动同步Sonarr数据 0 手动 1 同步最新 2 同步所有
	Tags             string `json:"tags"`
	GroupTemplateId  int32  `json:"group_template_id"` // 创建时使用的模板ID

}
