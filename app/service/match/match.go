package match

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/matchv2"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/app/service/sonarr"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/util/hash"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/xiaoyi510/regex_group"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type MediaMatchService struct {
}

// 过滤发布组
func FilterPushGroups(sourceItems []dbmodel.SourceItem, pushGroup []string) []dbmodel.SourceItem {
	if pushGroup == nil || len(pushGroup) == 0 {
		return sourceItems
	}

	// 两两组合pushGroup
	list := []string{}
	for _, v := range pushGroup {
		for _, v2 := range pushGroup {
			if v2 != v {
				list = append(list, "["+v+"&"+v2+"]")
			}
		}
		list = append(list, "["+v+"]")
	}

	ret := []dbmodel.SourceItem{}
	for _, item := range sourceItems {
		searchText := strings.ToLower(item.Title)
		searchText = strings.Replace(searchText, "【", "[", -1)
		searchText = strings.Replace(searchText, "】", "]", -1)
		search := false
		for _, searchPushGroup := range list {
			if strings.Contains(searchText, strings.ToLower(searchPushGroup)) {
				search = true
				break
			}
		}
		if search || len(list) == 0 {
			ret = append(ret, item)
		}

	}
	return ret
}

// 匹配媒体集数 季数信息
func ParseGroupMediaInfo(groupMedia *dbmodel.GroupMedia, useCache bool) *model.RssRoot {

	// 获取媒体使用的source
	var sourceItems []dbmodel.SourceItem
	sourceItems = sources.GetSourceItems(groupMedia.UseSource, useCache)
	// 判断媒体数据源中是否有数据
	if len(sourceItems) == 0 {
		return nil
	}

	// 判断发布组是否合规

	return ParseGroupMediaSourceItems(groupMedia, sourceItems)

}

// 将数据源项进行格式化
func ParseGroupMediaSourceItems(groupMedia *dbmodel.GroupMedia, items []dbmodel.SourceItem) *model.RssRoot {
	// 过滤字幕组
	items = FilterPushGroups(items, groupMedia.FilterPushGroup)
	// v1.5.6
	// 判断过滤字幕组后是否还有数据
	if len(items) == 0 {
		return nil
	}
	groupsMediaCache := &model.RssRoot{Version: "2.0"}
	groupsMediaCache.Channel.Link = appconf.AppConf.System.HttpAddr + "/rss/group_" + strconv.Itoa(int(groupMedia.GroupId)) + "/" + strconv.Itoa(groupMedia.Id) + ".xml"
	groupsMediaCache.Channel.Description = "XArr-Rss 分组订阅 " + groupMedia.MediaInfo.CnTitle
	groupsMediaCache.Channel.Title = "XArr-Rss-" + groupMedia.MediaInfo.CnTitle
	groupMedia.MediaInfo.SearchTitle = strings.Join(matchv2.MatchV2{}.JoinSonarrTitle(groupMedia.MediaInfo), "1666666666663")

	matchv2.MatchV2{}.Parse(groupMedia, items, &groupsMediaCache.Channel)

	// 保存缓存
	return groupsMediaCache
}

