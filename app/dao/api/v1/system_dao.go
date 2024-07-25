package v1

import (
	"XArr-Rss/app/components/torrent_title_parse"
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/app/global/variable"
	groups2 "XArr-Rss/app/service/groups"
	sources2 "XArr-Rss/app/service/sources"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

type Apiv1SystemDao struct {
}
type Apiv1SystemSaveReq struct {
	HttpAddr                    string `json:"http_addr"`
	HttpListenPort              int    `json:"http_listen_port"`
	GlobalProxy                 string `json:"global_proxy"`
	UserAgent                   string `json:"user_agent"`
	EchoTitleAnime              string `json:"echo_title_anime"`
	EchoTitleTv                 string `json:"echo_title_tv"`
	BackupDatabaseCount         int    `json:"backup_database_count"`
	BackupDatabaseTime          int64  `json:"backup_database_time"`
	SonarrUnmonitoredRmGmedia   string `json:"sonarr_unmonitored_rm_gmedia"`
	SonarrUnmonitoredNaddGmedia string `json:"sonarr_unmonitored_nadd_gmedia"`
}

type Apiv1SystemEditWordsRuleReq struct {
	WordsRule string `json:"words_rule"`
}

type Apiv1SystemGetRes struct {
	HttpAddr                    string `json:"http_addr"`
	HttpListenPort              int    `json:"http_listen_port"`
	GlobalProxy                 string `json:"global_proxy"`
	UserAgent                   string `json:"user_agent"`
	EchoTitleAnime              string `json:"echo_title_anime"`
	EchoTitleTv                 string `json:"echo_title_tv"`
	BackupDatabaseCount         int    `json:"backup_database_count"`
	BackupDatabaseTime          int64  `json:"backup_database_time"`
	SonarrUnmonitoredRmGmedia   string `json:"sonarr_unmonitored_rm_gmedia"`
	SonarrUnmonitoredNaddGmedia string `json:"sonarr_unmonitored_nadd_gmedia"`
}

func (d Apiv1SystemDao) Save(c *gin.Context) {
	var req Apiv1SystemSaveReq
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "请求参数错误",
		})
		return
	}

	if strings.Index(req.HttpAddr, "http") == -1 {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "映射后的地址必须为http[s]://开头",
		})
		return
	}
	if req.HttpListenPort <= 0 {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "端口号必须大于0",
		})
		return
	}
	appconf.AppConf.System.Port = req.HttpListenPort
	appconf.AppConf.System.GlobalProxy = req.GlobalProxy
	appconf.AppConf.System.UserAgent = req.UserAgent
	appconf.AppConf.System.EchoTitleAnime = req.EchoTitleAnime
	appconf.AppConf.System.EchoTitleTv = req.EchoTitleTv
	appconf.AppConf.System.BackupDatabaseCount = req.BackupDatabaseCount
	appconf.AppConf.System.BackupDatabaseTime = req.BackupDatabaseTime
	appconf.AppConf.System.SonarrUnmonitoredNaddGmedia = req.SonarrUnmonitoredNaddGmedia
	appconf.AppConf.System.SonarrUnmonitoredRmGmedia = req.SonarrUnmonitoredRmGmedia
	appconf.AppConf.System.HttpAddr = strings.TrimRight(req.HttpAddr, "/")

	err = appconf.AppConf.System.Save()
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
		"message": "保存成功,如果修改了端口号,请重启软件",
	})
}

// 导出配置项
func (d Apiv1SystemDao) Export(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	// 查询所有分组列表
	groups := groups2.GroupsService{}.GetGroupList()
	groupMedias := groups2.GroupMediaService{}.GetGroupMediaList()
	sources := sources2.SourcesService{}.GetSourcesList(false)

	// 进入页面重载配置
	appconf.AppConf.System.Reload()
	c.JSON(200, gin.H{
		"code":    0,
		"data":    gin.H{"system": appconf.AppConf.System, "service": appconf.AppConf.Service, "groups": groups, "group_medias": groupMedias, "sources": sources},
		"message": "success",
	})
}

// 导入配置项
func (d Apiv1SystemDao) Import(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	// 进入页面重载配置
	//appconf.AppConf.System.Reload()
	//c.JSON(200, gin.H{
	//	"code":    0,
	//	"data":    gin.H{"System": appconf.AppConf.System, "Service": appconf.AppConf.Service},
	//	"message": "success",
	//})
}

func (d Apiv1SystemDao) Get(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	// 进入页面重载配置
	appconf.AppConf.System.Reload()

	res := &Apiv1SystemGetRes{}
	res.HttpAddr = appconf.AppConf.System.HttpAddr
	res.GlobalProxy = appconf.AppConf.System.GlobalProxy
	res.HttpListenPort = appconf.AppConf.System.Port
	res.UserAgent = appconf.AppConf.System.UserAgent
	res.EchoTitleAnime = appconf.AppConf.System.EchoTitleAnime
	res.EchoTitleTv = appconf.AppConf.System.EchoTitleTv
	res.BackupDatabaseTime = appconf.AppConf.System.BackupDatabaseTime
	res.BackupDatabaseCount = appconf.AppConf.System.BackupDatabaseCount
	res.SonarrUnmonitoredNaddGmedia = appconf.AppConf.System.SonarrUnmonitoredNaddGmedia
	res.SonarrUnmonitoredRmGmedia = appconf.AppConf.System.SonarrUnmonitoredRmGmedia

	c.JSON(200, gin.H{
		"code":    0,
		"data":    res,
		"message": "success",
	})
}

func (d Apiv1SystemDao) GetWordsRule(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	// 进入页面重载配置
	appconf.AppConf.System.Reload()
	c.JSON(200, gin.H{
		"code":    0,
		"data":    appconf.AppConf.System.WordsRule,
		"message": "success",
	})
}

func (d Apiv1SystemDao) EditWordsRule(c *gin.Context) {
	var req Apiv1SystemEditWordsRuleReq

	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "请求参数错误",
		})
		return
	}
	// 校验值是否正确
	t := torrent_title_parse.TorrentTitleParse{}
	if err2 := t.ValidateRules(req.WordsRule); err2 != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    nil,
			"message": "规则错误:" + err2.Error(),
		})
		return
	}

	options.SetOption(options.WordsRule, req.WordsRule)

	appconf.AppConf.System.WordsRule = req.WordsRule

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "保存成功",
	})
}

// 检查token
func checkAuth(context *gin.Context) error {
	_, err := os.Stat("./conf/debug.txt")
	if err == nil {
		return nil
	}

	cookie, err := context.Request.Cookie("xarn_rss_token" + variable.HttpTokenAfter)
	if err != nil {
		log.Println("请清除浏览器缓存后重试", err)
		//variable.HttpTokenAfter += "b"
		return errors.New("请登录")
	}
	if !variable.CheckToken(cookie.Value) {
		return errors.New("登录已过期")
	}
	return nil
}
