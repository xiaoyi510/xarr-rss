package torrent_title_parse

import (
	"XArr-Rss/util/array"
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
	"strings"
)

var tvNameReg = regexp2.MustCompile(`(?:
(?<bilibili>\bbilibili\b)|
(?<bglobal>\bB\-global\b)|
(?<tvb>\btvb\b)|
(?<dsn>\bdsnp\b)|
(?<viutv>\bviutv\b)
)`, regexp2.Compiled|regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace)

func (this *TorrentTitleParse) ParseTvName(result *MatchResult) {
	// 匹配获取分组数据
	regGroups := regexp_ext.ParseGroups(tvNameReg, result.AnalyzeTitle)
	if regGroups == nil {
		return
	}
	ret := []string{}
	ret = regGroups.GetGroupValByName("0")

	// 替换标题
	result.AnalyzeTitle, _ = tvNameReg.Replace(result.AnalyzeTitle, "", -1, -1)
	if len(ret) > 0 {
		result.ProductionCompany = strings.Join(array.UniqueString(ret), ",")
	}
}
