package model

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"encoding/xml"
)

type RssRoot struct {
	XMLName xml.Name  `xml:"rss" json:"XMLName"`
	Version string    `xml:"version,attr" json:"version,omitempty"`
	Channel RssResult `xml:"channel" json:"channel"`
}

type RssResult struct {
	Title       string          `xml:"title" json:"title,omitempty"`
	Link        string          `xml:"link" json:"link,omitempty"`
	Description string          `xml:"description" json:"description,omitempty"`
	Item        []RssResultItem `xml:"item" json:"item,omitempty"`
}

type CDATA struct {
	Text string `xml:",cdata" json:"text,omitempty"`
}

type RssResultItem struct {
	Title         CDATA                  `xml:"title" json:"title"`
	OriginalTitle CDATA                  `xml:"originalTitle" json:"originalTitle"`
	OtherTitle    CDATA                  `xml:"other_title" json:"other_title"`
	CnTitle       string                 `xml:"cn_title" json:"cn_title"`
	PubDate       string                 `xml:"pubDate" json:"pubDate,omitempty"`
	Enclosure     RssResultItemEnclosure `xml:"enclosure" json:"enclosure"`
	Link          string                 `xml:"link" json:"link,omitempty"`
	Guid          RssResultItemGuid      `xml:"guid" json:"guid"`

	Season         int                             `xml:"season" json:"season,omitempty"`
	SourceId       string                          `xml:"source_id,omitempty" json:"source_id,omitempty"`
	OldMinEpisode  int                             `xml:"oldMinEpisode" json:"old_min_episode,omitempty"`
	OldMaxEpisode  int                             `xml:"oldMaxEpisode" json:"old_max_episode,omitempty"`
	MinEpisode     int                             `xml:"minEpisode" json:"min_episode,omitempty"`
	MaxEpisode     int                             `xml:"maxEpisode" json:"max_episode,omitempty"`
	OtherInfo      torrent_title_parse.MatchResult `xml:"-" json:"-"`
	XArrRssIndexer RssResultItemXArrRssIndexer     `xml:"XArrRssIndexer,omitempty"`
	OthderId       RssResultItemOtherId            `xml:"OthderId"`
}
type RssResultItemXArrRssIndexer struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}

type RssResultItemOtherId struct {
	Text   string `xml:",chardata"`
	TmdbId string `xml:"tmdbId,attr"`
	TvdbId string `xml:"tvdbId,attr"`
	ImdbId string `xml:"imdbId,attr"`
}

type RssResultItemGuid struct {
	IsPermaLink bool   `xml:"isPermaLink,attr" json:"isPermaLink,omitempty"`
	Text        string `xml:",cdata" json:"text,omitempty"`
}
type RssResultItemEnclosure struct {
	Type   string `xml:"type,attr" json:"type,omitempty"`
	Length string `xml:"length,attr,omitempty" json:"length,omitempty"`
	Url    string `xml:"url,attr" json:"url,omitempty"`
}

/*<rss version="2.0">
<channel>
<title>Mikan Project - 我的番组</title>
<link>http://mikanani.me/RSS/MyBangumi?token=OpeuUb3rOXk%2B0T%2Fv48yu6d8c%2BDWpH0uLFsTfOZNkc%2Bo%3D</link>
<description>Mikan Project - 我的番组</description>
<item>
<title>Ya Boy Kongming! - S01E04 - Chinese - WEBDL 1080p</title>
<pubDate>2022-04-21T22:04:43.854</pubDate>
<enclosure type="application/x-bittorrent" length="667177472" url="https://mikanani.me/Download/20220421/34b80022e46f52de51e4b05ae783f80353513710.torrent"/>
<link>https://mikanani.me/Home/Episode/34b80022e46f52de51e4b05ae783f80353513710</link>
<guid isPermaLink="true">https://mikanani.me/Home/Episode/34b80022e46f52de51e4b05ae783f80353513710</guid>
</item>
<item>
<title>The Greatest Demon Lord Is Reborn as a Typical Nobody - S01E03 - Chinese - WEBDL 1080p</title>
<pubDate>2022-04-20T21:45:07.654</pubDate>
<enclosure type="application/x-bittorrent" length="427147904" url="https://mikanani.me/Download/20220420/685eaa0e9548f6ddc1e45c47838d551b6daea21b.torrent"/>
<link>https://mikanani.me/Home/Episode/685eaa0e9548f6ddc1e45c47838d551b6daea21b</link>
<guid isPermaLink="true">https://mikanani.me/Home/Episode/685eaa0e9548f6ddc1e45c47838d551b6daea21b</guid>
</item>
<item>
<title>The Rising of the Shield Hero - S02E03 - Chinese - WEBDL 1080p</title>
<pubDate>2022-04-20T21:45:00</pubDate>
<enclosure type="application/x-bittorrent" length="584371392" url="https://mikanani.me/Download/20220420/5c65697f0705b69b7a070bca69076d2565950b39.torrent"/>
<link>https://mikanani.me/Home/Episode/5c65697f0705b69b7a070bca69076d2565950b39</link>
<guid isPermaLink="true">https://mikanani.me/Home/Episode/5c65697f0705b69b7a070bca69076d2565950b39</guid>
</item>
<item>
<title>Aharen-san wa Hakarenai - S01E03 - Chinese - WEBDL 1080p</title>
<pubDate>2022-04-16T01:27:58.818</pubDate>
<enclosure type="application/x-bittorrent" length="550408000" url="https://mikanani.me/Download/20220416/4be6d9337530a2a1e0b1145dd3dda8f3727d3aaa.torrent"/>
<link>https://mikanani.me/Home/Episode/4be6d9337530a2a1e0b1145dd3dda8f3727d3aaa</link>
<guid isPermaLink="true">https://mikanani.me/Home/Episode/4be6d9337530a2a1e0b1145dd3dda8f3727d3aaa</guid>
</item>
</channel>
</rss>*/
