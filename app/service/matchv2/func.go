package matchv2

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/array"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"strings"
)

// 处理种子大小
func (this MatchV2) ParseSourceItemLength(length string) string {
	if length == "0" || length == "" || helper.StrToInt(length) < 1024 {
		length = "996973766"
	}
	return length
}

// 检测和生成正则表达式
func (this MatchV2) checkAndGenReg(groupRegex dbmodel.GroupRegex) (error, *regexp2.Regexp) {
	if groupRegex.Reg == "" && groupRegex.MatchType != "auto" {
		return errors.New("正则表达式匹配规则为空"), nil
	}

	// 初始化正则方法
	reg, err := regexp2.Compile(strings.ReplaceAll(groupRegex.Reg, "(?P<", "(?<"), regexp2.IgnoreCase)
	// 判断正则是否正确
	if err != nil {
		return logsys.Error("正则数据错误:%s err：%v", "数据源同步至分组", groupRegex.Reg, err.Error()), nil
	}
	if reg == nil {
		return logsys.Error("正则数据错误：%v", "数据源同步至分组", groupRegex.Reg), nil
	}
	return nil, reg
}

// 将sonarr标题合并
func (this MatchV2) JoinSonarrTitle(sonarrMediaInfo dbmodel.Media) []string {
	//var b strings.Builder
	////b.Grow()
	//b.WriteString(sonarrMediaInfo.OriginalTitle)
	//b.WriteString("1666666666663")
	//b.WriteString(sonarrMediaInfo.CnTitle)
	//b.WriteString("1666666666663")
	//for _, itm := range sonarrMediaInfo.Titles {
	//	b.WriteString(itm)
	//	b.WriteString("1666666666663")
	//}
	//for _, itm := range sonarrMediaInfo.AlternateTitles {
	//	b.WriteString(itm)
	//	b.WriteString("1666666666663")
	//}

	a := append([]string{},
		sonarrMediaInfo.OriginalTitle,
		sonarrMediaInfo.CnTitle,
	)
	a = append(a, sonarrMediaInfo.Titles...)
	a = append(a, sonarrMediaInfo.AlternateTitles...)
	v := helper.ReplaceRegString(strings.Join(a, "1666666666663"))
	v = strings.ToLower(v)
	a = strings.Split(v, "1666666666663")
	v = ""
	a = array.UniqueString(a)
	a = array.FilterString(a)
	return a
	//
	//title = sonarrMediaInfo.OriginalTitle
	//if sonarrMediaInfo.CnTitle != "" {
	//	title += "-|-" + sonarrMediaInfo.CnTitle
	//}
	//if len(sonarrMediaInfo.Titles) > 0 {
	//	title += "-|-" + strings.Join(sonarrMediaInfo.Titles, "-|-")
	//}
	//if len(sonarrMediaInfo.AlternateTitles) > 0 {
	//	title += "-|-" + strings.Join(sonarrMediaInfo.AlternateTitles, "-|-")
	//}
	//return title

}

// 生成标题
func (this MatchV2) getTitle(groupMedia *dbmodel.GroupMedia, mediaParseInfo *torrent_title_parse.MatchResult, offset int) string {
	mediaInfo := groupMedia.MediaInfo
	mediaParseInfo.AbsoluteMaxEpisode = mediaParseInfo.AbsoluteMinEpisode
	if mediaParseInfo.MaxEpisode > mediaParseInfo.MinEpisode && mediaParseInfo.MaxEpisode > 0 {
		// 计算相差集数
		mediaParseInfo.AbsoluteMaxEpisode = mediaParseInfo.AbsoluteMinEpisode + (mediaParseInfo.MaxEpisode - mediaParseInfo.MinEpisode)
	}

	season := ""
	if mediaParseInfo.Season > 0 {
		season = fmt.Sprintf("S%02d", mediaParseInfo.Season) // 季
	}
	episode := ""
	abEpisode := ""
	if mediaParseInfo.MinEpisode > 0 {
		if mediaParseInfo.AbsoluteMaxEpisode > mediaParseInfo.AbsoluteMinEpisode {
			episode = fmt.Sprintf("E%02d-E%02d", mediaParseInfo.MinEpisode+offset, mediaParseInfo.MaxEpisode+offset)
			//title = append(title, season+episode) // 连续集
			// 获取绝对集
			//abEpisode = fmt.Sprintf("Episode %02d-%02d", mediaParseInfo.AbsoluteMinEpisode, mediaParseInfo.AbsoluteMaxEpisode)
			abEpisode = fmt.Sprintf("%02d-%02d", mediaParseInfo.AbsoluteMinEpisode+offset, mediaParseInfo.AbsoluteMaxEpisode+offset)
			//title = append(title, abEpisode)
		} else {
			episode = fmt.Sprintf("E%02d", mediaParseInfo.MinEpisode+offset)
			//title = append(title, season+episode) // 集
			// 获取绝对集
			if mediaParseInfo.AbsoluteMinEpisode > 0 {
				//abEpisode = fmt.Sprintf("Episode %02d", mediaParseInfo.AbsoluteMinEpisode)
				abEpisode = fmt.Sprintf("%02d", mediaParseInfo.AbsoluteMinEpisode+offset)
				//title = append(title, abEpisode)
			}
		}

	}

	releaseGroup := mediaParseInfo.ReleaseGroup
	if releaseGroup == "" {
		releaseGroup = "XArr"
	}
	retTitle := ""
	if mediaInfo.SeriesType == "anime" {
		if groupMedia.EchoTitleAnime == "" || !variable.ServerState.IsVip {
			if appconf.AppConf.System.EchoTitleAnime == "" {
				retTitle = `[{releaseGroup}][{chineseTitle}] {title} - {season}{episode} ({abEpisode}) [{language}][{quality}][{video}][{audio}][{mediaType}]`
			} else {
				retTitle = appconf.AppConf.System.EchoTitleAnime
			}
		} else {
			retTitle = groupMedia.EchoTitleAnime
		}

	} else {
		if groupMedia.EchoTitleTv == "" || !variable.ServerState.IsVip {
			if appconf.AppConf.System.EchoTitleTv == "" {
				retTitle = `[{releaseGroup}][{chineseTitle}] {title} - {season}{episode} [{language}][{quality}][{video}][{audio}][{mediaType}]`
			} else {
				retTitle = appconf.AppConf.System.EchoTitleTv
			}
		} else {
			retTitle = groupMedia.EchoTitleTv
		}
	}
	retTitle = helper.StrReplace(retTitle, []string{
		`{releaseGroup}`,
		`{chineseTitle}`,
		`{title}`,
		`{season}`,
		`{episode}`,
		`{abEpisode}`,
		`{language}`,
		`{quality}`,
		`{video}`,
		`{audio}`,
		`{mediaType}`,
		"[]",
		"()",
	}, []string{
		releaseGroup,
		mediaInfo.CnTitle,
		mediaInfo.OriginalTitle,
		season,
		episode,
		abEpisode,
		mediaParseInfo.Language,
		mediaParseInfo.QualityResolution,
		strings.ToUpper(strings.Join(mediaParseInfo.VideoEncode, " ")),
		strings.ToUpper(strings.Join(mediaParseInfo.AudioEncode, " ")),
		strings.ToUpper(strings.Join(mediaParseInfo.MediaType, " ")),
		"",
		"",
	})
	//logsys.Debug(retTitle, "匹配标题")
	return retTitle
	//return strings.Join(title, "") + " " + other
}
