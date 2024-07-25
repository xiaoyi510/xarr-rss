package match

import (
	"XArr-Rss/util/helper"
	"github.com/xiaoyi510/regex_group"
	"regexp"
)

func ParseMediaReleaseGroup(title string) string {
	title = helper.StrReplace(title, []string{"【", "】"}, []string{"[", "]"})
	reg := regexp.MustCompile(`(?m)^(?P<releaseGroup>\[.*?\]).*?`)
	match := regex_group.GetRegxGroupByRegOne(reg, title)
	if match == nil {
		return ""
	}
	// 获取匹配发布组
	return match.GetString("releaseGroup")
}
