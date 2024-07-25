package matchv2

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"XArr-Rss/util/regexp_ext"
	"errors"
	"github.com/dlclark/regexp2"
)

func (this MatchV2) parseManual(groupRegex dbmodel.GroupRegex, sonarrMediaInfo dbmodel.Media, sourceItem dbmodel.SourceItem, manualReg *regexp2.Regexp, mediaParseInfo *torrent_title_parse.MatchResult) error {

	matchString, _ := manualReg.MatchString(sourceItem.Title)
	if matchString {
		// 匹配成功
		//matchGroups := regex_group.GetRegxGroupByReg(manualReg, sourceItem.Title)
		matchGroups := regexp_ext.ParseGroups(manualReg, sourceItem.Title)
		if matchGroups == nil || len(matchGroups.Groups) == 0 {
			// 无匹配内容
			return errors.New("没有匹配到内容,请检查是否有设置命名分组")
		}
		// 获取集标记
		//episodeGroup, episodeOk := matchGroups[0]["episode"]
		episodeGroup, episodeOk := matchGroups.Groups["episode"]

		seasonGroup, seasonOk := matchGroups.Groups["season"]
		if variable.ServerState.IsVip {
			// 获取季标记
			if seasonOk {
				mediaParseInfo.Season = helper.StrToInt(seasonGroup[0])
			}
		}

		// 如果两个都没有 则返回错误
		if !episodeOk && !seasonOk {
			return logsys.Error("正则匹配数量错误,请检查是否有(?<episode>\\d+) 或者  (?<season>\\d+)：%v", "数据源同步至分组", groupRegex.Reg)
		}
		mediaParseInfo.MinEpisode = helper.StrToInt(episodeGroup[0])

		if variable.ServerState.IsVip {
			// 获取最大集标记
			maxEpisodeGroup, ok := matchGroups.Groups["max_episode"]
			if ok {
				mediaParseInfo.MaxEpisode = helper.StrToInt(maxEpisodeGroup[0])
			}
			// 增加发布组
			subgroup, ok := matchGroups.Groups["subgroup"]
			if ok {
				mediaParseInfo.ReleaseGroup = subgroup[0]
			}
		}
		return nil

	}

	return errors.New("没有匹配到")
}
