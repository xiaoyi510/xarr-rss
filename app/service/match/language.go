package match

import (
	"XArr-Rss/util/helper"
	"strings"
)

func (this MediaMatchService) ParseMediaLanguage(text string) string {
	l := "Unknown"
	text = strings.ToLower(text)
	text = helper.ReplaceRegStringSpace(text)
	text = helper.StrReplace(text, []string{"【", "】"}, []string{"[", "]"})

	var tArr = make(map[string]string)
	tArr = map[string]string{
		"big5":        "Chinese Traditional",
		"ncraws,baha": "Chinese Traditional",
		"cht":         "Chinese Traditional",
		"[cn]":        "Chinese Simplified",
		"[zh]":        "Chinese Simplified",
		"zh-cn":       "Chinese Simplified",
		"zhcn":        "Chinese Simplified",
		"chs":         "Chinese Simplified",
		"[gb]":        "Chinese Simplified",
		"[gb_mp4]":    "Chinese Simplified",
		"国漫":          "Chinese",
		"简体":          "Chinese Simplified",
		"简繁":          "Chinese Simplified",
		"简日":          "Chinese Simplified",
		"繁中":          "Chinese Traditional",
		"中繁":          "Chinese Traditional",
		"简中":          "Chinese Simplified",
		"中简":          "Chinese Simplified",
		"繁體":          "Chinese Traditional",
		"繁日":          "Chinese Traditional",
		"繁体":          "Chinese Traditional",
		"unknown":     "Unknown",
		"vietnamese":  "Vietnamese",
		"turkish":     "Turkish",
		"swedish":     "Swedish",
		"spanish":     "Spanish",
		"russian":     "Russian",
		"portuguese":  "Portuguese",
		"polish":      "Polish",
		"norwegian":   "Norwegian",
		"lithuanian":  "Lithuanian",
		"korean":      "Korean",
		"japanese":    "Japanese",
		"[jp]":        "Japanese",
		"jpa":         "Japanese",
		"italian":     "Italian",
		"icelandic":   "Icelandic",
		"hungarian":   "Hungarian",
		"hindi":       "Hindi",
		"hebrew":      "Hebrew",
		"greek":       "Greek",
		"german":      "German",
		"french":      "French",
		"flemish":     "Flemish",
		"finnish":     "Finnish",
		"english":     "English",
		"[en]":        "English",
		"dutch":       "Dutch",
		"danish":      "Danish",
		"czech":       "Czech",
		"chinese":     "Chinese",
		"bulgarian":   "Bulgarian",
		"arabic":      "Arabic",
	}
	for searchKey, tValue := range tArr {

		searchList := strings.Split(searchKey, ",")

		// 切割搜索 关键字
		i := 0
		for _, v := range searchList {
			if strings.Contains(text, v) {
				i++
				if i == len(searchList) {
					return tValue
				}
			}
		}

	}
	return l
}
