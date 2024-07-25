package v1

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/sdk/qbit"
	"XArr-Rss/app/sdk/qbit/client"
	"XArr-Rss/app/sdk/themoviedb"
	"XArr-Rss/app/sdk/transmission"
	"XArr-Rss/app/service/medias"
	"XArr-Rss/app/service/sonarr"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Apiv1ServiceDao struct {
}
type Apiv1ServiceDaoSaveReq struct {
	SonarrHost            string `json:"sonarr_host,omitempty"`
	SonarrApikey          string `json:"sonarr_apikey,omitempty"`
	SonarrVersion         string `json:"sonarr_version,omitempty"`
	SonarrRefreshTime     string `json:"sonarr_refresh_time,omitempty"`
	ThemoviedbApikey      string `json:"themoviedb_apikey,omitempty"`
	ThemoviedbProxyType   string `json:"themoviedb_proxy_type,omitempty"`
	ThemoviedbProxy       string `json:"themoviedb_proxy,omitempty"`
	QbittorrentHost       string `json:"qbittorrent_host,omitempty"`
	QbittorrentUsername   string `json:"qbittorrent_username,omitempty"`
	QbittorrentPassword   string `json:"qbittorrent_password,omitempty"`
	QbittorrentRenameFile string `json:"qbittorrent_rename_file,omitempty"`

	TransmissionHost     string `json:"transmission_host"`
	TransmissionUsername string `json:"transmission_username"`
	TransmissionPassword string `json:"transmission_password"`
}

func (d Apiv1ServiceDao) Save(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	var req Apiv1ServiceDaoSaveReq

	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "请求参数错误:" + err.Error(),
		})
		return
	}

	if strings.Index(req.SonarrHost, "http") == -1 {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "SonarrId 地址 必须为http[s]://开头",
		})
		return
	}
	if req.QbittorrentHost != "" {
		if strings.Index(req.QbittorrentHost, "http") == -1 {
			c.JSON(200, gin.H{
				"code":    500,
				"data":    nil,
				"message": "qBittorrent 地址 必须为http[s]://开头",
			})
			return
		}
		_, err = url.Parse(req.QbittorrentHost)
		if err != nil {
			c.JSON(200, gin.H{
				"code":    500,
				"data":    nil,
				"message": "qbit地址解析错误:" + err.Error(),
			})
			return
		}
	}

	appconf.AppConf.Service.Sonarr.Host = strings.TrimRight(req.SonarrHost, "/")
	appconf.AppConf.Service.Sonarr.Apikey = req.SonarrApikey
	appconf.AppConf.Service.Sonarr.Version = req.SonarrVersion
	appconf.AppConf.Service.Sonarr.RefreshTime = helper.StrToInt(req.SonarrRefreshTime)
	appconf.AppConf.Service.Themoviedb.Apikey = req.ThemoviedbApikey
	appconf.AppConf.Service.Themoviedb.Proxy = req.ThemoviedbProxy
	appconf.AppConf.Service.Themoviedb.ProxyType = helper.StrToInt(req.ThemoviedbProxyType)

	appconf.AppConf.Service.Qbittorrent.Host = req.QbittorrentHost
	appconf.AppConf.Service.Qbittorrent.Username = req.QbittorrentUsername
	appconf.AppConf.Service.Qbittorrent.Password = req.QbittorrentPassword
	appconf.AppConf.Service.Qbittorrent.RenameFile = req.QbittorrentRenameFile

	appconf.AppConf.Service.Transmission.Host = req.TransmissionHost
	appconf.AppConf.Service.Transmission.Username = req.TransmissionUsername
	appconf.AppConf.Service.Transmission.Password = req.TransmissionPassword

	err = appconf.AppConf.Service.Save()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "保存失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "保存成功",
	})
}

func (d Apiv1ServiceDao) Get(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	ret := Apiv1ServiceDaoSaveReq{
		SonarrHost:            appconf.AppConf.Service.Sonarr.Host,
		SonarrApikey:          appconf.AppConf.Service.Sonarr.Apikey,
		SonarrVersion:         appconf.AppConf.Service.Sonarr.Version,
		SonarrRefreshTime:     strconv.Itoa(appconf.AppConf.Service.Sonarr.RefreshTime),
		ThemoviedbApikey:      appconf.AppConf.Service.Themoviedb.Apikey,
		ThemoviedbProxy:       appconf.AppConf.Service.Themoviedb.Proxy,
		ThemoviedbProxyType:   strconv.Itoa(appconf.AppConf.Service.Themoviedb.ProxyType),
		QbittorrentHost:       appconf.AppConf.Service.Qbittorrent.Host,
		QbittorrentUsername:   appconf.AppConf.Service.Qbittorrent.Username,
		QbittorrentPassword:   appconf.AppConf.Service.Qbittorrent.Password,
		QbittorrentRenameFile: appconf.AppConf.Service.Qbittorrent.RenameFile,
		TransmissionHost:      appconf.AppConf.Service.Transmission.Host,
		TransmissionUsername:  appconf.AppConf.Service.Transmission.Username,
		TransmissionPassword:  appconf.AppConf.Service.Transmission.Password,
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    ret,
		"message": "success",
	})
}

