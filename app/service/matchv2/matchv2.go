package matchv2

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/util/hash"
	"XArr-Rss/util/logsys"
	"errors"
	"github.com/dlclark/regexp2"
	"github.com/spf13/cast"
	"net/url"
	"strconv"
	"strings"
)

type MatchV2 struct {
}

// 解析数据
func (this MatchV2) Parse(groupMedia *dbmodel.GroupMedia, sourceItems []dbmodel.SourceItem, groupCacheChannel *model.RssResult) error {

	// 判断XmlItem是否传递
	if groupCacheChannel == nil {
		return errors.New("分组缓存异常")
	}

	for k, _ := range sourceItems {
		// 正对souceItem格式化
		sourceItems[k].Enclosure.Length = this.ParseSourceItemLength(sourceItems[k].Enclosure.Length)
	}

	// 开始查询符合规则的数据
	for _, regStr := range groupMedia.Regex {
		// 格式化数据项
		err, resultList := this.ParseSourceItem(groupMedia, regStr, sourceItems)
		if err != nil {
			if strings.Contains(err.Error(), "非赞助会员") {
				return err
			}
			if strings.Contains(err.Error(), "正则匹配数量错误") {
				return err
			}
			continue
		}
		if resultList != nil {
			for _, v := range resultList {
				v.Enclosure.Length = this.ParseSourceItemLength(v.Enclosure.Length)
				v.XArrRssIndexer.Text = groupCacheChannel.Title + v.XArrRssIndexer.Text
				v.XArrRssIndexer.ID = strconv.Itoa(int(groupMedia.Id))
				groupCacheChannel.Item = append(groupCacheChannel.Item, v)
			}

		}
	}
	return nil
}

