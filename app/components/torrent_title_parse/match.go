package torrent_title_parse

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"XArr-Rss/util/opencc_ext"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"strings"
)

type TorrentTitleParse struct {
}

type MatchResult struct {
	OldTitle           string   `json:"old_title,omitempty"`          // 原始标题
	AnalyzeTitle       string   `json:"analyze_title,omitempty"`      // 解析后的标题
	SimpleTitle        string   `json:"simple_title,omitempty"`       // 简化标题
	AudioEncode        []string `json:"audio_encode,omitempty"`       // 音频编码
	VideoEncode        []string `json:"video_encode,omitempty"`       // 视频编码
	MediaType          []string `json:"media_type"`                   // 媒体类型
	MinEpisode         int      `json:"min_episode,omitempty"`        // 最小季数 或 单集集数
	MaxEpisode         int      `json:"max_episode,omitempty"`        // 最大集数
	Version            int      `json:"version,omitempty"`            // 版本号
	Season             int      `json:"season,omitempty"`             // 季数
	Title              string   `json:"title,omitempty"`              // 标题
	ReleaseGroup       string   `json:"release_group,omitempty"`      // 发布组 [NC_RAW]
	QualitySource      string   `json:"quality_source,omitempty"`     // 质量 webdl
	Resolution         string   `json:"resolution,omitempty"`         // 尺寸 1080p
	QualityResolution  string   `json:"quality_resolution,omitempty"` // 质量_尺寸 Webdl-1080p
	Language           string   `json:"language"`                     // 语言 chinese
	ProductionCompany  string   `json:"production_company"`           // 发布 公司
	Subtitles          string   `json:"subtitles"`                    // 字幕
	AbsoluteMinEpisode int      `json:"absolute_min_episode"`         // 原始最小集 绝对集
	AbsoluteMaxEpisode int      `json:"absolute_max_episode"`         // 原始最大集 绝对集
	OtherTag           []string `json:"other_tag"`                    // 其他的tag标签
}

func (this *TorrentTitleParse) Parse(title string) *MatchResult {
	// 自定义预处理词
	title = this.ReplaceTitleSelfRule(title, appconf.AppConf.System.WordsRule)

	result := &MatchResult{
		OldTitle:     title,
		AnalyzeTitle: title,
		AudioEncode:  []string{},
		VideoEncode:  []string{},
	}
	result.AnalyzeTitle = strings.Trim(result.AnalyzeTitle, " ")
	// 预处理标题
	result.AnalyzeTitle = helper.StrReplace(result.AnalyzeTitle, []string{"【", "】", "（", "）"}, []string{"[", "]", "(", ")"})
	result.AnalyzeTitle = helper.StrReplace(result.AnalyzeTitle, []string{"[email protected]"}, []string{" "})

	result.SimpleTitle = result.AnalyzeTitle
	// 提前格式化标题
	this.PreSubstitution(result)

	for _, report := range ReportTitleRegex {
		match, err := report.MatchString(result.SimpleTitle)
		if err == nil && match {
			// 解析质量 1080p
			this.ParseQuality(result)

			// 解析格式
			result.MediaType = this.ParseMediaType(result)

			// 解析字幕组
			this.ParseReleaseGroup(result)

			// 解析其他标签信息
			this.ParseOtherTag(result)
			// 解析语言
			result.Language = this.ParseLanguage(result)

			if len(result.MediaType) == 0 {
				// 解析格式
				result.MediaType = this.ParseMediaType(result)
			}
			// 解析出来hevc acc
			this.ParseHevcAcc(result)

			// 解析集数
			this.ParseEpisode(result, report)
			// 解析季数
			this.ParseSeason(result)

			// 解析发布网
			this.ParseTvName(result)

			// 解析字幕信息
			this.ParseSubtitles(result)

			// 修改空数据
			spaceReplaceReg := regexp2.MustCompile(`(?:\[\s*?\]|\(\s*?\))`, regexp2.Compiled|regexp2.IgnoreCase)
			replace, err := spaceReplaceReg.Replace(result.AnalyzeTitle, "", -1, -1)
			if err != nil {
				return nil
			}
			result.AnalyzeTitle = replace
			result.AnalyzeTitle, _ = opencc_ext.T2s.Convert(result.AnalyzeTitle)
			return result
		}
	}

	return result
}

func (this *TorrentTitleParse) ReplaceTitleSelfRule(title, rules string) string {
	// 获取规则列表
	ruleArr := strings.Split(rules, "\n")
	for _, rule := range ruleArr {
		spaceRule := strings.TrimSpace(rule)
		if spaceRule == "" || strings.HasPrefix(spaceRule, "#") {
			continue
		}
		// 判断是否正确
		ruleList := strings.Split(rule, " => ")
		if len(ruleList) < 2 {
			logsys.Error("规则格式错误[%s]", rule)
			return title
		}

		if len(ruleList) == 2 {
			title = strings.ReplaceAll(title, ruleList[0], ruleList[1])
			continue
		}

		// 判断是否为正则
		if strings.TrimSpace(ruleList[0]) == "{regexp}" {
			match, err := regexp2.Compile(ruleList[1], regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)
			if err != nil {
				logsys.Error("校验自定义规则异常[%s]:%s", rule, err.Error())
				return title
			}

			// 正则替换
			title2, err := match.Replace(title, ruleList[2], -1, -1)
			if err == nil {
				title = title2
			}
			continue
		}

	}

	return title
}

func (this *TorrentTitleParse) ValidateRules(rules string) error {
	// 获取规则列表
	ruleArr := strings.Split(rules, "\n")
	for _, rule := range ruleArr {
		if strings.TrimSpace(rule) == "" || strings.HasPrefix(rule, "#") {
			continue
		}
		// 判断是否正确
		ruleList := strings.Split(rule, " => ")
		if len(ruleList) < 2 {
			return errors.New(fmt.Sprintf("规则格式错误[%s]", rule))
		}

		// 判断是否为正则
		if ruleList[0] == "{regexp}" {
			_, err := regexp2.Compile(ruleList[1], regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)
			if err != nil {
				return errors.New(fmt.Sprintf("校验自定义规则异常[%s]:%s", rule, err.Error()))
			}
		}
	}

	return nil
}