func (d Apiv1ServiceDao) SonarrReload(c *gin.Context) {
	err := sonarr.SonarrService{}.SyncSonarrToLocal()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "刷新失败:" + err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code":    0,
			"data":    nil,
			"message": "刷新成功",
		})
	}
}
func (d Apiv1ServiceDao) SonarrReloadReal(c *gin.Context) {
	err := sonarr.SonarrService{}.SyncSonarrToLocal(true)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "刷新失败:" + err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code":    0,
			"data":    nil,
			"message": "刷新成功",
		})
	}
}

func (d Apiv1ServiceDao) SonarrTest(c *gin.Context) {
	err, sonarrStatus := sonarr.SonarrService{}.TestConnect()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "连接失败:" + err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code":    0,
			"data":    nil,
			"message": "连接成功:" + sonarrStatus.AppName + "(" + sonarrStatus.Version + ")",
		})
	}
}

func (d Apiv1ServiceDao) TmdbTest(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	if appconf.AppConf.Service.Themoviedb.Apikey != "" {
		// 获取tmdb 代理
		proxy := appconf.AppConf.Service.GetTMDBProxy()
		sdk := themoviedb.GetTheMoviedbSdk(appconf.AppConf.Service.Themoviedb.Apikey, proxy)
		err, _ := sdk.GetConfigurationRs()
		if err != nil {
			c.JSON(200, gin.H{
				"code":    502,
				"data":    nil,
				"message": "请求失败:" + err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"data":    nil,
			"message": "请求成功",
		})

	} else {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请输入ThemovieDb ApiKey",
		})
	}
}

func (d Apiv1ServiceDao) QbitTest(c *gin.Context) {
	err := checkQbit()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求失败:" + err.Error(),
		})
		return
	}

	qbit.Qbit.Client = &client.QbitClient{}
	qbit.Qbit.Client.Init(appconf.AppConf.Service.Qbittorrent.Host)
	// 解决qbit客户端保存未生效问题
	qbit.Qbit.Client.SetConf(appconf.AppConf.Service.Qbittorrent.Host)
	err = qbit.Qbit.Login(appconf.AppConf.Service.Qbittorrent.Username, appconf.AppConf.Service.Qbittorrent.Password)
	if err == nil {
		c.JSON(200, gin.H{
			"code":    0,
			"data":    nil,
			"message": "请求成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    502,
		"data":    nil,
		"message": "请求失败:" + err.Error(),
	})
}

func (d Apiv1ServiceDao) TmdbReload(c *gin.Context) {
	go medias.SyncMediaRemoteInfo()

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "任务已提交,请稍后查看",
	})
}

// 强制全部重新加载
func (d Apiv1ServiceDao) TmdbReloadRefresh(c *gin.Context) {
	go medias.SyncMediaRemoteInfo(true)

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "任务已提交,请稍后查看",
	})
}

func checkQbit() error {
	// 判断Qb是否输入了Key
	if appconf.AppConf.Service.Qbittorrent.Host == "" {
		return logsys.Error("请输入Qbit主机信息", "qbit")
	}
	if !strings.Contains(appconf.AppConf.Service.Qbittorrent.Host, "http") {
		return logsys.Error("请输入Qbit主机信息 如: http://xxx.xxx:6544", "qbit")
	}

	return nil

}

// 测试是否正常
func (d Apiv1ServiceDao) TransmissionTest(c *gin.Context) {
	tr := transmission.Transmission{}.GetClient(appconf.AppConf.Service.Transmission.Host, appconf.AppConf.Service.Transmission.Username, appconf.AppConf.Service.Transmission.Password)
	err := tr.Init()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "连接Transmission错误:" + err.Error(),
		})
	}

	err, version := tr.SessionGetVersion()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "连接Transmission错误:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "Transmission 版本号:" + version,
	})
}
