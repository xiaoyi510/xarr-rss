package torrent_title_parse

import (
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
	"strings"
)

const (
	Type_Unknown       = "Unknown"
	Type_Television    = "Television"
	Type_TelevisionRaw = "TelevisionRaw"
	Type_WebDl         = "WebDL"
	Type_WebRip        = "WebRip"
	Type_DVD           = "DVD"
	Type_HDTV          = "HDTV"
	Type_SDTV          = "SDTV"
	Type_Bluray        = "Bluray"
	Type_BlurayRaw     = "BlurayRaw"
)

const (
	Resolution_Unknown = "Unknown"
	Resolution_360p    = "360p"
	Resolution_480p    = "480p"
	Resolution_576p    = "576p"
	Resolution_720p    = "720p"
	Resolution_1080p   = "1080p"
	Resolution_2160p   = "2160p"
)

const (
	Quality_Unknown = "Unknown"
	Quality_SDTV    = "SDTV"
	Quality_RAWHD   = "Raw-HD"
	Quality_DVD     = "DVD"

	Quality_HDTV480p   = "HDTV-480p"
	Quality_WEBDL480p  = "WEBDL-480p"
	Quality_WEBRip480p = "WEBRip-480p"
	Quality_Bluray480p = "Bluray-480p"

	Quality_HDTV720p   = "HDTV-720p"
	Quality_WEBRip720p = "WEBRip-720p"
	Quality_WEBDL720p  = "WEBDL-720p"
	Quality_Bluray720p = "Bluray-720p"

	Quality_HDTV1080p        = "HDTV-1080p"
	Quality_WEBDL1080p       = "WEBDL-1080p"
	Quality_WEBRip1080p      = "WEBRip-1080p"
	Quality_Bluray1080p      = "Bluray-1080p"
	Quality_Bluray1080pRemux = "Bluray-1080p Remux"

	Quality_HDTV2160p        = "HDTV-2160p"
	Quality_WEBDL2160p       = "WEBDL-2160p"
	Quality_WEBRip2160p      = "WEBRip-2160p"
	Quality_Bluray2160p      = "Bluray-2160p"
	Quality_Bluray2160pRemux = "Bluray-2160p Remux"
)

// 正则匹配规则
var qualitySourceReg = regexp2.MustCompile(`\b(?:(?<bluray>BluRay|Blu-Ray|HD-?DVD|BDMux|BD(?!$))|(?<webdl>WEB[-_. ]DL|WEBDL|AmazonHD|iTunesHD|MaxdomeHD|NetflixU?HD|WebHD|[. ]WEB[. ](?:[xh]26[45]|DDP?5[. ]1)|[. ](?-i:WEB)$|\d+0p(?:[-. ]AMZN)?[-. ]WEB[-. ]|WEB-DLMux|\b\s\/\sWEB\s\/\s\b|AMZN[. ]WEB[. ])|(?<webrip>WebRip|Web-Rip|WEBMux)|(?<hdtv>HDTV)|(?<bdrip>BDRip)|(?<brrip>BRRip)|(?<dvd>DVD|DVDRip|NTSC|PAL|xvidvd)|(?<dsr>WS[-_. ]DSR|DSR)|(?<pdtv>PDTV)|(?<sdtv>SDTV)|(?<tvrip>TVRip))(?:\b|$|[ .])`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)
var webdlReg = regexp2.MustCompile(`\[WEB\]`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)

var animeReg = regexp2.MustCompile(`bd(?:720|1080|2160)|(?<=[-_. (\[])bd(?=[-_. )\]])`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)

var HighDefPdtvRegex = regexp2.MustCompile(`hr[-_. ]ws`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)
var OtherSourceRegex = regexp2.MustCompile(`(?<hdtv>HD[-_. ]TV)|(?<sdtv>SD[-_. ]TV)`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)
var CodecRegex = regexp2.MustCompile(`\b(?:(?<x264>x264)|(?<h264>h264)|(?<xvidhd>XvidHD)|(?<xvid>Xvid)|(?<divx>divx))\b`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)

var resolutionRegex = regexp2.MustCompile(`(\b|_)(?:(?<R360p>360p)|
(?<R480p>480p|640x480|848x480)|
(?<R576p>576p)|
(?<R720p>720p|1280x720)|
(?<R1080p>1080p\+?|1920x1080|1440p|FHD|1080i|4kto1080p|web1080p\+?)|
(?<R2160p>2160p|3840x2160|4k|4k[-_. ](?:UHD|HEVC|BD)|(?:UHD|HEVC|BD)[-_. ]4k))
(\b|$|_)`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)