// 正则匹配符合数据的资源
func (this MediaMatchService) ParseSourceItemsFromRegex(groupMedia dbmodel.GroupMedia, groupRegex dbmodel.GroupRegex, sourceItems []model.RssResultItem) (error, []model.RssResultItem) {
	err, reg := this.checkAndGenReg(groupRegex)
	if err != nil {
		return err, nil
	}

	// 如果不为自动并且含有错误就返回
	if groupRegex.MatchType != "auto" && err != nil {
		return logsys.Error("媒体正则表达式异常了吧:%s", "正则", err.Error()), nil
	}
	err, sonarrMediaInfo := sonarr.SonarrService{}.GetSonarrMediaById(groupMedia.SonarrId)
	if err != nil {
		return err, nil
	}

	groupCache := model.RssRoot{}

	for _, sourceItem := range sourceItems {
		// 正对souceItem格式化
		this.parseSourceItemLength(&sourceItem)

		// 判断匹配模式
		mediaParseInfo := &torrent_title_parse.MatchResult{}
		if groupRegex.MatchType == "auto" {
			// 自动
			err, mediaParseInfo = this.parseSourceItemsRegexAuto(groupMedia.SonarrId, groupRegex, sourceItem, reg, sonarrMediaInfo)
			//mediaParseInfo.AbsoluteMinEpisode = mediaParseInfov2.MinEpisode
			//mediaParseInfo.MinEpisode = mediaParseInfov2.MinEpisode
			//mediaParseInfo.MaxEpisode = mediaParseInfov2.MaxEpisode
			//mediaParseInfo.Language = mediaParseInfov2.Language
			//mediaParseInfo.Quality = mediaParseInfov2.QualityResolution
			//mediaParseInfo.Season = mediaParseInfov2.Season
		} else {
			// 手动
			err, mediaParseInfo = this.parseSourceItemsRegexManual(groupRegex, sourceItem, reg, sonarrMediaInfo)
			// 处理匹配语言和质量
			mediaParseInfo.Language = this.parseMediaLanguage(groupMedia.Language, sourceItem.Title.Text)
			mediaParseInfo.QualityResolution = this.parseMediaQuality(groupMedia.Quality, sourceItem.Title.Text)
		}

		// 没有匹配到信息就不执行了
		if err != nil {
			if strings.Contains(err.Error(), "非赞助会员") {
				return err, nil
			}
			continue
		}

		mediaParseInfo.Title = sourceItem.Title.Text

		//////////////////////////////处理标题

		var title = ""
		// 替换为Sonarr能识别的标题
		oldMaxEpisode := mediaParseInfo.AbsoluteMinEpisode
		if mediaParseInfo.MaxEpisode > mediaParseInfo.MinEpisode && mediaParseInfo.MaxEpisode > 0 {
			oldMaxEpisode = mediaParseInfo.AbsoluteMinEpisode + (mediaParseInfo.MaxEpisode - mediaParseInfo.MinEpisode)

			// 范围型
			title = fmt.Sprintf("%s - S%02dE%02d-%02d Episode %02d-%02d - %s - %s",
				//groupMedia.MediaInfo.CnTitle,
				groupMedia.MediaInfo.OriginalTitle,
				mediaParseInfo.Season,
				mediaParseInfo.MinEpisode+groupRegex.Offset,
				mediaParseInfo.MaxEpisode+groupRegex.Offset,
				mediaParseInfo.AbsoluteMinEpisode,
				oldMaxEpisode,
				mediaParseInfo.Language,
				mediaParseInfo.QualityResolution,
			)
		} else {
			title = fmt.Sprintf(
				"%s - S%02dE%02d Episode %d - %s - %s",
				//groupMedia.MediaInfo.CnTitle,
				groupMedia.MediaInfo.OriginalTitle,
				mediaParseInfo.Season,
				mediaParseInfo.MinEpisode+groupRegex.Offset,
				mediaParseInfo.AbsoluteMinEpisode,
				mediaParseInfo.Language,
				mediaParseInfo.QualityResolution,
			)

		}
		if mediaParseInfo.MaxEpisode == 0 {
			mediaParseInfo.MaxEpisode = mediaParseInfo.MinEpisode
		}

		groupCache.Channel.Item = append(groupCache.Channel.Item, model.RssResultItem{
			Title:         model.CDATA{Text: title},
			CnTitle:       sonarrMediaInfo.CnTitle,
			OriginalTitle: model.CDATA{Text: sourceItem.Title.Text},
			OtherTitle:    model.CDATA{Text: " other:" + strings.Join(sonarrMediaInfo.AlternateTitles, "-|-")},
			PubDate:       sourceItem.PubDate,
			OldMinEpisode: mediaParseInfo.AbsoluteMinEpisode,
			OldMaxEpisode: oldMaxEpisode,
			Season:        mediaParseInfo.Season,
			MinEpisode:    mediaParseInfo.MinEpisode + groupRegex.Offset,
			MaxEpisode:    mediaParseInfo.MaxEpisode + groupRegex.Offset,
			Enclosure: model.RssResultItemEnclosure{
				Type:   sourceItem.Enclosure.Type,
				Length: sourceItem.Enclosure.Length,
				Url:    appconf.AppConf.System.HttpAddr + "/api/v1/down?url=" + url.QueryEscape(sourceItem.Enclosure.Url) + "&title=" + url.QueryEscape(title) + "&source_id=" + url.QueryEscape(sourceItem.SourceId),
			},
			Link: sourceItem.Link,
			Guid: model.RssResultItemGuid{
				IsPermaLink: sourceItem.Guid.IsPermaLink,
				Text:        hash.Md5{}.HashString(sourceItem.Guid.Text + sourceItem.Enclosure.Url + sourceItem.Enclosure.Length + sourceItem.SourceId),
			},
			XArrRssIndexer: model.RssResultItemXArrRssIndexer{
				Text: " 数据源:" + sourceItem.SourceId,
				//Id:   groupMedia.Id,
			},
			OthderId: model.RssResultItemOtherId{
				TmdbId: groupMedia.MediaInfo.TmdbId,
				TvdbId: strconv.Itoa(groupMedia.MediaInfo.TvdbId),
				ImdbId: groupMedia.MediaInfo.ImdbId,
			},
		})

	}
	return nil, groupCache.Channel.Item
}

