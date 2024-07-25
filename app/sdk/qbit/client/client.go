package client

import (
	"XArr-Rss/app/model/qbit"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"strings"
	"time"
)

type QbitClient struct {
	conf       qbit.ModelQbit
	httpClient *http.Client
	jar        http.CookieJar
	cookies    []*http.Cookie
	IsLogin    bool
}

// Init 初始化
func (this *QbitClient) Init(serverUrl string) *QbitClient {
	this.conf.ServerUrl = serverUrl

	this.jar, _ = cookiejar.New(&cookiejar.Options{})

	// 取消SSL证书校验
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	this.httpClient = &http.Client{
		Jar:       this.jar,
		Transport: tr,
		Timeout:   30 * time.Second,
	}
	return this
}

func (this *QbitClient) SetConf(serverUrl string) {
	this.conf.ServerUrl = serverUrl

}

func (this *QbitClient) GetScheme() string {
	if strings.Contains(this.conf.ServerUrl, "https") {
		return "https"
	}
	return "http"
}

func (this *QbitClient) GetHost() string {
	sprintf := fmt.Sprintf("%s/api/v2",
		strings.TrimRight(this.conf.ServerUrl, "/"),
	)
	return sprintf
}

func (this *QbitClient) GetHostHeader() string {
	host := strings.ReplaceAll(this.conf.ServerUrl, "https://", "")
	host = strings.ReplaceAll(this.conf.ServerUrl, "http://", "")
	sprintf := fmt.Sprintf("%s:%s",
		strings.TrimRight(host, "/"),
	)
	return sprintf
}

func (this *QbitClient) Get(url string, data interface{}) (string, int) {
	tmpDataMap, _ := json.Marshal(data)

	dataMap := make(map[string]string)
	err := json.Unmarshal(tmpDataMap, &dataMap)
	if err != nil {
		return err.Error(), 0
	}

	reqData := url2.Values{}
	for k, v := range dataMap {
		reqData.Set(k, v)
	}

	url = this.GetHost() + "/" + url
	uri, err := url2.ParseRequestURI(url)
	if err != nil {
		return "", 0
	}
	uri.RawQuery = reqData.Encode()

	//fmt.Println("[GET] Url:" + url)
	request, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return err.Error(), 0
	}

	response, err := this.httpClient.Do(request)
	if err != nil {
		if response == nil {
			return err.Error(), -1
		}
		return err.Error(), response.StatusCode
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error(), response.StatusCode
	}

	return string(body), response.StatusCode
}

func (this *QbitClient) Post(url string, data interface{}) (string, int) {
	tmpDataMap, _ := json.Marshal(data)

	dataMap := make(map[string]interface{})
	err := json.Unmarshal(tmpDataMap, &dataMap)
	if err != nil {
		return err.Error(), 0
	}

	url = this.GetHost() + "/" + url
	//fmt.Println("[POST] Url:" + url)

	///////////////////////////

	form := url2.Values{}
	for k, v := range dataMap {
		form.Add(k, fmt.Sprintf("%v", v))
	}
	//fmt.Println(form.Encode())
	///////////////////////////
	request, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return err.Error(), 0
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := this.httpClient.Do(request)
	if err != nil {
		if response == nil {
			return err.Error(), -1
		}
		return err.Error(), response.StatusCode
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error(), response.StatusCode
	}
	this.cookies = response.Cookies()
	return string(body), response.StatusCode
}

func (this *QbitClient) GenHashs(hashs []string) string {
	return strings.Join(hashs, "|")
}

// SetCookie 设置Cookie
func (this *QbitClient) SetCookie(cookie []*http.Cookie) *QbitClient {
	cookieUrl := &url2.URL{
		Scheme:      this.GetScheme(),
		Opaque:      "",
		User:        nil,
		Host:        this.GetHostHeader(),
		Path:        "",
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}
	this.httpClient.Jar.SetCookies(cookieUrl, cookie)
	return this
}

func (this *QbitClient) GetCookie() string {
	var ret string
	for _, v := range this.cookies {
		ret += v.Name + "=" + v.Value + "; "
	}
	return ret
}
