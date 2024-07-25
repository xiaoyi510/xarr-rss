package regexp_ext

import (
	"XArr-Rss/util/logsys"
	"errors"
	"github.com/dlclark/regexp2"
	"strings"
)

type RegexpExtGroups struct {
	Groups map[string][]string
}

func ParseGroups(reg *regexp2.Regexp, text string) *RegexpExtGroups {
	//var qualitySourceReg = regexp2.MustCompile(`\b(?:(?<bluray>BluRay|Blu-Ray|HD-?DVD|BDMux|BD(?!$))|(?<webdl>WEB[-_. ]DL|WEBDL|AmazonHD|iTunesHD|MaxdomeHD|NetflixU?HD|WebHD|[. ]WEB[. ](?:[xh]26[45]|DDP?5[. ]1)|[. ](?-i:WEB)$|\d+0p(?:[-. ]AMZN)?[-. ]WEB[-. ]|WEB-DLMux|\b\s\/\sWEB\s\/\s\b|AMZN[. ]WEB[. ])|(?<webrip>WebRip|Web-Rip|WEBMux)|(?<hdtv>HDTV)|(?<bdrip>BDRip)|(?<brrip>BRRip)|(?<dvd>DVD|DVDRip|NTSC|PAL|xvidvd)|(?<dsr>WS[-_. ]DSR|DSR)|(?<pdtv>PDTV)|(?<sdtv>SDTV)|(?<tvrip>TVRip))(?:\b|$|[ .])`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)
	names := reg.GetGroupNames()
	retGroups := &RegexpExtGroups{}
	retGroups.Groups = make(map[string][]string, len(names))
	match, err := reg.FindStringMatch(text)
	if err != nil || match == nil {
		if err != nil {
			logsys.Debug("正则匹配分组错误:%s", "ParseGroups", err.Error())
		}
		return nil
	}
	for {
		if match == nil || err != nil {
			break
		}
		for _, v := range names {
			// 0 分组是全部匹配内容 不处理
			//if v == "0" {
			//	continue
			//}
			// 根据分组名获取信息
			groupInfo := match.GroupByName(v)
			if groupInfo != nil && groupInfo.String() != "" {
				retGroups.Groups[v] = append(retGroups.Groups[v], groupInfo.String())
			}
		}
		// 寻找下一个
		match, err = reg.FindNextMatch(match)
	}

	return retGroups
}

func (this *RegexpExtGroups) GetGroupValByName(name string) []string {
	v, ok := this.Groups[name]
	if ok {
		return v
	}
	return []string{}
}

func (this *RegexpExtGroups) GetGroupVals() map[string][]string {
	return this.Groups
}

func MatchString(str, regex string) error {
	reg, err := regexp2.Compile(regex, regexp2.IgnoreCase)
	if err != nil {
		return err
	}
	matchString, err := reg.MatchString(str)
	if err != nil {
		return err
	}
	if !matchString {
		return errors.New("未能匹配成功")
	}
	return nil
}

func MatchStringArr(str []string, regex string) error {
	reg, err := regexp2.Compile(regex, regexp2.IgnoreCase)
	if err != nil {
		return err
	}
	matchString, err := reg.MatchString(strings.Join(str, ","))
	if err != nil {
		return err
	}
	if !matchString {
		return errors.New("未能匹配成功")
	}
	return nil
}
