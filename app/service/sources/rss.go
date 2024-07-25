package sources

import (
	"XArr-Rss/app/model"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// 获取rss数据
func GetRssUrlResult(url, proxy, name, autoSearch, searchTitle string) (*model.RssRoot, error) {
	curlHelper := &helper.CurlHelper{}
	queryUri := url
	if autoSearch != "" {
		queryUri = GetAutoSearchUrl(url, autoSearch, searchTitle)
	}
	err, data, statusCode := curlHelper.Init(nil).SetProxy(proxy).Get(queryUri, nil, true)
	if err != nil {
		return nil, logsys.Error("同步数据源[%s]失败:%s", "数据源", name, err.Error())
	}
	if statusCode != 200 {
		return nil, logsys.Error("同步数据源[%s]失败:返回错误码%s", "数据源", name, statusCode)
	}

	var ret *model.RssRoot
	if strings.Contains(url, "://nyaa.si/") {
		// nyaa 支持
		ret, err = NyaaRssSource{}.Parse(string(data))
	} else if autoSearch != "" {
		// 自动检索
		ret, err = AutoSearchParse(url, queryUri, autoSearch, string(data))
	} else {
		// 常规解析
		ret, err = GeneralSource{}.Parse(string(data))
	}

	if err != nil {
		return nil, errors.New(err.Error())
	}

	if ret == nil {
		return nil, errors.New("获取数据错误")
	}

	return ret, err
}

func AutoSearchParse(url, queryUri, autoSearch string, data string) (*model.RssRoot, error) {
	switch autoSearch {
	case CAutoSearchMikan:
		// 蜜柑自动搜索解析
		return AutoSearchMikan{}.Parse(url, queryUri, data)

	default:
		return nil, errors.New("暂不支持")

	}
}

func GetAutoSearchUrl(baseUri, autoSearch, searchTitle string) string {
	switch autoSearch {
	case CAutoSearchMikan:
		if searchTitle == "" {
			return baseUri + "Home/Classic"
		}
		return baseUri + fmt.Sprintf("Home/Search?searchstr=%s", url.QueryEscape(searchTitle))
	default:
		return baseUri
	}
}
