package sources

import (
	"XArr-Rss/app/model"
	"XArr-Rss/util/date_util"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

type NyaaRssSource struct {
	XMLName xml.Name    `xml:"rss"`
	Channel NyaaChannel `xml:"channel"`
}

type NyaaChannel struct {
	Text  string `xml:",chardata"`
	Title struct {
		Text string `xml:",chardata"`
	} `xml:"title"`
	Description struct {
		Text string `xml:",chardata"`
	} `xml:"description"`
	Link struct {
		Text string `xml:",chardata"`
		//Href string `xml:"href,attr"`
		//Rel  string `xml:"rel,attr"`
		//Type string `xml:"type,attr"`
	} `xml:"link"`
	Item []NyaaChannelItem `xml:"item"`
}

type NyaaChannelItem struct {
	Text  string `xml:",chardata"`
	Title struct {
		Text string `xml:",chardata"`
	} `xml:"title"` // 标题
	Link struct {
		Text string `xml:",chardata"`
	} `xml:"link"` // 下载地址
	Guid struct {
		Text        string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
	PubDate struct {
		Text string `xml:",chardata"`
	} `xml:"pubDate"`
	Seeders struct {
		Text string `xml:",chardata"`
	} `xml:"seeders"`
	Leechers struct {
		Text string `xml:",chardata"`
	} `xml:"leechers"`
	Downloads struct {
		Text string `xml:",chardata"`
	} `xml:"downloads"`
	InfoHash struct {
		Text string `xml:",chardata"`
	} `xml:"infoHash"`
	CategoryId struct {
		Text string `xml:",chardata"`
	} `xml:"categoryId"`
	Category struct {
		Text string `xml:",chardata"`
	} `xml:"category"`
	Size struct {
		Text string `xml:",chardata"`
	} `xml:"size"`
	Comments struct {
		Text string `xml:",chardata"`
	} `xml:"comments"`
	Trusted struct {
		Text string `xml:",chardata"`
	} `xml:"trusted"`
	Remake struct {
		Text string `xml:",chardata"`
	} `xml:"remake"`
	Description struct {
		Text string `xml:",chardata"`
	} `xml:"description"`
}

func (this NyaaRssSource) Parse(data string) (*model.RssRoot, error) {
	parse := &NyaaRssSource{}
	err := xml.Unmarshal([]byte(data), parse)
	if err != nil {
		logsys.Error(err.Error(), "XML Parse")
		return nil, errors.New("不是正确的订阅数据格式")
	}
	ret := &model.RssRoot{
		Channel: model.RssResult{
			Title:       parse.Channel.Title.Text,
			Link:        parse.Channel.Link.Text,
			Description: parse.Channel.Description.Text,
		},
		Version: "2.0",
	}

	for _, v := range parse.Channel.Item {
		pubDate := v.PubDate.Text
		if v.PubDate.Text != "" {
			pubDate = v.PubDate.Text
		}
		if pubDate == "" {
			pubDate = date_util.TimeNowStr()
		}

		ret.Channel.Item = append(ret.Channel.Item, model.RssResultItem{
			Title: model.CDATA{
				Text: v.Title.Text,
			},
			PubDate: pubDate,
			Enclosure: model.RssResultItemEnclosure{
				Type:   "application/x-bittorrent",
				Length: mbToSize(v.Size.Text),
				Url:    v.Link.Text,
			},
			Link: v.Guid.Text,
			Guid: model.RssResultItemGuid{
				IsPermaLink: v.Guid.IsPermaLink == "true",
				Text:        v.Guid.Text,
			},
		})
	}

	return ret, nil
}

func mbToSize(mb string) string {
	mb = strings.ToUpper(mb)
	if strings.Contains(mb, "MB") {
		a := strings.Replace(mb, "MB", "", -1)
		return strconv.Itoa(helper.StrToInt(a) * 100)
	}
	return "0"
}
