package matchv2

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"errors"
	"strconv"
	"strings"
)

func (this MatchV2) parseTitleIsOk(sonarrMediaInfo dbmodel.Media, torrentTitleParse *torrent_title_parse.MatchResult, groupRegex dbmodel.GroupRegex) error {
	// 提取所有sonarr标题
	// 取出sonarr title
	sonarrTitleArr := this.JoinSonarrTitle(sonarrMediaInfo)
	//sonarrTitleArr := sonarrMediaInfo.SearchTitle
	//sonarrTitleArr := strings.Split(sonarrMediaInfo.SearchTitle, "1666666666663")
	if len(sonarrTitleArr) == 0 {
		return logsys.Error("自动匹配模式未能找到对应Sonarr数据可能被删除了,SonarrId:"+strconv.Itoa(sonarrMediaInfo.SonarrId), "分组数据自动匹配")
	}

	// 直接筛选标题
	pp := 0
	maxPP := 1
	title := (torrentTitleParse.AnalyzeTitle)

	for _, v := range sonarrTitleArr {

		// 小于4个字符的不处理
		if v == "" || len(v) <= 4 {
			continue
		}

		// 模糊匹配
		hasMatch, matchString := helper.MatchString(v, title)
		if hasMatch {
			title = strings.Replace(title, matchString, "", -1)
			pp++
			continue
		}

	}

	if pp >= maxPP {
		// 匹配标题成功
		this.ParseMediaSeasonAndEpisode(torrentTitleParse, groupRegex.Season, sonarrMediaInfo)
		return nil
	} else {
		return errors.New("未能匹配成功")
	}

	//var err error
	//var parseData bool // *torrent_title_parse.MatchResult
	//
	//err, parseData = match.ParseByAutoReg(torrentTitleParse.AnalyzeTitle, sonarrTitle, groupRegex.Reg)
	//if err != nil {
	//	if !strings.Contains(err.Error(), "未能匹配成功") {
	//		logsys.Error("自动匹配模式错误:%s", "分组数据自动匹配", err.Error())
	//	}
	//	return err
	//} else {
	//	if parseData != false {
	//		this.ParseMediaSeasonAndEpisode(torrentTitleParse, groupRegex.Season, sonarrMediaInfo)
	//		return nil
	//	}
	//
	//	return errors.New("没有匹配到")
	//}
}

// 计算媒体中对应的 季 和 集
func (this MatchV2) ParseMediaSeasonAndEpisode(torrentTitleParse *torrent_title_parse.MatchResult, manualSeason int, sonarrMediaInfo dbmodel.Media) {
	if manualSeason == -1 {
		// 需要自动搜索

		// 判断是否已经找到对应剧
		if torrentTitleParse.Season >= 0 {
			// 找到季信息 则数据不用处理
			infoDb := medias.MediaEpisodeService{}.GetSeasonEpisode(sonarrMediaInfo.SonarrId, torrentTitleParse.Season, torrentTitleParse.MinEpisode)
			if infoDb != nil && infoDb.SeasonNumber >= 0 {
				// 找到数据库里面的绝对集信息
				torrentTitleParse.AbsoluteMinEpisode = infoDb.AbsoluteEpisodeNumber
			} else {
				// 数据找不到绝对集信息
				// 获取媒体的季数最大集
				searched := false
				for _, seasonInfo := range sonarrMediaInfo.Seasons {
					if seasonInfo.SeasonNumber == torrentTitleParse.Season {
						// 找到这个相同季
						if torrentTitleParse.MinEpisode > seasonInfo.Statistics.EpisodeCount {
							// 超过了当前季的总集数 则 按照绝对集进行查询
							infoDb := medias.MediaEpisodeService{}.GetSeasonAbEpisode(sonarrMediaInfo.SonarrId, torrentTitleParse.Season, torrentTitleParse.MinEpisode)
							if infoDb != nil && infoDb.SeasonNumber >= 0 {
								// 找到数据库里面的绝对集信息
								torrentTitleParse.MinEpisode = infoDb.EpisodeNumber
								torrentTitleParse.AbsoluteMinEpisode = infoDb.AbsoluteEpisodeNumber
								searched = true
							}
						}

						break
					}
				}

				if searched == false {
					torrentTitleParse.AbsoluteMinEpisode = torrentTitleParse.MinEpisode
					torrentTitleParse.AbsoluteMaxEpisode = torrentTitleParse.AbsoluteMinEpisode
				}
			}

		} else {
			// 未找到季信息 只有单集信息 视为绝对集数
			// torrentTitleParse.MinEpisode

			//////////////////////////////////////////////
			// 判断数据库中是否有绝对集信息
			seasonNumber := 0
			var infoDb *dbmodel.MediaEpisodeList
			if len(sonarrMediaInfo.Seasons) > 0 {
				seasonNumber = sonarrMediaInfo.Seasons[len(sonarrMediaInfo.Seasons)-1].SeasonNumber
				infoDb = medias.MediaEpisodeService{}.GetSeasonAbEpisode(sonarrMediaInfo.SonarrId, seasonNumber, torrentTitleParse.MinEpisode)
			} else {
				logsys.Error("异常媒体:%d %s %s 没有剧集", "auto", sonarrMediaInfo.SonarrId, sonarrMediaInfo.OriginalTitle, sonarrMediaInfo.CnTitle)
			}

			if infoDb != nil && infoDb.SeasonNumber >= 0 {
				// 找到数据库中绝对集信息
				torrentTitleParse.Season = infoDb.SeasonNumber
				if torrentTitleParse.MaxEpisode > 0 && torrentTitleParse.MinEpisode > 0 {
					torrentTitleParse.MaxEpisode = torrentTitleParse.MaxEpisode - (torrentTitleParse.MinEpisode - infoDb.EpisodeNumber)
				}
				torrentTitleParse.MinEpisode = infoDb.EpisodeNumber
				torrentTitleParse.AbsoluteMinEpisode = infoDb.AbsoluteEpisodeNumber
			} else {
				// 找不到season季信息
				if sonarrMediaInfo.SeasonCount == 1 {
					torrentTitleParse.Season = 1
				}
				torrentTitleParse.AbsoluteMinEpisode = torrentTitleParse.MinEpisode
				torrentTitleParse.AbsoluteMaxEpisode = torrentTitleParse.AbsoluteMinEpisode
			}
			//////////////////////////////////////////////

		}

	} else {
		// 使用自定义季
		torrentTitleParse.Season = manualSeason
		torrentTitleParse.AbsoluteMinEpisode = torrentTitleParse.MinEpisode
		torrentTitleParse.AbsoluteMaxEpisode = torrentTitleParse.AbsoluteMinEpisode
		// 去找绝对集
		if torrentTitleParse.Season >= 0 {
			// 找到季信息 则数据不用处理
			infoDb := medias.MediaEpisodeService{}.GetSeasonEpisode(sonarrMediaInfo.SonarrId, torrentTitleParse.Season, torrentTitleParse.MinEpisode)
			if infoDb != nil && infoDb.SeasonNumber >= 0 {
				// 找到数据库里面的绝对集信息
				torrentTitleParse.AbsoluteMinEpisode = infoDb.AbsoluteEpisodeNumber
				torrentTitleParse.AbsoluteMaxEpisode = torrentTitleParse.AbsoluteMinEpisode
			}
		}

	}
}

