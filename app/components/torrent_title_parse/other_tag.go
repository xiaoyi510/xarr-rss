package torrent_title_parse

import (
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
)

var otherTagReg = regexp2.MustCompile(`(?:
(?<baha>\bBaha\b)|
(?<bglobal>\bB-Global\b)|
(?<cr>\bcr\b)|
(?<sentai>\bSentai\b)|
)`, regexp2.Compiled|regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace)

func (this *TorrentTitleParse) ParseOtherTag(result *MatchResult) {
	//// 匹配获取分组数据
	regGroups := regexp_ext.ParseGroups(otherTagReg, result.AnalyzeTitle)
	if regGroups == nil {
		return
	}
	//result.AnalyzeTitle, _ = subtitlesReg.Replace(result.AnalyzeTitle, "", -1, -1)

	ret := []string{}
	ret = regGroups.GetGroupValByName("baha")
	if len(ret) > 0 {
		result.OtherTag = append(result.OtherTag, "Baha")
	}
	ret = regGroups.GetGroupValByName("bglobal")
	if len(ret) > 0 {
		result.OtherTag = append(result.OtherTag, "B-Global")
	}
	ret = regGroups.GetGroupValByName("cr")
	if len(ret) > 0 {
		result.OtherTag = append(result.OtherTag, "cr")
	}
	ret = regGroups.GetGroupValByName("sentai")
	if len(ret) > 0 {
		result.OtherTag = append(result.OtherTag, "Sentai")
	}
	//
	//// 替换标题
	//if len(ret) > 0 {
	//	result.Subtitles = strings.Join(array.UniqueString(ret), ",")
	//}
}
