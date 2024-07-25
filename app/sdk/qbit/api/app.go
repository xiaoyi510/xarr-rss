package api

import (
	"XArr-Rss/app/model/qbit/app"
	"XArr-Rss/app/sdk/qbit/client"
	"encoding/json"
	"errors"
)

type App struct {
	client *client.QbitClient
}

func (this *App) SetClient(client *client.QbitClient) *App {
	this.client = client
	return this
}

// Preferences 首选项 GetUri application preferences
func (this *App) Preferences() (error, *app.ApiAppPreferencesRes) {
	// 调用API进行登录
	req := app.ApiAppPreferencesReq{}
	res, _ := this.client.Get("app/preferences", req)
	if res == "Forbidden" {
		return errors.New(res), nil
	}

	var resData app.ApiAppPreferencesRes
	json.Unmarshal([]byte(res), &resData)

	return nil, &resData
}

// Version GetUri application version
func (this *App) Version() (error, string) {
	// 调用API进行登录
	req := app.ApiAppVersionReq{}
	req.Test = "123"
	res, _ := this.client.Get("app/version", req)
	if res == "Forbidden" {
		return errors.New(res), ""
	}
	return nil, res
}

// WebApiVersion GetUri API version
func (this *App) WebApiVersion() (error, string) {
	// 调用API进行登录
	req := app.ApiAppWebApiVersionReq{}
	res, _ := this.client.Get("app/webapiVersion", req)
	if res == "Forbidden" {
		return errors.New(res), ""
	}
	return nil, res
}

// BuildInfo GetUri build info
func (this *App) BuildInfo() (error, *app.ApiAppBuildInfoRes) {
	// 调用API进行登录
	req := app.ApiAppBuildInfoReq{}
	res, _ := this.client.Get("app/buildInfo", req)
	if res == "Forbidden" {
		return errors.New(res), nil
	}
	var resData app.ApiAppBuildInfoRes

	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}

// Shutdown application
func (this *App) Shutdown() (error, string) {
	// 调用API进行登录
	req := app.ApiAppShutdownReq{}
	res, _ := this.client.Get("app/shutdown", req)
	if res == "Forbidden" {
		return errors.New(res), ""
	}
	return nil, res
}

// DefaultSavePath 获取默认保存文职
func (this *App) DefaultSavePath() (error, string) {
	// 调用API进行登录
	req := app.ApiAppShutdownReq{}
	res, _ := this.client.Get("app/defaultSavePath", req)
	if res == "Forbidden" {
		return errors.New(res), ""
	}
	return nil, res
}
