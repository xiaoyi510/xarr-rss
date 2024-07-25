package helper

import (
	"XArr-Rss/util/logsys"
	"crypto/tls"
	"errors"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CurlHttpHelper struct {
	Client *http.Client

	Transport *http.Transport
	Header    *http.Header

	Response struct {
		Header http.Header
		Status int
	}
}

// 获取客户端
func GetCurlHttpHelper(headers *http.Header) *CurlHttpHelper {
	return (&CurlHttpHelper{}).Init(headers)
}

// 获取默认HTTP客户端
func GetCurlHttpHelperDefault() *CurlHttpHelper {
	headers := &http.Header{}
	headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7;  XArr) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	return (&CurlHttpHelper{}).Init(headers)
}

// 初始化
func (this *CurlHttpHelper) Init(headers *http.Header) *CurlHttpHelper {
	// 取消SSL证书校验
	this.Transport = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   60 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		ExpectContinueTimeout: 60 * time.Second,
		IdleConnTimeout:       60 * time.Second,
	}

	this.Header = headers
	if this.Header == nil {
		this.Header = &http.Header{}
	}
	this.Client = &http.Client{Transport: this.Transport, Timeout: 60 * time.Second}

	return this
}

// 设置代理
func (this *CurlHttpHelper) SetHttpProxy(proxyUri string) (*CurlHttpHelper, error) {
	if proxyUri == "" {
		return this, nil
	}

	ProxyURL, err := url.Parse(proxyUri)
	if err != nil {
		return this, logsys.Error("代理格式错误:%s", "代理", proxyUri)
	} else {
		this.Transport.Proxy = http.ProxyURL(ProxyURL)
	}

	return this, nil
}

// 设置代理
func (this *CurlHttpHelper) SetProxy(proxyUri string) (*CurlHttpHelper, error) {
	if proxyUri == "" {
		return this, nil
	}
	// 判断是否为IP地址
	address := net.ParseIP(proxyUri)
	if address != nil {
		return this.SetHttpProxy("http://" + proxyUri)
	}
	if strings.Contains(proxyUri, "http") {
		return this.SetHttpProxy(proxyUri)
	}

	ProxyURL, err := url.Parse(proxyUri)
	if err != nil {

		return this, logsys.Error("代理格式错误:%s", "代理", proxyUri)
	}
	dialer, err := proxy.FromURL(ProxyURL, proxy.Direct)
	//dialer, err := proxy.SOCKS5("tcp", proxyUri, nil, proxy.Direct)
	if err != nil {
		return this, logsys.Error("socks5代理错误:%s", "代理", err.Error())
	}

	this.Transport.Dial = dialer.Dial

	return this, nil
}

// 清理请求信息
func (this *CurlHttpHelper) clear() {
	this.Response.Header = http.Header{}
	this.Response.Status = 0
}

// 设置默认Header头
func (this *CurlHttpHelper) setDefaultHeader() {
	if this.Header.Get("Content-Type") == "" {
		this.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
}

// Get 请求
func (this *CurlHttpHelper) Get(url string, data *url.Values, checkResponseStatus bool) (error, []byte) {
	this.clear()
	if data != nil {
		if strings.Contains(url, "?") {
			url += "&" + data.Encode()
		} else {
			url += "?" + data.Encode()
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, nil
	}
	if this.Header != nil {
		request.Header = *this.Header
	}

	response, err := this.Client.Do(request)
	if err != nil {
		return err, nil
	}
	defer response.Body.Close()
	defer this.Client.CloseIdleConnections()

	result, _ := ioutil.ReadAll(response.Body)
	if checkResponseStatus {
		if response.StatusCode != 200 {
			return errors.New("返回状态码异常:" + response.Status), nil
		}
	}
	this.Client.CloseIdleConnections()

	this.Response.Header = response.Header
	this.Response.Status = response.StatusCode
	return nil, result
}

// Post 请求
func (this *CurlHttpHelper) Post(url string, data *url.Values, checkResponseStatus bool) (error, []byte) {
	this.clear()

	var body *strings.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err, nil
	}
	this.setDefaultHeader()
	request.Header = *this.Header

	response, err := this.Client.Do(request)
	if err != nil {
		return err, nil
	}
	defer response.Body.Close()
	defer this.Client.CloseIdleConnections()

	result, _ := ioutil.ReadAll(response.Body)
	if checkResponseStatus {
		if response.StatusCode != 200 {
			return errors.New("返回状态码异常:" + response.Status), nil
		}
	}
	this.Client.CloseIdleConnections()

	this.Response.Header = response.Header
	this.Response.Status = response.StatusCode
	return nil, result
}

// Post 请求
func (this *CurlHttpHelper) PostString(url string, data string, checkResponseStatus bool) (error, []byte) {
	this.clear()

	var body *strings.Reader
	body = strings.NewReader(data)

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err, nil
	}
	this.setDefaultHeader()
	request.Header = *this.Header

	response, err := this.Client.Do(request)
	if err != nil {
		return err, nil
	}
	defer response.Body.Close()
	defer this.Client.CloseIdleConnections()

	result, _ := ioutil.ReadAll(response.Body)
	if checkResponseStatus {
		if response.StatusCode != 200 {
			return errors.New("返回状态码异常:" + response.Status), nil
		}
	}
	this.Client.CloseIdleConnections()

	this.Response.Header = response.Header
	this.Response.Status = response.StatusCode
	return nil, result
}

// 根据Proxy 去GET请求
func (this CurlHttpHelper) GetProxyResult(uri, proxy string) (error, []byte) {
	httpHelper := GetCurlHttpHelper(nil)
	_, err := httpHelper.SetProxy(proxy)
	if err != nil {
		return err, nil
	}
	return httpHelper.Get(uri, nil, false)
}