// 解析质量
func (this *TorrentTitleParse) ParseQuality(result *MatchResult) {
	// 匹配1080p

	// 获取质量类型
	// 获取匹配的质量尺寸
	resolution := this.ParseResolution(result)
	// 判断是否为原盘
	remuxMatch := this.ParseRemux(result)

	if this.matchQualitySource(resolution, remuxMatch, result) {
		return
	}

	// 匹配动画 蓝光
	if this.matchAnime(resolution, remuxMatch, result) {
		return
	}

	// 匹配webdl
	if this.matchWebdl(resolution, remuxMatch, result) {
		return
	}

	if this.matchContains(result) {
		return
	}

	if this.OtherSourceMatch(result) {
		return
	}

	result.QualitySource = Type_Unknown
	result.Resolution = Resolution_Unknown
}

func (this *TorrentTitleParse) matchQualitySource(resolution string, remuxMatch bool, result *MatchResult) bool {
	codecRegexGroups := regexp_ext.ParseGroups(CodecRegex, result.AnalyzeTitle)

	qualitySourceRegGroups := regexp_ext.ParseGroups(qualitySourceReg, result.AnalyzeTitle)
	if qualitySourceRegGroups == nil {
		return false
	}

	result.AnalyzeTitle, _ = qualitySourceReg.Replace(result.AnalyzeTitle, "", -1, -1)

	if len(qualitySourceRegGroups.GetGroupValByName("bluray")) > 0 {
		if codecRegexGroups == nil {
			//panic("!23")
		}
		// 蓝光匹配到了
		result.QualitySource = Type_Bluray
		result.Resolution = resolution
		if codecRegexGroups != nil && ((len(codecRegexGroups.GetGroupValByName("xvid")) > 0) ||
			(len(codecRegexGroups.GetGroupValByName("divx")) > 0)) {
			result.QualityResolution = Quality_Bluray480p
			return true

		}

		if resolution == Resolution_2160p {
			if remuxMatch {
				result.QualitySource = Quality_Bluray2160pRemux
			} else {
				result.QualitySource = Quality_Bluray2160p
			}
			return true
		}
		if resolution == Resolution_1080p {
			if remuxMatch {
				result.QualitySource = Quality_Bluray1080pRemux
			} else {
				result.QualitySource = Quality_Bluray1080p
			}
			return true
		}
		if resolution == Resolution_360p || resolution == Resolution_480p || resolution == Resolution_576p {
			result.QualitySource = Quality_Bluray480p
			return true
		}

		// Treat a remux without a source as 1080p, not 720p.
		if remuxMatch {
			result.QualitySource = Quality_Bluray1080pRemux
			return true
		}

		result.QualitySource = Quality_Bluray720p

		return true
	}
	if len(qualitySourceRegGroups.GetGroupValByName("webdl")) > 0 {
		// 匹配到了
		result.QualitySource = Type_WebDl
		result.Resolution = resolution

		if resolution == Resolution_2160p {
			result.QualityResolution = Quality_WEBDL2160p
			return true
		}

		if resolution == Resolution_1080p {
			result.QualityResolution = Quality_WEBDL1080p
			return true
		}

		if resolution == Resolution_720p {
			result.QualityResolution = Quality_WEBDL720p
			return true
		}

		if strings.Contains(result.AnalyzeTitle, "[WEBDL]") {
			result.QualityResolution = Quality_WEBDL720p
			result.Resolution = Resolution_720p
			return true
		}

		result.QualityResolution = Quality_WEBDL480p
		result.Resolution = Resolution_480p
		return true

	}
	if len(qualitySourceRegGroups.GetGroupValByName("webrip")) > 0 {
		// 匹配到了
		result.QualitySource = Type_WebRip
		result.Resolution = resolution
		if resolution == Resolution_2160p {
			result.QualityResolution = Quality_WEBRip2160p
			return true
		}

		if resolution == Resolution_1080p {
			result.QualityResolution = Quality_WEBRip1080p
			return true
		}

		if resolution == Resolution_720p {
			result.QualityResolution = Quality_WEBRip720p
			return true
		}

		result.QualityResolution = Quality_WEBRip480p
		result.Resolution = Resolution_480p
		return true
	}
	if len(qualitySourceRegGroups.GetGroupValByName("hdtv")) > 0 {
		// 匹配到了
		result.QualitySource = Type_HDTV
		result.Resolution = resolution
		//if MPEG2Regex.IsMatch(normalizedName) {
		//	result.QualitySource = Quality_RAWHD
		//	return
		//}

		if resolution == Resolution_2160p {
			result.QualityResolution = Quality_HDTV2160p
			return true
		}

		if resolution == Resolution_1080p {
			result.QualityResolution = Quality_HDTV1080p
			return true
		}

		if resolution == Resolution_720p {
			result.QualityResolution = Quality_HDTV720p
			return true
		}

		if strings.Contains(result.AnalyzeTitle, "[HDTV]") {
			result.Resolution = Resolution_720p
			result.QualityResolution = Quality_HDTV720p
			return true
		}

		result.QualityResolution = Quality_SDTV
		return true
	}
	if (len(qualitySourceRegGroups.GetGroupValByName("bdrip")) > 0) || (len(qualitySourceRegGroups.GetGroupValByName("brrip")) > 0) {
		// 匹配到了
		result.QualitySource = Type_BlurayRaw
		result.Resolution = resolution

		switch resolution {
		case Resolution_720p:
			result.QualityResolution = Quality_Bluray720p
			return true
		case Resolution_1080p:
			result.QualityResolution = Quality_Bluray1080p
			return true
		default:
			result.QualityResolution = Quality_Bluray480p
			return true
		}

	}
	if len(qualitySourceRegGroups.GetGroupValByName("dvd")) > 0 {
		result.QualitySource = Type_DVD
		result.Resolution = Resolution_Unknown
		result.QualityResolution = Quality_DVD
		return true
	}

	if (len(qualitySourceRegGroups.GetGroupValByName("pdtv")) > 0) || (len(qualitySourceRegGroups.GetGroupValByName("sdtv")) > 0) || (len(qualitySourceRegGroups.GetGroupValByName("dsr")) > 0) || (len(qualitySourceRegGroups.GetGroupValByName("tvrip")) > 0) {
		result.QualitySource = Type_HDTV
		result.Resolution = Resolution_Unknown

		if resolution == Resolution_1080p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "1080p") {
			result.Resolution = Resolution_1080p
			result.QualityResolution = Quality_HDTV1080p
			return true
		}

		if resolution == Resolution_720p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "720p") {
			result.Resolution = Resolution_720p

			result.QualityResolution = Quality_HDTV720p
			return true
		}
		matchString, err := HighDefPdtvRegex.MatchString(result.AnalyzeTitle)
		if err == nil && matchString {
			result.QualitySource = Type_HDTV
			result.Resolution = Resolution_720p
			result.QualityResolution = Quality_HDTV720p
			return true
		}

		result.QualityResolution = Quality_SDTV
		return true
	}
	//}
	return false
}

