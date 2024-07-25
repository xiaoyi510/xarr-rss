package sources

import (
	"XArr-Rss/app/model"
	"XArr-Rss/util/date_util"
	"XArr-Rss/util/logsys"
	"encoding/xml"
	"errors"
)

type GeneralSource struct {
	Channel GeneralChannel `xml:"channel"`
}

type GeneralChannel struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description string        `xml:"description"`
	Item        []GeneralItem `xml:"item"`
}
type GeneralTorrent struct {
	Xmlns         string `xml:"xmlns,attr"`
	Link          string `xml:"link"`
	ContentLength string `xml:"contentLength"`
	PubDate       string `xml:"pubDate"`
}

type GeneralItem struct {
	Title       string           `xml:"title"`
	Link        string           `xml:"link"`
	Description string           `xml:"description"`
	PubDate     string           `xml:"pubDate,omitempty"`
	Guid        GeneralGuid      `xml:"guid"`
	Enclosure   GeneralEnclosure `xml:"enclosure"`
	Torrent     GeneralTorrent   `xml:"torrent"`
}

type GeneralGuid struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Text        string `xml:",cdata"`
}

type GeneralEnclosure struct {
	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
	Url    string `xml:"url,attr"`
}

func (this GeneralSource) Parse(data string) (*model.RssRoot, error) {
	parse := &GeneralSource{}
	err := xml.Unmarshal([]byte(data), parse)
	if err != nil {
		logsys.Error(err.Error(), "XML Parse")
		return nil, errors.New("不是正确的订阅数据格式")
	}
	ret := &model.RssRoot{
		Channel: model.RssResult{
			Title:       parse.Channel.Title,
			Link:        parse.Channel.Link,
			Description: parse.Channel.Description,
		},
		Version: "2.0",
	}

	for _, v := range parse.Channel.Item {
		pubDate := v.PubDate
		if v.Torrent.PubDate != "" {
			pubDate = v.Torrent.PubDate
		}
		if pubDate == "" {
			pubDate = date_util.TimeNowStr()
		}

		ret.Channel.Item = append(ret.Channel.Item, model.RssResultItem{
			Title: model.CDATA{
				Text: v.Title,
			},
			PubDate: pubDate,
			Enclosure: model.RssResultItemEnclosure{
				Type:   v.Enclosure.Type,
				Length: v.Enclosure.Length,
				Url:    v.Enclosure.Url,
			},
			Link: v.Link,
			Guid: model.RssResultItemGuid{
				IsPermaLink: v.Guid.IsPermaLink,
				Text:        v.Guid.Text,
			},
		})
	}

	return ret, nil
}
