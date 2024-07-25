package torrent_title_parse

import (
	"XArr-Rss/util/helper"
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
	"strconv"
	"strings"
)

var seasonReg = regexp2.MustCompile(`(?:\b(?:season|Season|S|s)[\s-]?(?<season>\d+)(?:\b|E\d+)|
(?:第)[\s-]?(?<season_zh>[\d一二三四五六七八九十]+[季期])|
(?<season_pre>\b\d+nd\s+Season)|
(?<season_sp>(?:\b特别篇\b))
)`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)

// 解析季
func (this *TorrentTitleParse) ParseSeason(result *MatchResult) {

	matchString, err := seasonReg.MatchString(result.AnalyzeTitle)
	if err != nil || matchString == false {
		return
	}

	// 获取分组数据
	regGroups := regexp_ext.ParseGroups(seasonReg, result.AnalyzeTitle)
	if regGroups != nil {
		//result.AnalyzeTitle, _ = seasonReg.Replace(result.AnalyzeTitle, " ", -1, -1)
		season := ""

		if seasonGrp := regGroups.GetGroupValByName("season"); len(seasonGrp) > 0 {
			season = seasonGrp[0]
		} else if seasonPre := regGroups.GetGroupValByName("season"); len(seasonPre) > 0 {
			season = seasonPre[0]
		} else if seasonSp := regGroups.GetGroupValByName("season_sp"); len(seasonSp) > 0 {
			season = "0"
		} else if seasonZh := regGroups.GetGroupValByName("season_zh"); len(seasonZh) > 0 {
			season = seasonZh[0]
		} else {
			// 没有匹配到
			return
		}
		season = helper.StrReplace(season, []string{"季", "期"}, []string{""})
		arr := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
		for i, v := range arr {
			season = strings.Replace(season, v, strconv.Itoa(i+1), -1)
		}
		if season != "" {
			result.Season = helper.StrToInt(season)
			result.AnalyzeTitle, _ = seasonReg.Replace(result.AnalyzeTitle, " ", -1, -1)
			return
		}
	}

}