func (this *TorrentTitleParse) matchAnime(resolution string, remuxMatch bool, result *MatchResult) bool {
	animeMatchString, err := animeReg.MatchString(result.AnalyzeTitle)
	if err != nil {
		return false
	}
	if animeMatchString {
		if resolution == Resolution_360p || resolution == Resolution_480p ||
			resolution == Resolution_576p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "480p") {
			//result.ResolutionDetectionSource = QualityDetectionSource.Name
			result.QualityResolution = Quality_DVD
			result.Resolution = Resolution_Unknown

			return true
		}
		result.QualitySource = Type_Bluray
		if resolution == Resolution_1080p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "1080p") {
			//result.ResolutionDetectionSource = QualityDetectionSource.Name
			if remuxMatch {
				result.QualityResolution = Quality_Bluray1080pRemux
			} else {
				result.QualityResolution = Quality_Bluray1080p
			}
			result.Resolution = Resolution_1080p
			return true
		}

		if resolution == Resolution_2160p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "2160p") {
			//result.ResolutionDetectionSource = QualityDetectionSource.Name
			if remuxMatch {
				result.QualityResolution = Quality_Bluray2160pRemux
			} else {
				result.QualityResolution = Quality_Bluray2160p
			}
			result.Resolution = Resolution_2160p
			return true
		}

		// Treat a remux without a source as 1080p, not 720p.
		if remuxMatch {
			result.Resolution = Resolution_1080p

			result.QualityResolution = Quality_Bluray1080p
			return true
		}

		result.Resolution = Resolution_720p
		result.QualityResolution = Quality_Bluray720p
		return true
	}
	return false
}

