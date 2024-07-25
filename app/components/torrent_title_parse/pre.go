package torrent_title_parse

import (
	"XArr-Rss/util/helper"
	"github.com/dlclark/regexp2"
)

// 前置格式化

var PreSubstitutionRegex []ReplaceOption

type ReplaceOption struct {
	Regex   *regexp2.Regexp
	Replace string
	StartAt int
}

var SimpleTitleRegex = ReplaceOption{
	Regex:   regexp2.MustCompile(`(?:(480|720|1080|2160)[ip]|[xh][\W_]?26[45]|DD\W?5\W1|[<>?*|]|848x480|1280x720|1920x1080|3840x2160|4096x2160|(8|10)b(it)?|10-bit)\s*?`, regexp2.Compiled|regexp2.IgnoreCase),
	Replace: "",
	StartAt: -1,
}

var ReportTitleRegex = []*regexp2.Regexp{}

func init() {
	// 替换国漫为 gb
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`(?<left>.*)(\[国漫\])(?<right>.*)`, regexp2.Compiled),
		Replace: `${left}${right} [gb]`,
		StartAt: -1,
	})

	// [DAY字幕组][女神的露天咖啡厅/Megami no Café Terrace][第2话「与婆婆的约定」[WEB1080P] 去除 「与婆婆的约定」这个格式少一个 ] 所以替换为]
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`(\[第\d+话「[^\]]*?」)\[`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: `$1][`,
		StartAt: -1,
	})

	// 替换 03（Premium） 03（Onair） 03 为集数 去除 Premium  Onair
	//
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`(\b\d+)[\(\[](?:Premium|Onair)[\]\)]`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: `$1`,
		StartAt: -1,
	})

	// 替换jacket 将漫猫标题提到最前面 孤独摇滚！ [连载][DMG&LoliHouse][BOCCHI THE ROCK!][1080p][02][MKV/WEBRip][2022年10月]
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		//Regex:   regexp2.MustCompile(`(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)\[连载\]\[(?<subgroup>[^\]]+)\](.*)`, regexp2.Compiled),
		// 替换chineseTitle 为 title
		Regex:   regexp2.MustCompile(`(?<title>[^\]]+?)\[连载\]\[(?<subgroup>[^\]]+)\](.*)`, regexp2.Compiled),
		Replace: `[连载][${subgroup}][${title}]$1`,
		StartAt: -1,
	})
	// [4月新番]
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex: regexp2.MustCompile(`(?:
[\[\(](?:招募翻[译譯](?:校[對对]|片源)|
\[?\d月新番\]?|\d{4}年\d月[新]?番|
\d+年日剧|
TV\sSeries\(电视剧\)|
字幕[组社]招人内详|
招募(?:翻译|校[對对])|
周日版|
\b(?:仅限港澳台地区|无圣光版)\b|
\b(?:viutv粤语|特效歌词|tvb粤语)\b|
急招(?:校[對对]|翻译)、(?:[後后]期|校[對对]))[\]\)]|
(?:(\s+)?[★]?[\[\(]?\d+月新番[\]\(]?[★]?(\s+)?|无水印高清下载)|
\[(?:生[肉]?|粤)\]|
(?:\+)?小剧场|(?:\+)?小劇場|
^\[V\d修复版\]|
^\[(?:连载|TV|V2|短片|剧场)\]
)`, regexp2.Compiled|regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace),
		Replace: ``,
		StartAt: -1,
	})

	// 去除 2022.01.02 2022年2月12
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`.?\b(?<date>20\d{2}[\.年]\d{1,2}[\.月](?:\d{1,2})?)\b.?`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: ``,
		StartAt: -1,
	})

	// 去除 [2022][02] 集这种情况  (2022) [01]
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`[\[\(](?<day>20\d{2})[\]\)]\s*(\[\d+)`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: `$1`,
		StartAt: -1,
	})
	//  21[GB] 加空格
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`(\d+)(\[\b(?:gb|big5)\b\])`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: `$1 $2`,
		StartAt: -1,
	})

	// 去除 2022 2160p/WEB-DL 删除2022
	// House of the Dragon 2022 S01E06 2160p HMAX WEB-DL 10bit HDR DDP5.1 Atmos H265-BtsTV
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`(?<day>20\d{2})\s+((?:720|1080|2160)p|WEB-?DL|S\d+E\d+)`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: `$1`,
		StartAt: -1,
	})

	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`\.E(\d{2,4})\.\d{6}\.(.*-NEXT)$`, regexp2.Compiled),
		Replace: `.S01E$1.$2`,
		StartAt: -1,
	})
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`^(.*?\])(\[1080p\])(\[.*?)$`, regexp2.Compiled),
		Replace: `$1$3$2`,
		StartAt: -1,
	})

	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		Regex:   regexp2.MustCompile(`^\[(?<subgroup>[^\]]*?(?:LoliHouse|ZERO|Lilith-Raws)[^\]]*?)\](?<title>[^\[\]]+?)(?: - (?<episode>[0-9-]+)\s*|\[第?(?<episode>[0-9]+(?:-[0-9]+)?)话?(?:END|完|v\d)?\])\[`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: `[${subgroup}][${title}][${episode}][`,
		StartAt: -1,
	})

	// [爱恋&漫猫字幕组][莉可丽丝][Lycoris Recoil][06][1080p][MP4][BIG5][繁中] 替换为 [爱恋&漫猫字幕组] 莉可丽丝 Lycoris Recoil - 06 [1080p][MP4][BIG5][繁中]
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		//Regex:   regexp2.MustCompile(`^\[(?<subgroup>[^\]]+)\](?:\[国漫\])?(?:\s?★[^\[ -]+\s?)?\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)(?:\]\[|\s*[_/·]\s*))?(?<title>[^\]]+?)\](?:\[\d{4}\])?\[第?(?<episode>[0-9]+(?:-[0-9]+)?)话?(?:END|完|v\d)?\]`, regexp2.Compiled|regexp2.IgnoreCase),
		// 增加别名 alias
		//Regex:   regexp2.MustCompile(`^\[(?<subgroup>[^\]]+)\](?:\[国漫\])?(?:\s?★[^\[ -]+\s?)?(?<alias>\s?\[[^ -\]]+\]\s?)?\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)(?:\]\[|\s*[_/·]\s*))?(?<title>[^\]]+?)\](?:\[\d{4}\])?\[第?(?<episode>[0-9]+(?:-[0-9]+)?)[话集]?\s*(?:END|完|v\d|\s[^\]]*?[\u4E00-\u9FCC][^\]]*?)?\]`, regexp2.Compiled|regexp2.IgnoreCase),
		// 增加兼容 [DAY字幕组][女神的露天咖啡厅/Megami no Café Terrace][第5话「定期休息日」][简日双语][WEBrip][1080P][MP4]  「定期休息日」
		Regex:   regexp2.MustCompile(`^\[(?<subgroup>[^\]]+)\](?:\[国漫\])?(?:\s?★[^\[ -]+\s?)?(?<alias>\s?\[[^ -\]]+\]\s?)?\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)(?:\]\[|\s*[_/·]\s*))?(?<title>[^\]]+?)\](?:\[\d{4}\])?\[第?(?<episode>[0-9]+(?:-[0-9]+)?)[话集]?(?:\s*?「.*?」)?\s*(?:END|完|v\d|\s[^\]]*?[\u4E00-\u9FCC][^\]]*?)?\]`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: "[${subgroup}] ${chinesetitle} ${title} - ${episode} ",
		StartAt: -1,
	})

	// 处理以下规则
	//[GM-Team][斗罗大陆][Dou Luo Da Lu][Douro Mainland][220][AVC][GB][1080P] [gb]
	//[GM-Team][国漫][斗罗大陆][Dou Luo Da Lu][Douro Mainland][184-188][AVC][GB][1080P]
	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		//Regex:   regexp2.MustCompile(`^\[(?<subgroup>[^\]]+)\](?:\[国漫\])?(?:\s?★[^\[ -]+\s?)?\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)(?:\]\[|\s*[_/·]\s*))?(?<title>[^\]]+?)\](?:\[\d{4}\])?\[第?(?<episode>[0-9]+(?:-[0-9]+)?)话?(?:END|完|v\d)?\]`, regexp2.Compiled|regexp2.IgnoreCase),
		// 增加别名 alias
		Regex: regexp2.MustCompile(`^\[(?<subgroup>[^\]]+)\](?:\[国漫\])?(?:\s?★[^\[ -]+\s?)?\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)(?:\]\[|\s*[_/·]\s*))?(?<title>[^\]]+)(?:\]\[|\s*[_/·]\s*)(?<alias>[^\]]+)?(?:\]\[|\s*[_/·]\s*)?(?<episode>\d+(?:-\d+)?)\]`, regexp2.Compiled|regexp2.IgnoreCase),
		// 增加 (?!1080p)
		//Regex:   regexp2.MustCompile(`^\[(?<subgroup>[^\]]+)\](?:\[国漫\])?(?:\s?★[^\[ -]+\s?)?\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)(?:\]\[|\s*[_/·]\s*))?(?<title>[^\]]+)(?:\]\[|\s*[_/·]\s*)(?<alias>(?!1080p)(?:[^\]]+))?(?:\]\[|\s*[_/·]\s*)?(?<episode>\d+(?:-\d+)?)\]`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: "[${subgroup}] ${chinesetitle} ${title} - ${episode} ",
		StartAt: -1,
	})

	PreSubstitutionRegex = append(PreSubstitutionRegex, ReplaceOption{
		//Regex:   regexp2.MustCompile(`^\[(?<subgroup>[^\]]+)\](?:\[国漫\])?(?:\s?★[^\[ -]+\s?)?\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?)(?:\]\[|\s*[_/·]\s*))?(?<title>[^\]]+?)\](?:\[\d{4}\])?\[第?(?<episode>[0-9]+(?:-[0-9]+)?)话?(?:END|完|v\d)?\]`, regexp2.Compiled|regexp2.IgnoreCase),
		// 增加别名 alias
		Regex:   regexp2.MustCompile(`^^\[(?:(?<chinesetitle>[^\]]*?[\u4E00-\u9FCC][^\]]*?))(?:\]\[|\s*[_/·]\s*)(?<episode>\d+(?:-\d+)?)\]`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: "${chinesetitle} - ${episode} ",
		StartAt: -1,
	})

	///////////////////////////////////前置替换结束

	// 处理简单标题正则
	//Daily episodes without title (2018-10-12, 20181012) (Strict pattern to avoid false matches)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<airyear>19[6-9]\d|20\d\d)(?<sep>[-_]?)(?<airmonth>0\d|1[0-2])\k<sep>(?<airday>[0-2]\d|3[01])(?!\d)`, regexp2.Compiled|regexp2.IgnoreCase))

	//Multi-Part episodes without a title (S01E05.S01E06)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\W*S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:(?:[ex]){1,2}(?<episode>\d{1,3}(?!\d+)))+){2,}`, regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes without a title, Multi (S01E04E05, 1x04x05, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:(?:[-_]|[ex]){1,2}(?<episode>\d{2,3}(?!\d+))){2,})`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes without a title, Single (S01E05, 1x05)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:(?:[-_ ]?[ex])(?<episode>\d{2,3}(?!\d+))))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title Episode Absolute Episode Number ([SubGroup] Series Title Episode 01)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)(?<title>.+?)[-_. ](?:Episode)(?:[-_. ]+(?<absoluteepisode>(?<!\d+)\d{2,3}(\.\d{1,2})?(?!\d+)))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title Absolute Episode Number + Season+Episode
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\](?:_|-|\s|\.)?)(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+(?<absoluteepisode>\d{2,3}(\.\d{1,2})?))+(?:_|-|\s|\.)+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]){1,2}(?<episode>\d{2}(?!\d+)))+).*?(?<hash>[(\[]\w{8}[)\]])?(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title Season+Episode + Absolute Episode Number
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\](?:_|-|\s|\.)?)(?<title>.+?)(?:[-_\W](?<![()\[!]))+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]){1,2}(?<episode>\d{2}(?!\d+)))+)(?:(?:_|-|\s|\.)+(?<absoluteepisode>(?<!\d+)\d{2,3}(\.\d{1,2})?(?!\d+)))+.*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title Season+Episode
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\](?:_|-|\s|\.)?)(?<title>.+?)(?:[-_\W](?<![()\[!]))+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:[ex]|\W[ex]){1,2}(?<episode>\d{2}(?!\d+)))+)(?:\s|\.).*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title with trailing number Absolute Episode Number - Batch separated with tilde
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>[^-]+?)(?:(?<![-_. ]|\b[0]\d+) - )[-_. ]?(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+))~(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+))(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title with trailing number Absolute Episode Number  增加 |v\d
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>[^-]+?)(?:(?<![-_. ]|\b[0]\d+) - )(?:[-_. ]?(?<absoluteepisode>\d{2,3}(\.\d{1,2}|v\d)?(?!\d+)))+(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title with trailing number Absolute Episode Number  增加 |v\d
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\])?[-_. ]?(?<title>[^-]+?)(?:(?<![-_. ]|\b[0]\d+) - )(?:[-_. ]?(?<absoluteepisode>\d{2,3}(\.\d{1,2}|v\d)?(?!\d+)))+(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title with trailing number Absolute Episode Number
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>[^-]+?)(?:(?<![-_. ]|\b[0]\d+)[_ ]+)(?:[-_. ]?(?<absoluteepisode>\d{3}(\.\d{1,2})?(?!\d+)))+(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title - Absolute Episode Number
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>.+?)(?:(?<!\b[0]\d+))(?:[. ]-[. ](?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+|[-])))+(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title Absolute Episode Number - Absolute Episode Number (batches without full separator between title and absolute episode numbers)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>.+?)(?:(?<!\b[0]\d+))(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+|[-]))[. ]-[. ](?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+|[-]))(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title Absolute Episode Number
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>.+?)[-_. ]+\(?(?:[-_. ]?#?(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+)))+\)?(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - [SubGroup] Title Absolute Episode Number [梦蓝字幕组]New Doraemon 哆啦A梦新番[715][AVC][1080P][GB_JP]
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>.+?)[-_. ]+\(?(?:[-_. ]?#?(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+)))+\)?(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Multi-episode Repeated (S01E05 - S01E06, 1x05 - 1x06, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:(?:[ex]|[-_. ]e){1,2}(?<episode>\d{1,3}(?!\d+)))+){2,}`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Single episodes with a title (S01E05, 1x05, etc) and trailing info in slashes
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))(?:[ex]|\W[ex]|_){1,2}(?<episode>\d{2,3}(?!\d+|(?:[ex]|\W[ex]|_|-){1,2}\d+))).+?(?:\[.+?\])(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title Season EpisodeNumber + Absolute Episode Number [SubGroup]
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:[-_\W](?<![()\[!]))+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:[ex]|\W[ex]){1,2}(?<episode>(?<!\d+)\d{2}(?!\d+)))).+?(?:[-_. ]?(?<absoluteepisode>(?<!\d+)\d{3}(\.\d{1,2})?(?!\d+)))+.+?\[(?<subgroup>.+?)\](?:$|\.mkv)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Multi-Episode with a title (S01E05E06, S01E05-06, S01E05 E06, etc) and trailing info in slashes
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))(?:[ex]|\W[ex]|_){1,2}(?<episode>\d{2,3}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{2,3}(?!\d+)))+).+?(?:\[.+?\])(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title Absolute Episode Number [SubGroup] [Hash]? (Series Title Episode 99-100 [RlsGroup] [ABCD1234])
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)[-_. ]Episode(?:[-_. ]+(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+)))+(?:.+?)\[(?<subgroup>.+?)\].*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title Absolute Episode Number [SubGroup] [Hash]
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:_|-|\s|\.)+(?<absoluteepisode>\d{3}(\.\d{1,2})(?!\d+)))+(?:.+?)\[(?<subgroup>.+?)\].*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title Absolute Episode Number [Hash]
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:_|-|\s|\.)+(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+)))+(?:[-_. ]+(?<special>special|ova|ovd))?[-_. ]+.*?(?<hash>\[\w{8}\])(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with airdate AND season/episode number, capture season/epsiode only
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airdate>\d{4}\W+[0-1][0-9]\W+[0-3][0-9])(?!\W+[0-3][0-9])[-_. ](?:s?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+)))(?:[ex](?<episode>(?<!\d+)(?:\d{1,3})(?!\d+)))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with airdate AND season/episode number
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airyear>\d{4})\W+(?<airmonth>[0-1][0-9])\W+(?<airday>[0-3][0-9])(?!\W+[0-3][0-9]).+?(?:s?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+)))(?:[ex](?<episode>(?<!\d+)(?:\d{1,3})(?!\d+)))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	// Multi-episode with title (S01E05-06, S01E05-6)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:[-_\W](?<![()\[!]))+S(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))E(?<episode>\d{1,2}(?!\d+))(?:-(?<episode>\d{1,2}(?!\d+)))+[-_. ]`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with a title, Single episodes (S01E05, 1x05, etc) & Multi-episode (S01E05E06, S01E05-06, S01E05 E06, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))(?:[ex]|\W[ex]){1,2}(?<episode>\d{2,3}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{2,3}(?!\d+)))*)\W?(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with a title, 4 digit season number, Single episodes (S2016E05, etc) & Multi-episode (S2016E05E06, S2016E05-06, S2016E05 E06, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S(?<season>(?<!\d+)(?:\d{4})(?!\d+))(?:e|\We|_){1,2}(?<episode>\d{2,3}(?!\d+))(?:(?:\-|e|\We|_){1,2}(?<episode>\d{2,3}(?!\d+)))*)\W?(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with a title, 4 digit season number, Single episodes (2016x05, etc) & Multi-episode (2016x05x06, 2016x05-06, 2016x05 x06, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+(?<season>(?<!\d+)(?:\d{4})(?!\d+))(?:x|\Wx){1,2}(?<episode>\d{2,3}(?!\d+))(?:(?:\-|x|\Wx|_){1,2}(?<episode>\d{2,3}(?!\d+)))*)\W?(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	// Multi-season pack
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)[-_. ]+S(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))-S?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	// Partial season pack
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:\W+S(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))\W+(?:(?:Part\W?|(?<!\d+\W+)e)(?<seasonpart>\d{1,2}(?!\d+)))+)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Mini-Series with year in title, treated as season 1, episodes are labelled as Part01, Part 01, Part.1
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?\d{4})(?:\W+(?:(?:Part\W?|e)(?<episode>\d{1,2}(?!\d+)))+)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Mini-Series, treated as season 1, multi episodes are labelled as E1-E2
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:[-._ ][e])(?<episode>\d{2,3}(?!\d+))(?:(?:\-?[e])(?<episode>\d{2,3}(?!\d+)))+`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with airdate and part (2018.04.28.Part.2)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airyear>\d{4})[-_. ]+(?<airmonth>[0-1][0-9])[-_. ]+(?<airday>[0-3][0-9])(?![-_. ]+[0-3][0-9])[-_. ]+Part[-_. ]?(?<part>[1-9])`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Mini-Series, treated as season 1, episodes are labelled as Part01, Part 01, Part.1
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:\W+(?:(?:Part\W?|(?<!\d+\W+)e)(?<episode>\d{1,2}(?!\d+)))+)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Mini-Series, treated as season 1, episodes are labelled as Part One/Two/Three/...Nine, Part.One, Part_One
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:\W+(?:Part[-._ ](?<episode>One|Two|Three|Four|Five|Six|Seven|Eight|Nine)(?>[-._ ])))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Mini-Series, treated as season 1, episodes are labelled as XofY
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:\W+(?:(?<episode>(?<!\d+)\d{1,2}(?!\d+))of\d+)+)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Supports Season 01 Episode 03
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`(?:.*(?:\""|^))(?<title>.*?)(?:[-_\W](?<![()\[]))+(?:\W?Season\W?)(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:\W|_)+(?:Episode\W)(?:[-_. ]?(?<episode>(?<!\d+)\d{1,2}(?!\d+)))+`,
		regexp2.Compiled|regexp2.IgnoreCase))

	// Multi-episode with episodes in square brackets (Series Title [S01E11E12] or Series Title [S01E11-12])
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`(?:.*(?:^))(?<title>.*?)[-._ ]+\[S(?<season>(?<!\d+)\d{2}(?!\d+))(?:[E-]{1,2}(?<episode>(?<!\d+)\d{2}(?!\d+)))+\]`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Multi-episode release with no space between series title and season (S01E11E12)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`(?:.*(?:^))(?<title>.*?)S(?<season>(?<!\d+)\d{2}(?!\d+))(?:E(?<episode>(?<!\d+)\d{2}(?!\d+)))+`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Multi-episode with single episode numbers (S6.E1-E2, S6.E1E2, S6E1E2, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)[-_. ]S(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:[-_. ]?[ex]?(?<episode>(?<!\d+)\d{1,2}(?!\d+)))+`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Single episode season or episode S1E1 or S1-E1 or S1.Ep1 or S01.Ep.01
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`(?:.*(?:\""|^))(?<title>.*?)(?:\W?|_)S(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:\W|_)?Ep?[ ._]?(?<episode>(?<!\d+)\d{1,2}(?!\d+))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//3 digit season S010E05
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`(?:.*(?:\""|^))(?<title>.*?)(?:\W?|_)S(?<season>(?<!\d+)\d{3}(?!\d+))(?:\W|_)?E(?<episode>(?<!\d+)\d{1,2}(?!\d+))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//5 digit episode number with a title
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:(?<title>.+?)(?:_|-|\s|\.)+)(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+)))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>(?<!\d+)\d{5}(?!\d+)))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//5 digit multi-episode with a title
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:(?<title>.+?)(?:_|-|\s|\.)+)(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+)))(?:(?:[-_. ]{1,3}ep){1,2}(?<episode>(?<!\d+)\d{5}(?!\d+)))+`,
		regexp2.Compiled|regexp2.IgnoreCase))

	// Separated season and episode numbers S01 - E01
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:_|-|\s|\.)+S(?<season>\d{2}(?!\d+))(\W-\W)E(?<episode>(?<!\d+)\d{2}(?!\d+))(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	// Season and episode numbers in square brackets (single and mult-episode)
	// Series Title - [02x01] - Episode 1
	// Series Title - [02x01x02] - Episode 1
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?(?:[-_\W](?<![()\[!]))+\[(?<season>(?<!\d+)\d{1,2})(?:(?:-|x){1,2}(?<episode>\d{2}))+\].+?(?:\.|$)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	// Anime - Title with season number - Absolute Episode Number (Title S01 - EP14)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?S\d{1,2})[-_. ]{3,}(?:EP)?(?<absoluteepisode>\d{2,3}(\.\d{1,2})?(?!\d+|[-]))`,
		regexp2.Compiled|regexp2.IgnoreCase))
	//
	//// Anime - French titles with single episode numbers, with or without leading sub group ([RlsGroup] Title - Episode 1)
	//ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)[-_. ]+?(?:Episode[-_. ]+?)(?<absoluteepisode>\d{1}(\.\d{1,2})?(?!\d+))`,
	//	regexp2.Compiled|regexp2.IgnoreCase))

	// Anime - Absolute episode number in square brackets  [梦蓝字幕组]New Doraemon 哆啦A梦新番 [715][AVC][1080P][GB_JP] [-_. ]+? 不支持 `番 [715]` *? 支持
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)(?<title>.+?)[-_. ]+?[\[第](?<absoluteepisode>\d{1,4}(\.\d{1,2})?(?:v\d)?(?!\d+))[\]话集]`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//[DAY字幕组] 女神的露天咖啡厅 Megami no Café Terrace - 5 [简日双语][WEBrip][1080P][MP4]
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)(?<title>.+?)[-_. ]+?\s(?<absoluteepisode>\d{1,4}(\.\d{1,2})?(?:v\d)?(?!\d+))\s`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Season only releases
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)\W(?:S|Season|Saison)\W?(?<season>\d{1,2}(?!\d+))(\W+|_|$)(?<extras>EXTRAS|SUBPACK)?(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//4 digit season only releases
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)\W(?:S|Season|Saison)\W?(?<season>\d{4}(?!\d+))(\W+|_|$)(?<extras>EXTRAS|SUBPACK)?(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with a title and season/episode in square brackets
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+\[S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>(?<!\d+)\d{2}(?!\d+|i|p)))+\])\W?(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Supports 103/113 naming
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?(?:(?:[_.-](?<![()\[!]))+(?<season>(?<!\d+)[1-9])(?<episode>[1-9][0-9]|[0][1-9])(?![a-z]|\d+))+(?:[_.]|$)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//4 digit episode number
	//Episodes without a title, Single (S01E05, 1x05) AND Multi (S01E04E05, 1x04x05, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{4}(?!\d+|i|p)))+)(\W+|_|$)(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//4 digit episode number
	//Episodes with a title, Single episodes (S01E05, 1x05, etc) & Multi-episode (S01E05E06, S01E05-06, S01E05 E06, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{4}(?!\d+|i|p)))+)\W?(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with airdate (2018.04.28)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airyear>\d{4})[-_. ]+(?<airmonth>[0-1][0-9])[-_. ]+(?<airday>[0-3][0-9])(?![-_. ]+[0-3][0-9])`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with airdate (04.28.2018)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airmonth>[0-1][0-9])[-_. ]+(?<airday>[0-3][0-9])[-_. ]+(?<airyear>\d{4})(?!\d+)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with airdate (20180428)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?\W*(?<!\d+)(?<airyear>\d{4})(?<airmonth>[0-1][0-9])(?<airday>[0-3][0-9])(?!\d+)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Supports 1103/1113 naming
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)?(?:(?:[-_. ](?<![()\[!]))*(?<season>(?<!\d+|\(|\[|e|x)\d{2})(?<episode>(?<!e|x)\d{2}(?!p|i|\d+|\)|\]|\W\d+|\W(?:e|ep|x)\d+)))+([-_. ]+|$)(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Episodes with single digit episode number (S01E1, S01E5E6, etc)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.*?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]){1,2}(?<episode>\d{1}))+)+(\W+|_|$)(?!\\)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//iTunes Season 1\05 Title (Quality).ext
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:Season(?:_|-|\s|\.)(?<season>(?<!\d+)\d{1,2}(?!\d+)))(?:_|-|\s|\.)(?<episode>(?<!\d+)\d{1,2}(?!\d+))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//iTunes 1-05 Title (Quality).ext
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:(?<season>(?<!\d+)(?:\d{1,2})(?!\d+))(?:-(?<episode>\d{2,3}(?!\d+))))`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime Range - Title Absolute Episode Number (ep01-12)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:_|\s|\.)+(?:e|ep)(?<absoluteepisode>\d{2,3}(\.\d{1,2})?)-(?<absoluteepisode>(?<!\d+)\d{1,2}(\.\d{1,2})?(?!\d+|-)).*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title Absolute Episode Number (e66)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:(?:_|-|\s|\.)+(?:e|ep)(?<absoluteepisode>\d{2,4}(\.\d{1,2})?))+.*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title Episode Absolute Episode Number (Series Title Episode 01)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)[-_. ](?:Episode)(?:[-_. ]+(?<absoluteepisode>(?<!\d+)\d{2,3}(\.\d{1,2})?(?!\d+)))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime Range - Title Absolute Episode Number (1 or 2 digit absolute episode numbers in a range, 1-10)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)[_. ]+(?<absoluteepisode>(?<!\d+)\d{1,2}(\.\d{1,2})?(?!\d+))-(?<absoluteepisode>(?<!\d+)\d{1,2}(\.\d{1,2})?(?!\d+|-))(?:_|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title Absolute Episode Number
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:[-_. ]+(?<absoluteepisode>(?<!\d+)\d{2,3}(\.\d{1,2})?(?!\d+|[ip])))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Anime - Title {Absolute Episode Number}
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+(?<absoluteepisode>(?<!\d+)\d{2,3}(\.\d{1,2})?(?!\d+|[ip])))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//Extant, terrible multi-episode naming (extant.10708.hdtv-lol.mp4)
	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)[-_. ](?<season>[0]?\d?)(?:(?<episode>\d{2}){2}(?!\d+))[-_. ]`,
		regexp2.Compiled|regexp2.IgnoreCase))

	ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+(?<absoluteepisode>(?<!\d+)\d{2,3}(\.\d{1,2})?(?!\d+|[ip])))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.Compiled|regexp2.IgnoreCase))

	//[国漫]妖神记第253集无水印高清下载
	//ReportTitleRegex = append(ReportTitleRegex, regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+)第(?<episode>\d+)[话集](?:$|\.)?`,
	//	regexp2.Compiled|regexp2.IgnoreCase))
	//[幻樱字幕组][7月新番][继母的拖油瓶是我的前女友 Mamahaha no Tsurego ga Motokano Datta][01v2][GB_MP4][1280X720]
}

