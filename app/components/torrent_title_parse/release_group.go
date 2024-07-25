package torrent_title_parse

import (
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
)

var WebsitePrefixRegex = regexp2.MustCompile(`^\[\s*[-a-z]+(\.[a-z]{2,})+\s*\][- ]*|^www\.[a-z]+\.(?:com|net|org)[ -]*`, regexp2.Compiled|regexp2.IgnoreCase)
var AnimeReleaseGroupRegex = regexp2.MustCompile(`^(?:\[(?<subgroup>(?!\s).+?(?<!\s))\](?:_|-|\s|\.)?)`, regexp2.Compiled|regexp2.IgnoreCase)
var CleanReleaseGroupRegex = regexp2.MustCompile(`^(.*?[-._ ](S\d+E\d+)[-._ ])|(-(RP|1|NZBGeek|Obfuscated|Scrambled|sample|Pre|postbot|xpost|Rakuv[a-z0-9]*|WhiteRev|BUYMORE|AsRequested|AlternativeToRequested|GEROV|Z0iDS3N|Chamele0n|4P|4Planet|AlteZachen))+$`, regexp2.Compiled|regexp2.IgnoreCase)
var ReleaseGroupRegex = regexp2.MustCompile(`-(?<releasegroup>[a-z0-9]+(?!.+?(?:480p|720p|1080p|2160p)))(?<!.*?WEB-DL|Blu-Ray|480p|720p|1080p|2160p|DTS-HD|DTS-X|DTS-MA|DTS-ES)(?:\b|[-._ ]|$)|[-._ ]\[(?<releasegroup>[a-z0-9]+)\]$`, regexp2.Compiled|regexp2.IgnoreCase)

func (this *TorrentTitleParse) ParseReleaseGroup(result *MatchResult) {
	//title = title.Trim();
	//title = RemoveFileExtension(title);
	//for _, replace := range PreSubstitutionRegex {
	//	s, err := replace.Regex.Replace(result.AnalyzeTitle, replace.Replace, -1, -1)
	//	if err != nil {
	//		continue
	//	}
	//	result.AnalyzeTitle = s
	//	break
	//}

	//animeMatch, err := AnimeReleaseGroupRegex.FindStringMatch(result.AnalyzeTitle)
	//if err != nil || animeMatch == nil {
	//	return
	//}
	result.AnalyzeTitle, _ = WebsitePrefixRegex.Replace(result.AnalyzeTitle, "", -1, -1)

	regGroups := regexp_ext.ParseGroups(AnimeReleaseGroupRegex, result.AnalyzeTitle)
	if regGroups == nil {
		return
	}

	if subgroup := regGroups.GetGroupValByName("subgroup"); len(subgroup) > 0 {
		result.AnalyzeTitle, _ = AnimeReleaseGroupRegex.Replace(result.AnalyzeTitle, "", -1, -1)
		result.ReleaseGroup = subgroup[0]
		return
	}

	result.AnalyzeTitle, _ = CleanReleaseGroupRegex.Replace(result.AnalyzeTitle, "", -1, -1)

	//releaseGroups := regexp_ext.ParseGroups(ReleaseGroupRegex, result.AnalyzeTitle)
	//if regGroups == nil {
	//	return
	//}
	//if releaseGroup := releaseGroups.GetGroupValByName("releasegroup"); len(releaseGroup) > 0 {
	//
	//	return releaseGroup[len(releaseGroup)-1]
	//}

	//var matches, _ = ReleaseGroupRegex.FindStringMatch(result.AnalyzeTitle)
	//
	//if matches.GroupCount() > 0 {
	//	var group = matches.OfType < Match > ().Last().Groups["releasegroup"].Value
	//
	//	var groupIsNumeric int
	//
	//	if (int.TryParse(group, out groupIsNumeric))
	//	{
	//	return null;
	//	}
	//
	//	if InvalidReleaseGroupRegex.IsMatch(group) {
	//		return null
	//	}
	//
	//	return group
	//}
	//
	//return null
}