// 处理种子大小
func (this MediaMatchService) parseSourceItemLength(sourceItem *model.RssResultItem) {
	if sourceItem.Enclosure.Length == "0" || sourceItem.Enclosure.Length == "" || helper.StrToInt(sourceItem.Enclosure.Length) < 1024 {
		sourceItem.Enclosure.Length = "996973766"
	}
}

// 手工正则匹配
func (this MediaMatchService) parseSourceItemsRegexManual(regStr dbmodel.GroupRegex, sourceItem model.RssResultItem, reg *regexp.Regexp, sonarrMediaInfo dbmodel.Media) (error, *torrent_title_parse.MatchResult) {
	ret := &torrent_title_parse.MatchResult{}

	// 判断能否匹配到标题信息
	if reg.MatchString(sourceItem.Title.Text) {
		// 提取集数信息
		matchGroups := regex_group.GetRegxGroupByReg(reg, sourceItem.Title.Text)
		if matchGroups == nil || len(matchGroups) == 0 {
			// 无匹配内容
			return errors.New("没有匹配到内容,请检查是否有设置命名分组"), ret
		}
		// 获取集标记
		episodeGroup, ok := matchGroups[0]["episode"]
		if !ok {
			return logsys.Error("正则匹配数量错误,请检查是否有(?<episode>\\d+)：%v", "数据源同步至分组", regStr.Reg), ret
		}
		minEpisode := helper.StrToInt(episodeGroup.SubMatch)

		// 获取最大集标记
		maxEpisodeGroup, ok := matchGroups[0]["max_episode"]
		if ok {
			ret.MaxEpisode = helper.StrToInt(maxEpisodeGroup.SubMatch)
		}

		// 获取季标记
		seasonGroup, ok := matchGroups[0]["season"]
		if !ok {
			ret.Season = helper.StrToInt(seasonGroup.SubMatch)
		}

		// 提取对应第几集
		//allString := reg.FindAllStringSubmatch(sourceItem.Title.Text, -1)
		//if len(allString) == 0 || len(allString[0]) < 2 {
		//	return logsys.Error("正则匹配数量错误,请检查是否有(?<episode>\\d+)：%v", "数据源同步至分组", regStr.Reg), ret
		//}
		// 获取集数
		//minEpisode, _ := strconv.Atoi(allString[0][1])
		//if len(allString[0]) >= 3 {
		//	if !variable.ServerState.IsVip {
		//		return logsys.Error("非赞助会员用户只支持单集规则", "数据源同步至分组", regStr.Reg), ret
		//	}
		//	maxEpisode, _ := strconv.Atoi(allString[0][2])
		//	ret.MaxEpisode = maxEpisode
		//}

		// 设置设置的季数信息
		ret.Season, ret.MinEpisode = this.ParseMediaSeason(minEpisode, regStr.Season, sonarrMediaInfo)
		// minEpisode 可能会因为 offset变化
		ret.AbsoluteMinEpisode = minEpisode

		return nil, ret
	}
	return errors.New("没有匹配到"), ret
}

