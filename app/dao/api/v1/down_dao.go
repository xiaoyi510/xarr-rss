package v1

import (
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/download"
	"XArr-Rss/app/service/sources"
	"XArr-Rss/app/service/torrent"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"bytes"
	"encoding/base32"
	"encoding/hex"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/jackpal/bencode-go"
	"log"
	"strings"
)

type Apiv1DownDao struct {
}

// Sonarr开始下载种子
func (this Apiv1DownDao) Down(c *gin.Context) {
	url := c.Query("url")
	sourceId := c.Query("source_id")
	title := c.Query("title")
	if title == "" {
		log.Println("获取标题失败")
		c.JSON(502, gin.H{})
		return
	}
	repMh := regexp2.MustCompile(`\s*?[:：]\s*`,
		regexp2.IgnoreCase|regexp2.Compiled|regexp2.IgnorePatternWhitespace)
	title, _ = repMh.Replace(title, " - ", -1, -1)

	//title = helper.StrReplace(title, []string{
	//	" : ",
	//	": ",
	//	":",
	//}, []string{
	//	" - ",
	//	" - ",
	//	"-",
	//})

	// 查询代理信息
	proxy := ""
	if sourceId != "" {
		// 查询数据源是否使用代理
		sourceInfo, err := sources.SourcesService{}.GetSourcesInfo(sourceId)
		if err != nil {
			log.Println("查询数据源错误", err.Error())
			c.JSON(502, gin.H{})
			return
		}
		if sourceInfo.UseProxy == 1 {
			proxy = options.GetOption(options.OptionsGlobalProxy)
		}
		if sourceInfo.DownloadPasskey != "" {
			if !strings.Contains(url, "passkey") {
				if strings.Contains(url, "?") {
					url += "&passkey=" + sourceInfo.DownloadPasskey
				} else {
					url += "?passkey=" + sourceInfo.DownloadPasskey
				}
			}
		}
	} else {
		log.Println("没有查询到数据源ID,无法使用代理")
	}
	if url == "" {
		log.Println("下载地址为空")
		c.JSON(502, gin.H{})
		return
	}
	log.Println("开始试图下载种子文件", url, title, c.GetHeader("user-agent"))
	hash := this.ParseTorrentHash(url, proxy, variable.ServerState.IsVip)
	if len(hash) == 0 {
		c.JSON(404, gin.H{})
		return
	}

	downloadInfo := download.DownloadService{}.GetDownInfo(hash)
	if downloadInfo != nil {
		download.DownloadService{}.DelDownInfo(downloadInfo)
		// 修改为重新监听
		//downloadInfo.Status = dbmodel.DownloadListStatusWait
		//downloadInfo.Process = 0
		//downloadInfo.Title = title
		//downloadInfo.CreatedAt = time.Now()
		//downloadInfo.UpdatedAt = time.Now()
		//downloadInfo.DownClient = ""
		log.Println("发现已有种子被监听,但是重新下载了,本次将监听任务删除后重新添加", url, title)

		//err := download.DownloadService{}.SaveDownloadInfo(downloadInfo)
		//if err != nil {
		//	log.Println("修改监听任务失败了")
		//}
		//c.Redirect(302, url)
		//return
	}
	//if strings.Index(c.GetHeader("user-agent"), "Sonarr/") > -1 {
	//	c.Redirect(302, url)
	//	return
	//}
	// 保存到任务
	download.DownloadService{}.AddDownloadInfo(&dbmodel.DownloadList{
		Hash:    hash,
		Title:   title,
		Status:  dbmodel.DownloadListStatusWait,
		Process: 0,
	})

	logsys.Info("已加入下载监听:"+hash, "监听")
	c.Redirect(302, url)
}

