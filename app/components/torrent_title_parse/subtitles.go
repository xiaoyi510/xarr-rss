package torrent_title_parse

import (
	"XArr-Rss/util/array"
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
	"strings"
)

var subtitlesReg = regexp2.MustCompile(`(?:
(?<srt>\bsrt\b)
)`, regexp2.Compiled|regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace)

func (this *TorrentTitleParse) ParseSubtitles(result *MatchResult) {
	// 匹配获取分组数据
	regGroups := regexp_ext.ParseGroups(subtitlesReg, result.AnalyzeTitle)
	if regGroups == nil {
		return
	}
	ret := []string{}
	ret = regGroups.GetGroupValByName("0")

	// 替换标题
	result.AnalyzeTitle, _ = subtitlesReg.Replace(result.AnalyzeTitle, "", -1, -1)
	if len(ret) > 0 {
		result.Subtitles = strings.Join(array.UniqueString(ret), ",")
	}
}
