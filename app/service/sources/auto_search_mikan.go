package sources

import (
	"XArr-Rss/app/model"
	"XArr-Rss/util/date_util"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

type AutoSearchMikan struct {
}

func (t AutoSearchMikan) Parse(url, queryUri, data string) (*model.RssRoot, error) {
	if data == "" {
		return nil, errors.New("格式化数据缺失")
	}
	ret := &model.RssRoot{
		Channel: model.RssResult{
			Title:       "蜜柑自动检索",
			Link:        "",
			Description: "XArr 自动检索",
		},
		Version: "2.0",
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("#sk-container .fadeIn > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		link, _ := s.Find("td:nth-child(1) > a.magnet-link-wrap").Attr("href")
		title := s.Find("td:nth-child(1) .magnet-link-wrap").Text()
		length := s.Find("td:nth-child(2)").Text()
		pubDate := s.Find("td:nth-child(3)").Text()
		href, _ := s.Find("td:nth-child(4)>a").Attr("href")
		//fmt.Printf("种子文件名 %d: %s ,%s\n", i, title, href)

		// 判断是否是列表页面
		if strings.Contains(queryUri, "Home/Classic") {
			// 按照class的形式走
			pubDate = s.Find("td:nth-child(1)").Text()
			pubDate = strings.Replace(pubDate, "昨天", date_util.PreDayDate(), -1)
			pubDate = strings.Replace(pubDate, "今天", date_util.TodayDate(), -1)
			link, _ = s.Find("td:nth-child(3) > a.magnet-link-wrap").Attr("href")
			if link != "" && strings.Index(link, "/Home/Episode") == 0 {
				link = url + link
			}
			length = s.Find("td:nth-child(4)").Text()
			href, _ = s.Find("td:nth-child(5)>a").Attr("href")
			title = s.Find("td:nth-child(3) .magnet-link-wrap").Text()

		}

		if pubDate == "" || title == "" || href == "" {
			return
		}

		ret.Channel.Item = append(ret.Channel.Item, model.RssResultItem{
			Title: model.CDATA{
				Text: title,
			},
			PubDate: pubDate,
			Enclosure: model.RssResultItemEnclosure{
				Type:   "application/x-bittorrent",
				Length: length,
				Url:    url + strings.TrimLeft(href, "/"),
			},
			Link: link,
			Guid: model.RssResultItemGuid{
				IsPermaLink: false,
				Text:        "", //hash.Md5{}.HashString(url + link + length + pubDate),
			},
		})
	})

	return ret, nil
}
