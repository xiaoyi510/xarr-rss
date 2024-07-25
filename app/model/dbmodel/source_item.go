package dbmodel

// 数据源表
type SourceItem struct {
	Model
	SourceId  int                  `json:"source_id"  gorm:"uniqueIndex:hashF"`
	Hash      string               `json:"hash"  gorm:"uniqueIndex:hashF"`
	PubDate   string               `json:"pub_date"`
	Enclosure SourceItemEnclosure  `gorm:"serializer:json" json:"enclosure"`
	Content   string               `json:"content"`
	Link      string               `json:"link"`
	Guid      string               `json:"guid"`
	Title     string               `json:"title"`
	ParseInfo *SourceItemParseInfo `json:"parse_info" gorm:"serializer:json"` // 匹配的数据结果
}

type SourceItemEnclosure struct {
	Type   string `json:"type"`
	Length string `json:"length"`
	Url    string `json:"url"`
}

// 内容格式化
type SourceItemParseInfo struct {
	AnalyzeTitle string   `json:"analyze_title,omitempty"` // 解析后剩余的标题
	AudioEncode  []string `json:"audio_encode,omitempty"`  // 音频编码
	VideoEncode  []string `json:"video_encode,omitempty"`  // 视频编码
	MediaType    []string `json:"media_type,omitempty"`    // 媒体类型
	MinEpisode   int      `json:"min_episode,omitempty"`   // 最小季数 或 单集集数
	MaxEpisode   int      `json:"max_episode,omitempty"`   // 最大集数
	Version      int      `json:"version,omitempty"`       // 版本号
	Season       int      `json:"season,omitempty"`        // 季数
	ReleaseGroup string   `json:"release_group,omitempty"` // 发布组 [NC_RAW]
	//QualitySource      string   `json:"quality_source"`       // 质量 webdl
	//Resolution         string   `json:"resolution"`           // 尺寸 1080p
	QualityResolution string `json:"quality_resolution,omitempty"` // 质量_尺寸 Webdl-1080p
	Language          string `json:"language,omitempty"`           // 语言 chinese
	ProductionCompany string `json:"production_company,omitempty"` // 发布 公司
	Subtitles         string `json:"subtitles,omitempty"`          // 字幕
	//AbsoluteMinEpisode int    `json:"absolute_min_episode"` // 原始最小集 绝对集
	//AbsoluteMaxEpisode int    `json:"absolute_max_episode"` // 原始最大集 绝对集
}
