package qbit

import (
	"XArr-Rss/app/global/conf/options"
	"XArr-Rss/app/sdk/qbit/api"
	"XArr-Rss/app/sdk/qbit/client"
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var Qbit = qbitApi{}

type qbitApi struct {
	Client   *client.QbitClient
	Auth     api.Auth
	App      api.App
	Torrents api.Torrents
	Sync     api.Sync
	IsLogin  bool
}

func (this *qbitApi) GetApp() *api.App {
	return (&api.App{}).SetClient(this.Client)
}
func (this *qbitApi) GetTorrents() *api.Torrents {
	return (&api.Torrents{}).SetClient(this.Client)
}

func (this *qbitApi) GetAuth() *api.Auth {
	return (&api.Auth{}).SetClient(this.Client)
}
func (this *qbitApi) GetSync() *api.Sync {
	return (&api.Sync{}).SetClient(this.Client)
}

// Login 第一次运行 同步所有的数据
func (this *qbitApi) Login(username, password string) error {
	//>> 登录获取Cookie
	res, ck := this.GetAuth().Login(username, password)
	if res == "Ok." {
		//>> 写出Cookie
		err := options.SetOption(options.Qbittorrent_Cookie, ck)
		//err := os.WriteFile(appconf.AppConf.ConfDir+"/qbit.ck.txt", []byte(ck), 0666)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("登录失败:" + res)
}

// CheckLogin 检测登录是否有效
func (this *qbitApi) CheckLogin(username, password string) bool {
	//if username == "" || password == "" {
	//	log.Println("请输入Qbit账号密码")
	//	return false
	//}
	//log.Println("开始检测Cookie")
	// 判断缓存是否可用
	if this.CheckCookie() == false {
		log.Println("Cookie无效 开始使用账号密码登录")
		err := this.Login(username, password)
		if err != nil {
			log.Println("登录失败", err)
			this.IsLogin = false
			if err.Error() == "登录失败:身份认证失败次数过多，您的 IP 地址已被封禁。" {
				log.Println("由于IP地址被封禁,5分钟后重试")
				time.Sleep(5 * time.Minute)
			}
			if err.Error() == "登录失败:Unauthorized" {
				log.Println("由于登录失败,3分钟后重试")
				time.Sleep(3 * time.Minute)
			}
			return false
		}
	}
	this.IsLogin = true
	return true
}

// CheckCookie 判断cookie是否可用
func (this *qbitApi) CheckCookie() bool {
	cookie := options.GetOption(options.Qbittorrent_Cookie)
	if len(cookie) <= 0 {
		return false
	}
	this.Client.SetCookie(this.parseCookie(string(cookie)))

	err, _ := this.GetApp().Version()
	if err != nil {
		return false
	}
	return true
}

// 解析Cookie
func (this *qbitApi) parseCookie(cookie string) []*http.Cookie {
	rawCookies := cookie
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\nHost: "+this.Client.GetHostHeader()+"\r\n\r\n", rawCookies)

	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))

	if err == nil {
		return req.Cookies()
	}
	return nil
}
