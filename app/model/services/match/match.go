package match

// 匹配的媒体数据
type MediaMatchInfo struct {
	Title         string
	OriginalTitle string
	OtherTitle    string
	MinEpisode    int
	MaxEpisode    int
	OldMinEpisode int `xml:"oldMinEpisode" json:"old_min_episode"`
	OldMaxEpisode int `xml:"oldMaxEpisode" json:"old_max_episode"`
	Language      string
	Quality       string
	Season        int
}