// 解析种子hash 多个已逗号分隔
func (this Apiv1DownDao) ParseTorrentHash(uri, proxy string, mustUseQb bool) string {
	// 取出hash
	if !mustUseQb {
		//compile, err := regexp.Compile(`[retData-z0-9A-Z]{40}`)
		//if err != nil {
		//	return ""
		//}
		//hash := compile.FindString(uri)
		hash := download.DownloadService{}.GetUrlTorrentHash(uri)
		if hash != "" {
			return hash
		}
	}

	if !variable.ServerState.IsVip {
		logsys.Info("非赞助会员用户,不可匹配磁力连接", "下载")
		return ""
	}

	// 判断是否为磁链
	if strings.Contains(uri, "magnet:") {
		// 直接解析hahs
		// 提取 urn:btih:
		// magnet:?xt=urn:btih:ZZ3HO4MWSKGZMAFXHQPTTW6JB24HGRFZ&amp;dn=&amp;tr=http%3A%2F%2F104.143.10.186%3A8000%2Fannounce&amp;tr=udp%3A%2F%2F104.143.10.186%3A8000%2Fannounce&amp;tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&amp;tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&amp;tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&amp;tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.prq.to%2Fannounce&amp;tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce&amp;tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud&amp;tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&amp;tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&amp;tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&amp;tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&amp;tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&amp;tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&amp;tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&amp;tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce&amp;tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce&amp;tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce&amp;tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce&amp;tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce&amp;tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce&amp;tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce&amp;tr=http%3A%2F%2Ft.nyaatracker.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Fanidex.moe%3A6969%2Fannounce&amp;tr=https%3A%2F%2Fopentracker.acgnx.se%2Fannounce&amp;tr=http%3A%2F%2Fopentracker.acgnx.com%3A6869%2Fannounce&amp;tr=http%3A%2F%2Fopen.miotracker.com%2Fannounce&amp;tr=http%3A%2F%2Ftracker.anirena.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce&amp;tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce&amp;tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&amp;tr=http%3A%2F%2F0205.uptm.ch%3A6969%2Fannounce&amp;tr=http%3A%2F%2F1337.abcvg.info%3A80%2Fannounce&amp;tr=http%3A%2F%2F178.175.143.27%3A80%2Fannounce&amp;tr=http%3A%2F%2F78.30.254.12%3A2710%2Fannounce&amp;tr=http%3A%2F%2F91.217.91.21%3A3218%2Fannounce&amp;tr=http%3A%2F%2F93.92.64.5%3A80%2Fannounce&amp;tr=http%3A%2F%2F%5B2001%3A1b10%3A1000%3A8101%3A0%3A242%3Aac11%3A2%5D%3A6969%2Fannounce&amp;tr=http%3A%2F%2F%5B2001%3A470%3A1%3A189%3A0%3A1%3A2%3A3%5D%3A6969%2Fannounce&amp;tr=http%3A%2F%2F%5B2a04%3Aac00%3A1%3A3dd8%3A%3A1%3A2710%5D%3A2710%2Fannounce&amp;tr=http%3A%2F%2Faaa.army%3A8866%2Fannounce&amp;tr=http%3A%2F%2Fankeschwarz.net%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fatrack.pow7.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Fbandari.org%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fbobbialbano.com%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fbt.pusacg.org%3A8080%2Fannounce&amp;tr=http%3A%2F%2Fdn42.smrsh.net%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fexplodie.org%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fgrifon.info%3A80%2Fannounce&amp;tr=http%3A%2F%2Fh4.trakx.nibba.trade%3A80%2Fannounce&amp;tr=http%3A%2F%2Fipv6-publictracker.zooki.xyz%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fns349743.ip-91-121-106.eu%3A80%2Fannounce&amp;tr=http%3A%2F%2Fopen.acgnxtracker.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Fopen.touki.ru%3A80%2Fannounce.php&amp;tr=http%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fpow7.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Fretracker.hotplug.ru%3A2710%2Fannounce&amp;tr=http%3A%2F%2Fretracker.krs-ix.ru%3A80%2Fannounce&amp;tr=http%3A%2F%2Fretracker.sevstar.net%3A2710%2Fannounce&amp;tr=http%3A%2F%2Fretracker.spark-rostov.ru%3A80%2Fannounce&amp;tr=http%3A%2F%2Frt.tace.ru%3A80%2Fannounce&amp;tr=http%3A%2F%2Fsecure.pow7.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Ft1.pow7.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Ft2.pow7.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Fthetracker.org%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftorrentsmd.com%3A8080%2Fannounce&amp;tr=http%3A%2F%2Ftorrenttracker.nwc.acsalaska.net%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftr.cili001.com%3A8070%2Fannounce&amp;tr=http%3A%2F%2Ftracker.anonwebz.xyz%3A8080%2Fannounce&amp;tr=http%3A%2F%2Ftracker.birkenwald.de%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker.bittor.pw%3A1337%2Fannounce&amp;tr=http%3A%2F%2Ftracker.bt4g.com%3A2095%2Fannounce&amp;tr=http%3A%2F%2Ftracker.bz%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker.dutchtracking.nl%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.files.fm%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker.gbitt.info%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.ipv6tracker.ru%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.kuroy.me%3A5944%2Fannounce&amp;tr=http%3A%2F%2Ftracker.lelux.fi%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.moeking.me%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker.noobsubs.net%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.nyap2p.com%3A8080%2Fannounce&amp;tr=http%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&amp;tr=http%3A%2F%2Ftracker.skyts.net%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker.sloppyta.co%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftracker.ygsub.com%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker.yoshi210.com%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker.zerobytes.xyz%3A1337%2Fannounce&amp;tr=http%3A%2F%2Ftracker.zum.bi%3A6969%2Fannounce&amp;tr=http%3A%2F%2Ftracker2.dler.org%3A80%2Fannounce&amp;tr=http%3A%2F%2Ftrun.tom.ru%3A80%2Fannounce&amp;tr=http%3A%2F%2Fvpn.flying-datacenter.de%3A6969%2Fannounce&amp;tr=http%3A%2F%2Fvps02.net.orel.ru%3A80%2Fannounce&amp;tr=http%3A%2F%2Fvps2.avc.cx%3A7171%2Fannounce&amp;tr=http%3A%2F%2Fwww.wareztorrent.com%3A80%2Fannounce&amp;tr=https%3A%2F%2F1337.abcvg.info%3A443%2Fannounce&amp;tr=https%3A%2F%2F2.tracker.eu.org%3A443%2Fannounce&amp;tr=https%3A%2F%2F3.tracker.eu.org%3A443%2Fannounce&amp;tr=https%3A%2F%2Faaa.army%3A8866%2Fannounce&amp;tr=https%3A%2F%2Fatrack-asia-s1.114913.xyz%3A443%2Fannounce&amp;tr=https%3A%2F%2Fopen.kickasstracker.com%3A443%2Fannounce&amp;tr=https%3A%2F%2Fopentracker.acgnx.se%3A443%2Fannounce&amp;tr=https%3A%2F%2Fpublictracker.pp.ua%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.bt-torrentHash.com%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.gbitt.info%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.hama3.net%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.imgoingto.icu%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.lelux.fi%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.nitrix.me%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.sloppyta.co%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftracker.tamersunion.org%3A443%2Fannounce&amp;tr=https%3A%2F%2Ftrakx.herokuapp.com%3A443%2Fannounce&amp;tr=https%3A%2F%2Fw.wwwww.wtf%3A443%2Fannounce
		err, querys := helper.ParseUrl(uri)
		if err != nil {
			logsys.Error("磁力解析失败:%s", "下载", err.Error())
			return ""
		}
		xts, ok := querys["xt"]
		if !ok {
			logsys.Error("磁力找到xt失败", "下载")
			return ""
		}
		ret := []string{}
		for _, xt := range xts {
			if strings.Index(xt, "urn:btih:") == 0 {
				logsys.Info("找到磁力信息v1:%v", "磁力", xt)

				//v1
				xt = strings.ReplaceAll(xt, "urn:btih:", "")
				if len(xt) != 40 {
					logsys.Info("磁力解析2:%v", "磁力", xt)
					sha1Data, err := base32.StdEncoding.DecodeString(xt)
					if err != nil {
						logsys.Error("磁力信息解密失败:%s", "磁力", err.Error())
						//return ""
					} else {
						ret = append(ret, hex.EncodeToString(sha1Data))
					}

				} else {
					ret = append(ret, xt)
				}
			} else if strings.Index(xt, "urn:btmh:") == 0 {
				//v2
				xt = strings.ReplaceAll(xt, "urn:btmh:", "")
				// 判断长度是否需要处理
				xtLen := len(xt)

				if xtLen == 68 {
					// 切割前4位
					enType := xt[0:4]
					if enType == "1220" {
						ret = append(ret, xt[4:])
					} else {
						logsys.Error("磁力信息匹配失败:%s", "磁力", xt)
					}
				} else if xtLen == 64 {
					ret = append(ret, xt)
				} else {
					// 匹配失败
					logsys.Error("磁力信息匹配失败:%s", "磁力", xt)
				}
			}
		}

		return strings.Join(ret, ",")

	}

	return this.ParseOnlineTorrent(uri, proxy)
}