func (this *TorrentTitleParse) matchWebdl(resolution string, remuxMatch bool, result *MatchResult) bool {
	webdlMatchString, err := webdlReg.MatchString(result.AnalyzeTitle)
	if err != nil {
		return false
	}
	if webdlMatchString {
		result.AnalyzeTitle, _ = webdlReg.Replace(result.AnalyzeTitle, "", -1, -1)

		if resolution == Resolution_360p || resolution == Resolution_480p ||
			resolution == Resolution_576p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "480p") {
			//result.ResolutionDetectionSource = QualityDetectionSource.Name
			result.QualityResolution = Quality_WEBDL480p
			result.QualitySource = Type_WebDl
			result.Resolution = Resolution_480p

			return true
		}
		result.QualitySource = Type_WebDl
		if resolution == Resolution_1080p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "1080p") {
			//result.ResolutionDetectionSource = QualityDetectionSource.Name
			result.QualityResolution = Quality_WEBDL1080p
			result.Resolution = Resolution_1080p
			return true
		}

		if resolution == Resolution_2160p || strings.Contains(strings.ToLower(result.AnalyzeTitle), "2160p") {
			//result.ResolutionDetectionSource = QualityDetectionSource.Name
			result.QualityResolution = Quality_WEBDL2160p

			result.Resolution = Resolution_2160p
			return true
		}

		// Treat a remux without a source as 1080p, not 720p.
		if remuxMatch {
			result.Resolution = Resolution_1080p

			result.QualityResolution = Quality_WEBDL1080p
			return true
		}

		result.Resolution = Resolution_720p
		result.QualityResolution = Quality_WEBDL720p
		return true
	}
	if resolution != Resolution_Unknown {
		result.QualitySource = Type_Unknown
		if remuxMatch {
			result.QualitySource = Type_BlurayRaw
		} else {
			//try
			//{
			//	var quality = MediaFileExtensions.GetQualityForExtension(name.GetPathExtension());
			//
			//	if (quality != Quality_Unknown)
			//	{
			//		result.SourceDetectionSource = QualityDetectionSource.Extension;
			//		source = Quality_Source;
			//	}
			//}
			//catch (ArgumentException ex)
			//{
			//Logger.Debug(ex, "Unable to parse quality from extension");
			//}
		}

		if resolution == Resolution_2160p {
			result.Resolution = Resolution_2160p

			if result.QualitySource == Type_Unknown {
				result.QualitySource = Type_WebDl
				result.QualityResolution = Quality_WEBDL2160p
			} else {
				result.QualityResolution = result.QualitySource + "-2160p"
			}
			return true
		}
		if resolution == Resolution_1080p {
			result.Resolution = Resolution_1080p
			if result.QualitySource == Type_Unknown {
				result.QualitySource = Type_WebDl
				result.QualityResolution = Quality_WEBDL1080p
			} else {
				result.QualityResolution = result.QualitySource + "-1080p"
			}
			return true
		}
		if resolution == Resolution_720p {
			result.Resolution = Resolution_720p
			if result.QualitySource == Type_Unknown {
				result.QualitySource = Type_WebDl
				result.QualityResolution = Quality_WEBDL720p
			} else {
				result.QualityResolution = result.QualitySource + "-720p"
			}
			return true
		}
		if resolution == Resolution_360p || resolution == Resolution_480p {
			result.Resolution = Resolution_360p

			if result.QualitySource == Type_Unknown {
				result.QualitySource = Type_SDTV

				result.QualityResolution = Quality_SDTV
			} else {
				result.QualityResolution = result.QualitySource + "-480p"
			}
			return true
		}

	}
	return false
}

