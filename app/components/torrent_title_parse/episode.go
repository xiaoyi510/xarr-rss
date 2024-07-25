package torrent_title_parse

import (
	"XArr-Rss/util/helper"
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
)

var episodeReg = regexp2.MustCompile(`(?!^|\[20\d\d\])
(?:[\[\sEe【第_\(]|Ep)
(?<!season\s)(?<episode>\d+)(?:集|话)?[\-\~]{0,1}
(?:s\d+)?e?
(?<max_episode>\d+)?
(?<version>v\d+)?
(?:[\)】集话\]\s](?:Fin)?|\s?Fin|$)`,
	regexp2.IgnoreCase|regexp2.Compiled|regexp2.IgnorePatternWhitespace)

// 解析第几集
func (this *TorrentTitleParse) ParseEpisode(result *MatchResult, reg *regexp2.Regexp) {
	matchString, err := episodeReg.MatchString(result.AnalyzeTitle)
	if err != nil || matchString == false {
		return
	}
	// 获取分组数据
	regGroups := regexp_ext.ParseGroups(episodeReg, result.AnalyzeTitle)
	if regGroups != nil {
		// 替换分组内容
		result.AnalyzeTitle, _ = episodeReg.Replace(result.AnalyzeTitle, " ", -1, -1)

		// 判断是否有集数
		if episode := regGroups.GetGroupValByName("episode"); len(episode) > 0 {
			result.MinEpisode = helper.StrToInt(episode[len(episode)-1])
		}

		// 判断是否有最大集数
		if maxEpisode := regGroups.GetGroupValByName("max_episode"); len(maxEpisode) > 0 {
			result.MaxEpisode = helper.StrToInt(maxEpisode[0])
		}
	}
}
