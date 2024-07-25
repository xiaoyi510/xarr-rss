package medias

type Apiv1MediasTestReq struct {
	//Titles          []string `json:"titles"`
	UseSource       []string `json:"use_source"` // 使用那些数据源
	Reg             string   `json:"reg"`
	RegType         int      `json:"regType"`
	Season          int      `json:"season"`
	FilterPushGroup []string `json:"filter_push_group"` // 过滤发布组
	Offset          int      `json:"offset"`
	Language        string   `json:"language"`
	Quality         string   `json:"quality"`
	MatchType       string   `json:"matchType"`
	SonarrId        int      `json:"sonarr_id"`
	EchoTitleAnime  string   `json:"echo_title_anime"`
	EchoTitleTv     string   `json:"echo_title_tv"`
}
