package dbmodel

type MediaEpisodeList struct {
	Model
	SonarrId              int    `json:"sonarr_id,omitempty"  gorm:"uniqueIndex:idx_member"`
	SeasonNumber          int    `json:"season_number,omitempty" gorm:"uniqueIndex:idx_member"`  // 季号
	EpisodeNumber         int    `json:"episode_number,omitempty" gorm:"uniqueIndex:idx_member"` // 剧集第几集
	EpisodeTitle          string `json:"episode_title,omitempty"`                                // 剧集标题
	AirDate               string `json:"air_date,omitempty"`
	AirDateUtc            string `json:"air_date_utc,omitempty"` // 开播UTC时间
	EpisodeId             int    `json:"episode_id,omitempty"`
	AbsoluteEpisodeNumber int    `json:"absolute_episode_number,omitempty"` // 绝对的剧集号
	HasFile               int    `json:"has_file"`                          // 是否已有文件
}
