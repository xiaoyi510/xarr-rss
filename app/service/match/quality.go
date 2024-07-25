package match

import (
	"XArr-Rss/util/helper"
	"strings"
)

func (this MediaMatchService) ParseMediaQuality(text string) string {
	text = strings.ToLower(text)
	text = helper.ReplaceRegStringSpace(text)

	// 判断压制类型
	var t = "WebDL"
	var tArr = make(map[string]string)
	tArr = map[string]string{
		"bdrip":        "BDRip",
		"webrip":       "WebRip",
		"webdl":        "WebDL",
		"bluray":       "Bluray",
		"hdtv":         "HDTV",
		"sdtv":         "SDTV",
		"dvd":          "DVD",
		"nc-raws,baha": "WebDL",
		"[raw]":        "Raw",
	}
	for searchKey, tValue := range tArr {
		searchKey = helper.StrReplace(searchKey, []string{"【", "】"}, []string{"[", "]"})
		searchKey = helper.ReplaceRegStringSpace(searchKey)
		searchList := strings.Split(searchKey, ",")

		// 切割搜索 关键字
		i := 0
		for _, v := range searchList {
			text2 := text
			if v == "raw" && strings.Contains(text, "ncraws") {
				text2 = strings.Replace(text, "ncraws", "", -1)
			}

			if strings.Contains(text2, v) {
				i++
				if i == len(searchList) {
					t = tValue
				}
			}
		}
	}

	quality := "Unknown"
	var qualityArr = make(map[string]string)
	qualityArr = map[string]string{
		"2160p":      "2160p",
		"1080premux": "1080p Remux",
		"1080p+":     "1080p+",
		"1080p":      "1080p",
		"1920x1080":  "1080p",
		"720p":       "720p",
		"480p":       "480p",
		"hd":         "HD",
	}
	for searchKey, tValue := range qualityArr {
		if strings.Contains(text, searchKey) {
			quality = tValue
			break
		}
	}
	if quality == "Unknown" {
		return t
	}
	if t != "" {
		return t + "-" + quality

	}
	return quality

}