// 自动匹配
func (this MediaMatchService) parseSourceItemsRegexAuto(sonarrId int, groupRegex dbmodel.GroupRegex, sourceItem model.RssResultItem, reg *regexp.Regexp, sonarrMediaInfo dbmodel.Media) (error, *torrent_title_parse.MatchResult) {
	ret := &torrent_title_parse.MatchResult{}
	// 提取所有sonarr标题
	// 取出sonarr title
	sonarrTitle := sonarrMediaInfo.OriginalTitle
	if sonarrMediaInfo.CnTitle != "" {
		sonarrTitle += "-|-" + sonarrMediaInfo.CnTitle
	}
	if len(sonarrMediaInfo.Titles) > 0 {
		sonarrTitle += "-|-" + strings.Join(sonarrMediaInfo.Titles, "-|-")
	}
	if len(sonarrMediaInfo.AlternateTitles) > 0 {
		sonarrTitle += "-|-" + strings.Join(sonarrMediaInfo.AlternateTitles, "-|-")
	}

	if sonarrTitle == "" {
		return logsys.Error("自动匹配模式未能找到对应Sonarr数据可能被删除了,SonarrId:"+strconv.Itoa(sonarrId), "分组数据自动匹配"), ret
	}
	var err error
	var parseData bool // *torrent_title_parse.MatchResult
	regTitleText := sourceItem.Title.Text

	// 解析标题信息
	parse := torrent_title_parse.TorrentTitleParse{}
	ret = parse.Parse(regTitleText)

	for true {
		// 处理匹配规则
		if ret.AnalyzeTitle != "" {
			err, parseData = ParseByAutoReg(ret.AnalyzeTitle, sonarrTitle, groupRegex.Reg)
		} else {
			logsys.Debug("不支持的文件标题:%s", "分组数据自动匹配", regTitleText)
			//err, parseData = ParseByAutoReg(regTitleText, sonarrTitle, groupRegex.Reg)
		}
		if err != nil {
			if !strings.Contains(err.Error(), "未能匹配成功") {
				logsys.Error("自动匹配模式错误:%s", "分组数据自动匹配", err.Error())
			}
			return err, ret
		} else {
			break

			//if parseData == nil {
			//	break
			//}

			// 把标题给他删了
			//if (parseData.MinEpisode) > sonarrMediaInfo.TotalEpisodeCount {
			//	//logsys.Debug("剧集匹配成功但是集数不在Sonarr总集数范畴中 集数:%d Sonarr总集数:%d", "分组数据自动匹配", parseData.MinEpisode, sonarrMediaInfo.TotalEpisodeCount)
			//	regTitleText = strings.Replace(regTitleText, strconv.Itoa(parseData.MinEpisode), " ", -1)
			//} else {
			//	break
			//}

		}
	}

	if parseData != false {
		//ret.Title = parseData.Title
		//ret.MinEpisode = parseData.MinEpisode
		//ret.MaxEpisode = parseData.MaxEpisode
		//
		//// 自动匹配季度和集数
		ret.AbsoluteMinEpisode = ret.MinEpisode
		ret.Season, ret.MinEpisode = this.ParseMediaSeason(ret.MinEpisode, groupRegex.Season, sonarrMediaInfo)

		return nil, ret
	}

	return errors.New("没有匹配到"), ret
}

// 检测和生成正则表达式
func (this MediaMatchService) checkAndGenReg(groupRegex dbmodel.GroupRegex) (error, *regexp.Regexp) {
	if groupRegex.Reg == "" && groupRegex.MatchType != "auto" {
		return errors.New("正则表达式匹配规则为空"), nil
	}
	// 初始化正则方法
	reg, err := regexp.Compile(groupRegex.Reg)
	// 判断正则是否正确
	if err != nil {
		return logsys.Error("正则数据错误:%s err：%v", "数据源同步至分组", groupRegex.Reg, err.Error()), nil
	}
	if reg == nil {
		return logsys.Error("正则数据错误：%v", "数据源同步至分组", groupRegex.Reg), nil
	}
	return nil, reg
}

func (this MediaMatchService) parseMediaLanguage(language string, text string) string {
	if language == "-1" {
		if variable.ServerState.IsVip {
			language = this.ParseMediaLanguage(text)
		} else {
			logsys.Error("非赞助会员用户不可以使用自动匹配语言", "会员")
		}
	}
	return language
}

func (this MediaMatchService) parseMediaQuality(quality string, text string) string {
	if quality == "-1" {
		if variable.ServerState.IsVip {
			quality = this.ParseMediaQuality(text)
		} else {
			logsys.Error("非赞助会员用户不可以使用自动匹配质量", "会员")
		}
	}
	return quality
}

