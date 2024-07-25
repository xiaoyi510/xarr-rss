package torrent_title_parse

import (
	"XArr-Rss/util/regexp_ext"
	"github.com/dlclark/regexp2"
)

// 视频编码识别：AVC,H264,x264,HEVC,H265,x265
//色深（一般标注在视频编码之后）：8bit,10bit
//音频编码：AAC,FLAC,OPUS
var hevcReg = regexp2.MustCompile(`\b(?<aac>aac)(?:x\d)?|
(?<ass>ass)|
(?<ddp>ddp5.1)|
(?<atmos>atmos)|
(?<hevc>hevc(?:[-_](?:8|10)bit)?)|
(?<flac>flac)|
(?<opus>opus)|
(?<h264>h264|x264)|
(?<h265>h265|x265)|
(?<xvid>xvid)(?:\-afg)?|
(?<avc>avc(?:[-_\s](?:8|10)bit)?)
\b`, regexp2.IgnoreCase|regexp2.IgnorePatternWhitespace|regexp2.Compiled)

// 解析出来hevc acc
func (this *TorrentTitleParse) ParseHevcAcc(result *MatchResult) {
	matchString, err := hevcReg.MatchString(result.AnalyzeTitle)
	if err != nil || matchString == false {
		return
	}

	hevcGroup := regexp_ext.ParseGroups(hevcReg, result.AnalyzeTitle)
	if hevcGroup == nil {
		return
	}
	// 音频
	for _, v := range []string{"aac", "flac", "opus", "ddp"} {
		audio := hevcGroup.GetGroupValByName(v)
		if len(audio) > 0 {
			result.AudioEncode = append(result.AudioEncode, audio...)
		}
	}

	// 视频
	for _, v := range []string{"avc", "hevc", "h264", "h265", "xvid", "atmos"} {
		video := hevcGroup.GetGroupValByName(v)
		if len(video) > 0 {
			result.VideoEncode = append(result.VideoEncode, video...)
		}
	}

	result.AnalyzeTitle, _ = hevcReg.Replace(result.AnalyzeTitle, "", -1, -1)

}