func (this *TorrentTitleParse) PreSubstitution(result *MatchResult) {

	for _, replace := range PreSubstitutionRegex {
		s, err := replace.Regex.Replace(result.AnalyzeTitle, replace.Replace, replace.StartAt, -1)
		if err != nil {
			return
		}
		result.AnalyzeTitle = s
	}
	result.AnalyzeTitle = helper.StrReplace(result.AnalyzeTitle, []string{
		"枫雪字幕_",
		"_1080P_",
	}, []string{
		"[枫雪字幕] ",
		" [1080P] ",
	})

	result.SimpleTitle, _ = SimpleTitleRegex.Regex.Replace(result.AnalyzeTitle, SimpleTitleRegex.Replace, SimpleTitleRegex.StartAt, -1)
	stReg := ReplaceOption{
		Regex:   regexp2.MustCompile(`(\[.*?\])`, regexp2.Compiled|regexp2.IgnoreCase),
		Replace: " $1 ",
		StartAt: 1,
	}
	tmp, err := stReg.Regex.Replace(result.SimpleTitle, stReg.Replace, stReg.StartAt, 1)
	if err == nil && tmp != "" {
		result.SimpleTitle = tmp
	}

	result.SimpleTitle = helper.StrReplace(result.SimpleTitle, []string{` [] `, "[]"}, []string{"[]", ""})

}
