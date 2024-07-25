package torrent_title_parse

import (
	"XArr-Rss/util/array"
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
)

var medieTypeReg = regexp2.MustCompile(`\b(?:(gb|big5)_)?(?:(?<mp4>mp4)|(?<mkv>mkv))\b`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)

// 媒体格式查询 mp4什么的
func (this *TorrentTitleParse) ParseMediaType(result *MatchResult) []string {

	regGroups := regexp_ext.ParseGroups(medieTypeReg, result.AnalyzeTitle)
	if regGroups == nil {
		return []string{}
	} else {
		ret := []string{}
		if mp4V := regGroups.GetGroupValByName("mp4"); len(mp4V) > 0 {
			// 删除mp4
			ret = append(ret, mp4V[0])
		}
		if mkvV := regGroups.GetGroupValByName("mkv"); len(mkvV) > 0 {
			// 删除mp4
			ret = append(ret, mkvV[0])
		}

		if len(ret) > 0 {
			result.AnalyzeTitle, _ = medieTypeReg.Replace(result.AnalyzeTitle, "$1", -1, -1)
			return array.UniqueString(ret)
		}
		return []string{}
	}

}