// ParseSourceItem 解析数据源中的媒体数据
// @var groupMedia 分组里面的媒体数据
// @var groupRegex 分组中的配置方式
// @var sourceItems 需要匹配的内容
func (this MatchV2) ParseSourceItem(groupMedia *dbmodel.GroupMedia, groupRegex dbmodel.GroupRegex, sourceItems []dbmodel.SourceItem) (err error, res []model.RssResultItem) {
	err, manualReg := this.checkAndGenReg(groupRegex)
	if err != nil {
		return err, nil
	}

	// 如果不为自动并且含有错误就返回
	if groupRegex.MatchType != "auto" && err != nil {
		return logsys.Error("媒体正则表达式异常了吧:%s", "正则", err.Error()), nil
	}

	// 查询媒体信息
	sonarrMediaInfo := groupMedia.MediaInfo

	// 定义返回的数据项
	groupCache := model.RssRoot{}

	// 数据源开始一个一个匹配
	for _, sourceItem := range sourceItems {
		// 先匹配标题是否涵盖出来需要的内容
		torrentTileParse := &torrent_title_parse.MatchResult{
			OldTitle:     sourceItem.Title,
			AnalyzeTitle: sourceItem.Title,
			AudioEncode:  []string{},
			VideoEncode:  []string{},
		}

		if sourceItem.ParseInfo != nil {
			torrentTileParse = &torrent_title_parse.MatchResult{
				OldTitle:     sourceItem.Title,
				AnalyzeTitle: sourceItem.ParseInfo.AnalyzeTitle,
				AudioEncode:  sourceItem.ParseInfo.AudioEncode,
				VideoEncode:  sourceItem.ParseInfo.VideoEncode,
				MediaType:    sourceItem.ParseInfo.MediaType,
				MinEpisode:   sourceItem.ParseInfo.MinEpisode,
				MaxEpisode:   sourceItem.ParseInfo.MaxEpisode,
				Version:      sourceItem.ParseInfo.Version,
				Season:       sourceItem.ParseInfo.Season,
				ReleaseGroup: sourceItem.ParseInfo.ReleaseGroup,
				//QualitySource:      sourceItem.ParseInfo.QualitySource,
				//Resolution:         sourceItem.ParseInfo.Resolution,
				QualityResolution:  sourceItem.ParseInfo.QualityResolution,
				Language:           sourceItem.ParseInfo.Language,
				ProductionCompany:  sourceItem.ParseInfo.ProductionCompany,
				Subtitles:          sourceItem.ParseInfo.Subtitles,
				AbsoluteMinEpisode: sourceItem.ParseInfo.MinEpisode,
				//AbsoluteMaxEpisode: sourceItem.ParseInfo.AbsoluteMaxEpisode,
			}
		} else if groupRegex.MatchType == "auto" {
			logsys.Debug("需要重新解析 或 不支持:%s", "匹配标题", sourceItem.Title)
			//parse := torrent_title_parse.TorrentTitleParse{}
			//torrentTileParse = parse.Parse(sourceItem.Title)
			// 解析不了的 不录入
			continue
		}

		// 解析标题中的信息
		mediaParseInfo := &torrent_title_parse.MatchResult{}
		if groupRegex.MatchType == "auto" {
			if !variable.ServerState.IsVip {
				return logsys.Error("赞助会员才能使用自动匹配功能", "匹配"), nil

			}
			if !this.matchRegTitle(groupRegex.Reg, groupRegex.RegType, sourceItem.Title) {
				continue
			}

			// 自动必然是支持会员的
			if torrentTileParse.MinEpisode == 0 {
				// 匹配错误
				//logsys.Debug("匹配错误:%s %s", "自动匹配", sourceItem.Title.Text, torrentTileParse.AnalyzeTitle)
				//return errors.New("匹配错误"), nil
				continue
			}
			// 自动检测标题是否匹配
			err = this.parseTitleIsOk(sonarrMediaInfo, torrentTileParse, groupRegex)
			if err == nil {
				// 判断是否需要自动搜索语言
				if groupMedia.Language != "-1" {
					torrentTileParse.Language = groupMedia.Language
				}
				// 自动搜索质量
				if groupMedia.Quality != "-1" {
					torrentTileParse.QualityResolution = groupMedia.Quality
				}

				mediaParseInfo = torrentTileParse

				// 只有在Season为0并且不是特殊指定要S00的情况下，才考虑将其设置为1
				if mediaParseInfo.Season == 0 && groupRegex.Season != 0 {
					if sonarrMediaInfo.SeasonCount <= 1 {
						mediaParseInfo.Season = 1
					}
				}
			}
		} else {
			// 手动 支持 episode max_episode season 分组命名
			torrentTileParseTmp := &torrent_title_parse.MatchResult{}
			err = this.parseManual(groupRegex, sonarrMediaInfo, sourceItem, manualReg, torrentTileParseTmp)
			if err == nil {
				if groupRegex.Season == -1 {
					torrentTileParseTmp.Season = torrentTileParse.Season
				}
				this.ParseMediaSeasonAndEpisode(torrentTileParseTmp, groupRegex.Season, sonarrMediaInfo)
				// 自动检测标题是否匹配
				//err = this.parseTitleIsOk(sonarrMediaInfo, torrentTileParseTmp, groupRegex)
				//if err == nil {
				// 判断是否需要自动搜索语言
				if groupMedia.Language == "-1" {
					torrentTileParseTmp.Language = torrentTileParse.Language
				} else {
					torrentTileParseTmp.Language = groupMedia.Language
				}
				// 自动搜索质量
				if groupMedia.Quality == "-1" {
					torrentTileParseTmp.QualityResolution = torrentTileParse.QualityResolution
				} else {
					torrentTileParseTmp.QualityResolution = groupMedia.Quality
				}

				mediaParseInfo = torrentTileParseTmp
				//}
			}

		}

		// 没有匹配到信息就不执行了
		if err != nil {
			if strings.Contains(err.Error(), "非赞助会员") {
				return err, nil
			}
			if strings.Contains(err.Error(), "正则匹配数量错误") {
				return err, nil
			}
			continue
		}

		// 判断语言是否未找到

		var rssItemTitle = ""
		// 替换为Sonarr能识别的标题
		rssItemTitle = this.getTitle(groupMedia, mediaParseInfo, groupRegex.Offset)

		if mediaParseInfo.MaxEpisode == 0 {
			mediaParseInfo.MaxEpisode = mediaParseInfo.MinEpisode
		}
		if mediaParseInfo.AbsoluteMinEpisode == 0 {
			logsys.Info("t", "f")
		}

		// 计算整季多集长度
		if mediaParseInfo.Season >= 0 && mediaParseInfo.MinEpisode == 0 && mediaParseInfo.MaxEpisode == 0 {
			if sourceItem.Enclosure.Length == "996973766" {
				for _, season := range sonarrMediaInfo.Seasons {
					if season.SeasonNumber == mediaParseInfo.Season {
						sourceItem.Enclosure.Length = strconv.Itoa(996973766 * season.Statistics.TotalEpisodeCount)
					}
				}
			}
		}

		// 计算连续剧集文件长度
		if mediaParseInfo.Season >= 0 && mediaParseInfo.MinEpisode > 0 && mediaParseInfo.MaxEpisode > mediaParseInfo.MinEpisode {
			if sourceItem.Enclosure.Length == "996973766" {
				for _, season := range sonarrMediaInfo.Seasons {
					if season.SeasonNumber == mediaParseInfo.Season {
						sourceItem.Enclosure.Length = strconv.Itoa(996973766 * (mediaParseInfo.MaxEpisode - mediaParseInfo.MinEpisode))
					}
				}
			}
		}

		// 如果没有长度则默认
		if sourceItem.Enclosure.Length == "" {
			sourceItem.Enclosure.Length = "996973766"
		}

		// 格式化到group
		groupCache.Channel.Item = append(groupCache.Channel.Item, model.RssResultItem{
			Title:         model.CDATA{Text: rssItemTitle},
			CnTitle:       sonarrMediaInfo.CnTitle,
			OriginalTitle: model.CDATA{Text: sourceItem.Title},
			OtherTitle:    model.CDATA{Text: " other:" + strings.Join(sonarrMediaInfo.AlternateTitles, "-|-")},
			PubDate:       sourceItem.PubDate,
			OldMinEpisode: mediaParseInfo.AbsoluteMinEpisode,
			OldMaxEpisode: mediaParseInfo.AbsoluteMaxEpisode,
			Season:        mediaParseInfo.Season,
			MinEpisode:    mediaParseInfo.MinEpisode + groupRegex.Offset,
			MaxEpisode:    mediaParseInfo.MaxEpisode + groupRegex.Offset,
			Enclosure: model.RssResultItemEnclosure{
				Type:   sourceItem.Enclosure.Type,
				Length: sourceItem.Enclosure.Length,
				Url:    appconf.AppConf.System.HttpAddr + "/api/v1/down?url=" + url.QueryEscape(sourceItem.Enclosure.Url) + "&title=" + url.QueryEscape(rssItemTitle) + "&source_id=" + url.QueryEscape(strconv.Itoa(sourceItem.SourceId)),
			},
			Link: sourceItem.Link,
			Guid: model.RssResultItemGuid{
				IsPermaLink: false,
				Text:        hash.Md5{}.HashString(sourceItem.Guid + sourceItem.Enclosure.Url + sourceItem.Enclosure.Length + cast.ToString(sourceItem.SourceId)),
			},
			XArrRssIndexer: model.RssResultItemXArrRssIndexer{
				Text: " 数据源:" + strconv.Itoa(sourceItem.SourceId),
				//Id:   groupMedia.Id,
			},
			OthderId: model.RssResultItemOtherId{
				TmdbId: groupMedia.MediaInfo.TmdbId,
				TvdbId: strconv.Itoa(groupMedia.MediaInfo.TvdbId),
				ImdbId: groupMedia.MediaInfo.ImdbId,
			},
			OtherInfo: *mediaParseInfo,
			SourceId:  strconv.Itoa(sourceItem.SourceId),
		})
	}
	return nil, groupCache.Channel.Item
}