func (this Apiv1DownDao) ParseOnlineTorrent(uri, proxy string) string {
	defer func() {
		if e := recover(); e != nil {
			logsys.Error("err:%v", "格式化在线种子", e)
		}
	}()
	// 下载种子
	curlHelper := &helper.CurlHelper{}
	err, result, _ := curlHelper.Init(nil).SetProxy(proxy).Get(uri, nil, true)
	//err, result, _ := helper.CurlHelper{}.GetUri(uri, nil, nil, false)
	if err != nil {
		logsys.Error("获取种子数据错误["+uri+"]:"+err.Error(), "种子")
		return ""
	}
	if len(result) <= 0 {
		logsys.Error("获取种子数据错误!", "种子")
		return ""
	}
	hashStr := []string{}
	a := this.GetHash(result, "v1")
	if a != "" {
		hashStr = append(hashStr, a)
	}
	// 提取v2 种子hash
	a = this.GetHash(result, "v2")
	if a != "" {
		hashStr = append(hashStr, a)
	}
	return strings.Join(hashStr, ",")
}

func (this Apiv1DownDao) GetHash(result []byte, v string) string {

	t, err := bencode.Decode(bytes.NewReader(result))
	if err != nil {
		logsys.Error("解码种子信息异常:%s", "种子", err.Error())
		return ""
	}

	a := t.(map[string]interface{})["info"]
	hashStr := ""
	if v == "" || v == "v1" {
		hashStr, _ = torrent.TorrentHash(a)
	} else {
		hashStr, _ = torrent.TorrentHashV2(a)
	}
	if len(hashStr) != 0 {
		// 保存种子
		return hashStr
	}
	return ""
}
