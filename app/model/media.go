package model

type Media struct {
	SonarrId          string        `json:"sonarrId"`
	TitleSlug         string        `json:"title_slug"` // 用于标题信息
	CnTitle           string        `json:"cn_title"`
	AlternateTitles   []string      `json:"alternate_titles"` // Sonarr里面的其他标题
	Titles            []string      `json:"titles"`           // 所有标题
	Title             string        `json:"title"`
	SeasonCount       int           `json:"season_count"`
	TotalEpisodeCount int           `json:"totalEpisodeCount"`
	EpisodeCount      int           `json:"episodeCount"`
	Year              int           `json:"year"`
	TvdbId            int           `json:"tvdbId"`
	TvRageId          int           `json:"tvRageId"`
	TvMazeId          int           `json:"tvMazeId"`
	ImdbId            string        `json:"imdbId"`
	TmdbId            string        `json:"tmdbId"`
	ErrTime           int64         `json:"err_time"`
	Seasons           []MediaSeason `json:"seasons"`
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