// 匹配标题内容
func (this MatchV2) matchRegTitle(matchReq string, regType int, title string) bool {
	if matchReq != "" {
		if regType != dbmodel.REG_TYPE_REGEXP {
			// 过滤包含
			// 有匹配规则
			matchCount := 0
			matchAndArr := strings.Split(matchReq, ",")
			for _, matchAnd := range matchAndArr {
				// 判断是否有|
				if strings.Contains(matchAnd, "|") {
					matchOrArr := strings.Split(matchAnd, "|")
					matchOrCount := 0

					for _, matchOr := range matchOrArr {
						// 如果有一个实现效果 则跳出处理
						if strings.Contains(title, matchOr) {
							matchOrCount++
							break
						}
					}
					if matchOrCount == 0 {
						return false
					}
					matchCount++

				} else {
					// 单独一个
					if strings.Contains(title, matchAnd) {
						matchCount++
					} else {
						// 不包含指定字符串
						// 返回
						return false
					}
				}
			}
			if len(matchAndArr) == matchCount {
				return true
			}
			return false
		} else {
			reg, err := regexp2.Compile(matchReq, regexp2.IgnoreCase)
			if err != nil {
				logsys.Error("设置的匹配规则有问题:%s", "匹配规则", err.Error())
				return false
			}
			matchString, err := reg.MatchString(title)
			if err != nil {
				logsys.Error("设置的匹配规则有问题 匹配文字失败:%s", "匹配规则", err.Error())
				return false
			}
			return matchString
		}

	} else {
		// 不筛选返回
		return true
	}
}
