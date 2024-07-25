package dbmodel

const (
	AutoInsertSonarrNone = iota
	AutoInsertSonarrNew
	AutoInsertSonarrAuto
	AutoInsertSonarrTags
)

type Group struct {
	Model
	Name               string `json:"name" gorm:"unique"`
	AutoInsertSonarr   int    `json:"auto_insert_sonarr"`    // 自动同步Sonarr数据 0 手动 1 同步最新 2 同步所有 3 同步指定Tag
	LastInsertSonarrId int    `json:"last_insert_sonarr_id"` // 最后同步的一次id
	Tags               string `json:"tags"`                  // 需要同步那些指定的Tag
	GroupTemplateId    int32  `json:"group_template_id"`     // 创建时使用的模板ID
}