// 计算媒体第几季
func (this MatchV2) ParseMediaSeason(minEpisode, parseedSeason, season int, sonarrMediaInfo dbmodel.Media) (int, int) {
	if season == -1 {
		if !variable.ServerState.IsVip {
			logsys.Error("只非赞助会员用户不可以使用自动匹配季", "会员")
			return season, minEpisode
		}

		// 判断是否已经匹配到结果了
		season = parseedSeason
		if season >= 0 {
			return season, minEpisode
		}

		// 处理判断是否为多季 排除0季
		isMulite := 0
		seasonFirstEpisode := 0
		for _, v3 := range sonarrMediaInfo.Seasons {
			if v3.SeasonNumber > 0 {
				isMulite++
				if v3.SeasonNumber == 1 {
					seasonFirstEpisode = v3.Statistics.TotalEpisodeCount
				}
			}
		}
		// 如果为多季 并且 当前标题中的集数 小于第一季的总集数 则直接返回最新季和对应匹配的集数
		if isMulite > 1 && seasonFirstEpisode >= minEpisode {
			// 多季
			return sonarrMediaInfo.Seasons[len(sonarrMediaInfo.Seasons)-1].SeasonNumber, minEpisode
		}

		//////////////////////////////////////////////
		// 判断数据库中是否有绝对集信息
		infoDb := medias.MediaEpisodeService{}.GetSeasonAbEpisode(sonarrMediaInfo.SonarrId, sonarrMediaInfo.Seasons[len(sonarrMediaInfo.Seasons)-1].SeasonNumber, minEpisode)
		if infoDb != nil {
			// 找到绝对集信息
			return infoDb.SeasonNumber, infoDb.EpisodeNumber
		}

		//////////////////////////////////////////////

		// 根据每季集数 计算对应季
		tempSeansonEpisodeCount := 0
		for _, v2 := range sonarrMediaInfo.Seasons {
			if v2.SeasonNumber > 0 {
				// 将每季的集数相加 如果集数匹配 则为对应的数据
				tempSeansonEpisodeCount += v2.Statistics.TotalEpisodeCount
				if minEpisode <= tempSeansonEpisodeCount {
					season = v2.SeasonNumber
					// 修改minEpisode 为对应季的数据
					minEpisode = minEpisode - (tempSeansonEpisodeCount - v2.Statistics.TotalEpisodeCount)
					break
				}
			}
		}

	}
	return season, minEpisode
}
