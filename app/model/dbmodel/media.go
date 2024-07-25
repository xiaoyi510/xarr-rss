package dbmodel

type Media struct {
	SonarrId          int           `json:"sonarr_id" gorm:"primarykey;autoIncrement:false"`
	TitleSlug         string        `json:"title_slug"` // 用于标题信息
	CnTitle           string        `json:"cn_title"`
	Overview          string        `json:"overview"`
	AlternateTitles   []string      `json:"alternate_titles" gorm:"serializer:json"` // Sonarr里面的其他标题
	Titles            []string      `json:"titles" gorm:"serializer:json"`           // 所有标题
	OriginalTitle     string        `json:"original_title"`
	SeasonCount       int           `json:"season_count"`
	TotalEpisodeCount int           `json:"totalEpisodeCount"`
	EpisodeCount      int           `json:"episodeCount"`
	Monitored         int           `json:"monitored"`
	Year              int           `json:"year"`
	TvdbId            int           `json:"tvdbId" gorm:"column:tvdb_id"`
	TvRageId          int           `json:"tvRageId" gorm:"column:tvrage_id"`
	TvMazeId          string        `json:"tvMazeId" gorm:"column:tvmaze_id"`
	ImdbId            string        `json:"imdbId" gorm:"column:imdb_id"`
	TmdbId            string        `json:"tmdbId" gorm:"column:tmdb_id"`
	Image             string        `json:"image" gorm:"column:image"`
	ErrTime           int64         `json:"err_time"`
	SeriesType        string        `json:"series_type"`
	Seasons           []MediaSeason `json:"seasons"  gorm:"serializer:json"`
	Tags              []string      `json:"tags" gorm:"serializer:json"` // 标签
	SearchTitle       string        `gorm:"-" json:"-"`
}

type MediaSeason struct {
	SeasonNumber int  `json:"seasonNumber"`
	Monitored    bool `json:"monitored"`
	Statistics   struct {
		EpisodeFileCount  int `json:"episodeFileCount"`
		EpisodeCount      int `json:"episodeCount"`
		TotalEpisodeCount int `json:"totalEpisodeCount"`
	} `json:"statistics"`
}
