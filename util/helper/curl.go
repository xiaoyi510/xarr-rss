package helper

import (
	"XArr-Rss/util/logsys"
	"crypto/tls"
	"errors"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CurlHelper struct {
	Client    *http.Client
	HttpProxy func(*http.Request) (*url.URL, error)
	Header    []CurlHeader
	Transport *http.Transport
}
type CurlHeader struct {
	Name  string
	Value string
}

func (this *CurlHelper) SetHttpProxy(proxyUri string) *CurlHelper {
	if proxyUri == "" {
		return this
	}

	ProxyURL, err := url.Parse(proxyUri)
	if err != nil {
		logsys.Error("代理格式错误:%s", "代理", proxyUri)
		return this
	} else {
		this.Transport.Proxy = http.ProxyURL(ProxyURL)
	}

	return this
}

func (this *CurlHelper) SetProxy(proxyUri string) *CurlHelper {
	if proxyUri == "" {
		return this
	}
	if strings.Contains(proxyUri, "http") {
		return this.SetHttpProxy(proxyUri)
	}
	ProxyURL, err := url.Parse(proxyUri)
	if err != nil {
		logsys.Error("代理格式错误:%s", "代理", proxyUri)
		return this
	}
	dialer, err := proxy.FromURL(ProxyURL, proxy.Direct)
	//dialer, err := proxy.SOCKS5("tcp", proxyUri, nil, proxy.Direct)
	if err != nil {
		logsys.Error("socks5代理错误:%s", "代理", err.Error())
		return this
	}

	this.Transport.Dial = dialer.Dial

	return this
}

func (this *CurlHelper) Init(headers []CurlHeader) *CurlHelper {
	// 取消SSL证书校验
	this.Transport = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		ExpectContinueTimeout: 60 * time.Second,
	}

	if headers != nil && len(headers) > 0 {
		this.Header = headers
	} else {
		this.Header = append(this.Header, CurlHeader{
			Name:  "User-Agent",
			Value: "Mozilla/5.0 (Macintosh Go; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		})
	}
	this.Client = &http.Client{Transport: this.Transport, Timeout: 60 * time.Second}

	return this
}
func (this *CurlHelper) Get(url string, data *url.Values, checkResponseStatus bool) (error, []byte, int) {
	client := this.Client
	if data != nil {
		if strings.Contains(url, "?") {
			url += "&" + data.Encode()
		} else {
			url += "?" + data.Encode()
		}
	}
	logsys.Debug("HTTP GET Uri:%s", "HTTP", url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, nil, 0
	}
	if this.Header != nil && len(this.Header) > 0 {
		for _, v := range this.Header {
			request.Header.Set(v.Name, v.Value)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return err, nil, 0
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()

	result, _ := ioutil.ReadAll(response.Body)
	if checkResponseStatus {
		if response.StatusCode != 200 {
			return errors.New("返回状态码异常:" + response.Status), nil, response.StatusCode
		}
	}
	client.CloseIdleConnections()

	return nil, result, response.StatusCode
}

func (this CurlHelper) GetUri(url string, data *url.Values, headers []CurlHeader, checkResponseStatus bool) (error, []byte, int) {

	// 取消SSL证书校验
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   50,
	}

	client := &http.Client{Transport: tr, Timeout: 60 * time.Second}

	if data != nil {
		if strings.Contains(url, "?") {
			url += "&" + data.Encode()
		} else {
			url += "?" + data.Encode()
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, nil, 0
	}
	if headers != nil {
		for _, v := range headers {
			request.Header.Set(v.Name, v.Value)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return err, nil, 0
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()

	result, _ := ioutil.ReadAll(response.Body)
	if checkResponseStatus {
		if response.StatusCode != 200 {
			return errors.New("返回状态码异常:" + response.Status), nil, response.StatusCode
		}
	}
	client.CloseIdleConnections()

	return nil, result, response.StatusCode
}

func (this CurlHelper) Post(url string, data *url.Values, headers []CurlHeader) (error, []byte) {

	// 取消SSL证书校验
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	var body *strings.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err, nil
	}
	if len(headers) == 0 {
		headers = append(headers, CurlHeader{
			Name:  "Content-Type",
			Value: "application/x-www-form-urlencoded; charset=UTF-8",
		})
	}

	for _, v := range headers {
		request.Header.Set(v.Name, v.Value)
	}
	response, err := client.Do(request)
	if err != nil {
		return err, nil
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	result, _ := ioutil.ReadAll(response.Body)
	return nil, result
}

func (this CurlHelper) PostReturnHeader(url string, data *url.Values, headers []CurlHeader) (error, []byte, http.Header, int) {

	// 取消SSL证书校验
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	var body *strings.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err, nil, nil, 0
	}
	if len(headers) == 0 {
		headers = append(headers, CurlHeader{
			Name:  "Content-Type",
			Value: "application/x-www-form-urlencoded; charset=UTF-8",
		})
	}

	for _, v := range headers {
		request.Header.Set(v.Name, v.Value)
	}
	response, err := client.Do(request)
	if err != nil {
		return err, nil, nil, 0
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	result, _ := ioutil.ReadAll(response.Body)
	return nil, result, response.Header, response.StatusCode
}

func (this CurlHelper) PostStringReturnHeader(url string, data string, headers []CurlHeader) (error, []byte, http.Header, int) {

	// 取消SSL证书校验
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	var body *strings.Reader

	body = strings.NewReader(data)

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err, nil, nil, 0
	}
	if len(headers) == 0 {
		headers = append(headers, CurlHeader{
			Name:  "Content-Type",
			Value: "application/x-www-form-urlencoded; charset=UTF-8",
		})
	}

	for _, v := range headers {
		request.Header.Set(v.Name, v.Value)
	}
	response, err := client.Do(request)
	if err != nil {
		return err, nil, nil, 0
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	result, _ := ioutil.ReadAll(response.Body)
	return nil, result, response.Header, response.StatusCode
}

func (this CurlHelper) PostString(url string, data string, headers []CurlHeader) (error, []byte) {
	// 取消SSL证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err, nil
	}
	for _, v := range headers {
		request.Header.Set(v.Name, v.Value)
	}
	response, err := client.Do(request)
	if err != nil {
		return err, nil
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	result, _ := ioutil.ReadAll(response.Body)
	return nil, result
}
