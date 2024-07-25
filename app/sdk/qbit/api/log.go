package api

import (
	"XArr-Rss/app/model/qbit/app"
	"XArr-Rss/app/sdk/qbit/client"
	"errors"
)

type Log struct {
	client *client.QbitClient
}

func (this *Log) SetClient(client *client.QbitClient) *Log {
	this.client = client
	return this
}

func (this *Log) GetLog() (error, string) {
	// 调用API进行登录
	req := app.ApiAppShutdownReq{}
	res, _ := this.client.Get("log/main", req)
	if res == "Forbidden" {
		return errors.New(res), ""
	}
	return nil, res
}