func (this *TorrentTitleParse) matchContains(result *MatchResult) bool {

	// 找不到尺寸
	if strings.Contains(result.AnalyzeTitle, "848x480") {
		if strings.Contains(result.AnalyzeTitle, "dvd") {
			//result.SourceDetectionSource = QualityDetectionSource.Name
			result.QualitySource = Type_DVD
			result.Resolution = Resolution_Unknown
			result.QualityResolution = Quality_DVD
		} else if strings.Contains(strings.ToLower(result.AnalyzeTitle), "bluray") {
			//result.SourceDetectionSource = QualityDetectionSource.Name
			//result.Quality = Quality_Bluray480p
			result.QualitySource = Type_Bluray
			result.Resolution = Resolution_480p
			result.QualityResolution = Quality_Bluray480p
		} else {
			//result.Quality = Quality_SDTV
			result.QualitySource = Type_HDTV
			result.Resolution = Resolution_Unknown
			result.QualityResolution = Quality_SDTV
		}

		return true
	}

	if strings.Contains(strings.ToLower(result.AnalyzeTitle), "1280x720") {
		result.Resolution = Resolution_720p

		if strings.Contains(strings.ToLower(result.AnalyzeTitle), "bluray") {
			result.QualitySource = Type_Bluray
			result.QualityResolution = Quality_Bluray720p
		} else {
			result.QualitySource = Type_HDTV
			result.QualityResolution = Quality_HDTV720p
		}

		return true
	}

	if strings.Contains(strings.ToLower(result.AnalyzeTitle), "1920x1080") {
		result.Resolution = Resolution_1080p

		if strings.Contains(strings.ToLower(result.AnalyzeTitle), "bluray") {

			result.QualitySource = Type_Bluray
			result.QualityResolution = Quality_Bluray1080p
		} else {
			result.QualitySource = Type_HDTV
			result.QualityResolution = Quality_HDTV1080p
		}

		return true
	}

	if strings.Contains(strings.ToLower(result.AnalyzeTitle), "bluray720p") {
		result.QualitySource = Type_Bluray
		result.Resolution = Resolution_720p
		result.QualityResolution = Quality_Bluray720p

		return true
	}

	if strings.Contains(strings.ToLower(result.AnalyzeTitle), "bluray1080p") {
		result.QualitySource = Type_Bluray
		result.Resolution = Resolution_1080p
		result.QualityResolution = Quality_Bluray1080p

		return true
	}

	if strings.Contains(strings.ToLower(result.AnalyzeTitle), "bluray2160p") {

		result.QualitySource = Type_Bluray
		result.Resolution = Resolution_2160p
		result.QualityResolution = Quality_Bluray2160p

		return true
	}

	return false
}

// 匹配质量
func (this *TorrentTitleParse) ParseResolution(result *MatchResult) string {
	regGroups := regexp_ext.ParseGroups(resolutionRegex, result.AnalyzeTitle)
	if regGroups == nil {
		return Resolution_Unknown
	}

	// 删除搜索到的内容
	result.AnalyzeTitle, _ = resolutionRegex.Replace(result.AnalyzeTitle, "", -1, -1)

	if len(regGroups.GetGroupValByName("R360p")) > 0 {
		return Resolution_360p
	}
	if len(regGroups.GetGroupValByName("R480p")) > 0 {
		return Resolution_480p
	}
	if len(regGroups.GetGroupValByName("R576p")) > 0 {
		return Resolution_576p
	}
	if len(regGroups.GetGroupValByName("R720p")) > 0 {
		return Resolution_720p
	}
	if len(regGroups.GetGroupValByName("R1080p")) > 0 {
		return Resolution_1080p
	}
	if len(regGroups.GetGroupValByName("R2160p")) > 0 {
		return Resolution_2160p
	}

	return Resolution_Unknown
}

// 判断是否为蓝光原盘
func (this *TorrentTitleParse) ParseRemux(result *MatchResult) bool {
	regex := regexp2.MustCompile(`\b(?<remux>(BD)?[-_. ]?Remux)\b`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)
	res, err := regex.MatchString(result.AnalyzeTitle)
	if err != nil || res == false {
		return false
	}
	result.AnalyzeTitle, _ = regex.Replace(result.AnalyzeTitle, "", -1, -1)
	return res
}

func (this *TorrentTitleParse) OtherSourceMatch(result *MatchResult) bool {
	regGroups := regexp_ext.ParseGroups(OtherSourceRegex, result.AnalyzeTitle)
	if regGroups == nil {
		return false
	}

	if len(regGroups.GetGroupValByName("sdtv")) > 0 {

		result.QualitySource = Type_HDTV
		result.Resolution = Resolution_Unknown
		result.QualityResolution = Quality_SDTV

		return true // Quality_SDTV
	}
	if len(regGroups.GetGroupValByName("hdtv")) > 0 {
		result.QualitySource = Type_HDTV
		result.Resolution = Resolution_720p
		result.QualityResolution = Quality_HDTV720p
		return true //Quality_HDTV720p
	}

	return false // Quality_Unknown
}
