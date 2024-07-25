package torrent_title_parse

import (
	"XArr-Rss/util/array"
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
	"strings"
)

const (
	Language_Unknown             = "Unknown"
	Language_English             = "English"
	Language_French              = "French"
	Language_Spanish             = "Spanish"
	Language_German              = "German"
	Language_Italian             = "Italian"
	Language_Danish              = "Danish"
	Language_Dutch               = "Dutch"
	Language_Japanese            = "Japanese"
	Language_Icelandic           = "Icelandic"
	Language_Chinese             = "Chinese"             // 中文
	Language_Chinese_Simplified  = "Chinese Simplified"  // 简体
	Language_Chinese_Traditional = "Chinese Traditional" // 繁体
	Language_Russian             = "Russian"
	Language_Polish              = "Polish"
	Language_Vietnamese          = "Vietnamese"
	Language_Swedish             = "Swedish"
	Language_Norwegian           = "Norwegian"
	Language_Finnish             = "Finnish"
	Language_Turkish             = "Turkish"
	Language_Portuguese          = "Portuguese"
	Language_Flemish             = "Flemish"
	Language_Greek               = "Greek"
	Language_Korean              = "Korean"
	Language_Hungarian           = "Hungarian"
	Language_Hebrew              = "Hebrew"
	Language_Lithuanian          = "Lithuanian"
	Language_Czech               = "Czech"
	Language_Arabic              = "Arabic"
	Language_Hindi               = "Hindi"
)

var CaseSensitiveLanguageRegex = regexp2.MustCompile(`(?<lithuanian>\bLT\b)|(?<czech>\bCZ\b)|(?<polish>\bPL\b)`, regexp2.Compiled)
var LanguageRegex = regexp2.MustCompile(`(?:\W|_)
(?<italian>\b(?:ita|italian)\b)|
(?<german>german\b|videomann)|
(?<flemish>flemish)|
(?<greek>greek)|
(?<french>(?:\W|_)(?:FR)(?:\W|_))|
(?<russian>\brus\b)|
(?<hungarian>\b(?:HUNDUB|HUN)\b)|
(?<hebrew>\bHebDub\b)|
(?<baha>\bbaha\b)|
(?<japanese>\b(?:jp|附日字)\b)|
(?<polish>\b(?:PL\W?DUB|DUB\W?PL|LEK\W?PL|PL\W?LEK)\b)|
(?<chinese>\[(?:繁中|中繁)?(?:CH[ST]|BIG5|GB)\]|\[(?:CH[ST]\s+CH[ST])\]|
	\bgb\b|
	[繁简]+日(?:雙語|双语|多语)?|
\b[简繁]中\b|
[简繁][体體][内內][封嵌]|简体中字|[简繁日中]+外挂|[简繁日中]+内[封嵌]|简体|国漫|简|繁[体體]|繁|内封字幕|外挂字幕|(?:中文)?字幕|\[(?:gb|big5)_(?:jp|mp4)\]|\bbig5\b)`, regexp2.Compiled|regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace)

// 搜索匹配对应语言
var ContainsArr = map[string]string{
	"chinese": Language_Chinese,

	"french":     Language_French,
	"spanish":    Language_Spanish,
	"danish":     Language_Danish,
	"dutch":      Language_Dutch,
	"japanese":   Language_Japanese,
	"icelandic":  Language_Icelandic,
	"korean":     Language_Korean,
	"russian":    Language_Russian,
	"polish":     Language_Polish,
	"vietnamese": Language_Vietnamese,
	"swedish":    Language_Swedish,
	"norwegian":  Language_Norwegian,
	"finnish":    Language_Finnish,
	"turkish":    Language_Turkish,
	"portuguese": Language_Portuguese,
	"hungarian":  Language_Hungarian,
	"hebrew":     Language_Hebrew,
	"arabic":     Language_Arabic,
	"hindi":      Language_Hindi,
	"english":    Language_English,
	"mandarin":   Language_Chinese,
	"cantonese":  Language_Chinese,
}
var groupArr = map[string]string{
	"lithuanian": Language_Lithuanian,
	"czech":      Language_Czech,
	"polish":     Language_Polish,
	"italian":    Language_Italian,
	"german":     Language_German,
	"flemish":    Language_Flemish,
	"greek":      Language_Greek,
	"french":     Language_French,
	"russian":    Language_Russian,
	"dutch":      Language_Dutch,
	"hungarian":  Language_Dutch,
	"hebrew":     Language_Hebrew,
	"chinese":    Language_Chinese,
	"japanese":   Language_Japanese,
}

func (this *TorrentTitleParse) ParseLanguage(result *MatchResult) string {
	result.AnalyzeTitle = strings.ToLower(result.AnalyzeTitle)

	tempTitle := result.AnalyzeTitle
	reg := regexp2.MustCompile(`.*?[_. ](S\d{2}(?:E\d{2,4})*[_. ].*)`, regexp2.IgnoreCase|regexp2.Compiled)
	tempTitle, _ = reg.Replace(tempTitle, "$1", -1, -1)
	oldTempTitle := tempTitle

	// 循环查询关键字
	for keyword, languageText := range ContainsArr {
		if strings.Contains(tempTitle, keyword) {
			tempTitle = strings.Replace(tempTitle, keyword, "", -1)
			result.AnalyzeTitle = strings.Replace(result.AnalyzeTitle, oldTempTitle, tempTitle, -1)
			return languageText
		}
	}

	// 使用正则匹配
	regexLanguage, tempTitle := this.RegexLanguage(result, tempTitle)

	if regexLanguage != Language_Unknown {
		result.AnalyzeTitle = strings.Replace(result.AnalyzeTitle, oldTempTitle, tempTitle, -1)
		return regexLanguage
	}

	return Language_Unknown
}

// 匹配语言规则
func (this *TorrentTitleParse) RegexLanguage(result *MatchResult, title string) (string, string) {
	// Case insensitive
	caseSensitiveGroups := regexp_ext.ParseGroups(CaseSensitiveLanguageRegex, title)
	if caseSensitiveGroups != nil {
		title, _ = CaseSensitiveLanguageRegex.Replace(title, "", -1, -1)
		for groupName, language := range groupArr {
			if len(caseSensitiveGroups.GetGroupValByName(groupName)) > 0 {
				return language, title
			}
		}
	}

	// 根据分组 循环查询是否有合法数据
	var retLang = []string{}
	caseLanguageGroups := regexp_ext.ParseGroups(LanguageRegex, title)
	if caseLanguageGroups != nil {
		title, _ = LanguageRegex.Replace(title, "", -1, -1)
		// 返回可识别的语言
		for groupName, language := range groupArr {
			// 特殊处理nc-raws
			if strings.Contains(strings.ToLower(result.ReleaseGroup), "nc-raws") || strings.Contains(strings.ToLower(result.ReleaseGroup), "まひろ🍥") {
				if len(caseLanguageGroups.GetGroupValByName("baha")) > 0 {
					retLang = append(retLang, groupArr["chinese"])
				}
			} else if len(caseLanguageGroups.GetGroupValByName(groupName)) > 0 {
				retLang = append(retLang, language)
			}
		}
		if len(retLang) > 0 {
			return strings.Join(array.UniqueString(retLang), ","), title
		}
	}
	return Language_Unknown, title
}