// 计算媒体第几季
func (this MediaMatchService) ParseMediaSeason(minEpisode, season int, sonarrMediaInfo dbmodel.Media) (int, int) {
	if season == -1 {
		if !variable.ServerState.IsVip {
			logsys.Error("只非赞助会员用户不可以使用自动匹配季", "会员")
			return season, minEpisode
		}
		//err, searchSeason := ParseMediaSeason(title)
		//if err == nil {
		//	if searchSeason != "" {
		//		season = helper.StrToInt(searchSeason)
		//	}
		//}
		if season > 0 {
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

// 匹配标题中的第几集
func ParseMediaTitleEpisode(title string) (err error, retRegex, retSubGrop, retMinEpisode, retMaxEpisode string) {
	// 剔除 1080
	reg := regexp2.MustCompile(`\b(1080p|1920x1080|720p|1280x720|480p|2160p)\b`, regexp2.IgnoreCase|regexp2.Compiled)

	// 替换为特殊符号
	tempTitle, err := reg.ReplaceFunc(title, func(m regexp2.Match) string {
		return strings.Repeat(string([]byte{127}), m.Length)
	}, -1, -1)

	//log.Println(tempTitle, err)

	if strings.Contains(tempTitle, `(?<episode>\d+)`) {
		return errors.New("已有集标记,无需重复识别"), "", "", "", ""
	}
	if !variable.ServerState.IsVip {
		return errors.New("非赞助会员用户不可以使用当前功能"), "", "", "", ""
	}

	//myReg, err := regexp.Compile(`(?m)^.*?.*?(?:[\[\sEe【第_\(]|Ep)(?P<episode>\d+)[\-\~]{0,1}e?(?P<max_episode>\d+)?(?P<version>v\d+)?([\)】集话\]\s]|Fin).*?$`)
	myReg, err := regexp.Compile(`(?m)^(?:\[连载\])?[【\[](?P<subgroup>[^\]】]+).*?(?:[\[\sEe【第_\(]|Ep)(?P<episode>\d+)[\-\~]{0,1}e?(?P<max_episode>\d+)?(?P<version>v\d+)?([\)】集话\]\s]|Fin).*?$`)
	if err != nil {
		return errors.New("匹配错误异常:" + err.Error()), "", "", "", ""
	}

	if !myReg.MatchString(strings.ToLower(tempTitle)) {
		return errors.New("未能匹配成功:" + title), "", "", "", ""
	}
	matchData := regex_group.GetRegxGroupByRegOne(myReg, strings.ToLower(tempTitle))
	//log.Println(matchData.GetString("episode"))
	//log.Println(matchData.GetString("max_episode"))
	//log.Println(matchData.GetString("version"))
	// 查询子查询
	indexs := myReg.FindStringSubmatchIndex(strings.ToLower(tempTitle))
	// 3*2 3 代表几个子查询
	if len(indexs) < 5*2 {
		return errors.New("未能匹配成功:" + title), "", "", "", ""
	}
	//allMatches := myReg.FindStringSubmatch(title)

	edipose := ""
	maxEpisode := ""
	subgroup := ""

	if matchData.Get("max_episode") != nil && matchData.Get("max_episode").SubMatch != "" {
		maxEpisode = matchData.Get("max_episode").SubMatch
		// 处理替换字符串
		if matchData.Get("max_episode").Offset < 0 {
			marshal, _ := json.Marshal(matchData.Get("max_episode"))
			logsys.Error("解析错误了:%s,%s", "解析", string(marshal), title)
		} else {
			startS := title[0:matchData.Get("max_episode").Offset] // 开始的位置
			startE := title[matchData.Get("max_episode").Size:]    // 开始的结束

			title = startS + `-----max_episode-----` + startE
		}
	}

	if matchData.Get("episode") != nil && matchData.Get("episode").SubMatch != "" {
		edipose = matchData.Get("episode").SubMatch
		// 处理替换字符串
		if matchData.Get("episode").Offset < 0 {
			marshal, _ := json.Marshal(matchData.Get("episode"))
			logsys.Error("解析错误了:%s,%s", "解析", string(marshal), title)
		} else {
			//log.Println(title[0:matchData.Get("episode").Offset], tempTitle[0:matchData.Get("episode").Offset])
			startS := title[0:matchData.Get("episode").Offset] // 开始的位置
			startE := title[matchData.Get("episode").Size:]    // 开始的结束
			title = startS + `-----episode-----` + startE
		}
	}

	if matchData.Get("subgroup") != nil && matchData.Get("subgroup").SubMatch != "" {
		subgroup = matchData.Get("subgroup").SubMatch
		// 处理替换字符串
		if matchData.Get("subgroup").Offset < 0 {
			marshal, _ := json.Marshal(matchData.Get("subgroup"))
			logsys.Error("解析错误了:%s,%s", "解析", string(marshal), title)
		} else {
			startS := title[0:matchData.Get("subgroup").Offset] // 开始的位置
			startE := title[matchData.Get("subgroup").Size:]    // 开始的结束

			title = startS + `-----subgroup-----` + startE
		}
	}
	a, _ := regexp.Compile(`([\\^$.*+?()[\]{}|])`)
	title = a.ReplaceAllString(title, "\\$1")
	title = helper.StrReplace(title, []string{
		"-----subgroup-----",
		"-----max_episode-----",
		"-----episode-----",
	}, []string{
		`(?<subgroup>[^\]】]+)`,
		`(?<max_episode>\d+)`,
		`(?<episode>\d+)`,
	})
	return nil, title, subgroup, edipose, maxEpisode
}

// 匹配季度
func ParseMediaSeason(title string) (error, string) {
	//myReg, err := regexp.Compile(`^.*?(?:season|Season|S|s|第)[\s-]*?(?<season>[\d一二三四五六七八九十]+)[季]?.*?$`)

	regStrArr := []string{
		`^.*?(?:season|Season|S|s)[\s-]?(?<season>(\d+)).*?$`,
		`^.*?(?:第)[\s-]?(?<season>[\d一二三四五六七八九十]+[季期]).*?$`,
	}
	for _, regStr := range regStrArr {
		one := regex_group.GetRegxGroupOne(regStr, title)
		if one == nil {
			continue
		}
		season := one.GetString("season")
		if season == "" {
			continue
		}

		//myReg, err := regexp.Compile(v)
		//if err != nil {
		//	return errors.New("匹配错误异常:" + err.Error()), ""
		//}
		//if !myReg.MatchString(title) {
		//	continue
		//}
		//allMatches := myReg.FindStringSubmatch(title)
		// 去掉季
		//allMatches[1] = strings.Replace(allMatches[1], "季", "", -1)

		season = helper.StrReplace(season, []string{"季", "期"}, []string{""})
		arr := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
		for i, v := range arr {
			//allMatches[1] = strings.Replace(allMatches[1], v, strconv.Itoa(i+1), -1)
			season = strings.Replace(season, v, strconv.Itoa(i+1), -1)
		}
		return nil, season
	}
	return errors.New("未能匹配成功:" + title), ""

}

// 自动匹配正则
func ParseByAutoReg(title, sonarrTitle, matchReq string) (error, bool) {
	// 判断是否筛选条件是否为空
	pp := 0
	maxPP := 1
	//oldTitle := title
	if matchReq != "" {
		// 过滤包含
		// 有匹配规则
		matchAndArr := strings.Split(matchReq, ",")
		for _, matchAnd := range matchAndArr {
			// 判断是否有|
			if strings.Contains(matchAnd, "|") {
				matchOrArr := strings.Split(matchAnd, "|")
				for _, matchOr := range matchOrArr {
					// 如果有一个实现效果 则跳出处理
					if strings.Contains(title, matchOr) {
						pp++
						break
					}
				}
				if pp == 0 {
					return errors.New("未能匹配成功"), false
				}
				pp = 0

			} else {
				// 单独一个
				if strings.Contains(title, matchAnd) {
					//pp++
				} else {
					// 不包含指定字符串
					// 返回
					return errors.New("未能匹配成功"), false
				}
			}
		}
	}

	///////////////////////////////去除标题
	// 分割sonarrTitle
	sonarrTitleArr := strings.Split(sonarrTitle, "-|-")

	// 处理原有标题
	//title2 := helper.ReplaceRegString(title)

	for _, v := range sonarrTitleArr {
		v = helper.ReplaceRegString(v)
		if v == "" || len(v) <= 4 {
			continue
		}

		// 模糊匹配
		hasMatch, matchString := helper.MatchString(strings.ToLower(v), strings.ToLower(title))
		if hasMatch {
			title = strings.Replace(title, matchString, "", -1)
			pp++
			continue
		}
		// 判断标题是否存在
		//if strings.Contains(title2, v) {
		//	// 匹配成功
		//	// 万能匹配
		//	pp++
		//	break
		//}
	}
	///////////////////////////////去除标题

	if pp >= maxPP {
		//// 都匹配成功了 再来匹配万能
		//// 万能匹配
		//err, _, min, max := ParseMediaTitleEpisode(title)
		//if err == nil {
		return nil, true
		//}
		//return err, nil
	}
	return nil, false
}
