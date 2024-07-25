package api

import (
	"XArr-Rss/app/model/qbit/transfer"
	"XArr-Rss/app/sdk/qbit/client"
	"encoding/json"
	"errors"
	"strconv"
)

type Transfer struct {
	client *client.QbitClient
}

func (this *Transfer) SetClient(client *client.QbitClient) *Transfer {
	this.client = client
	return this
}

// Info
func (this *Transfer) Info() (error, *transfer.ApiTransferInfoRes) {
	req := transfer.ApiTransferInfoReq{}
	res, _ := this.client.Get("transfer/maindata", req)
	if res == "Forbidden" {
		return errors.New(res), nil
	}

	var resData transfer.ApiTransferInfoRes
	json.Unmarshal([]byte(res), &resData)

	return nil, &resData
}

// SpeedLimitsMode GetUri alternative speed limits state
func (this *Transfer) SpeedLimitsMode() (error, int) {
	req := transfer.ApiTransferSpeedLimitsModeReq{}
	res, _ := this.client.Get("transfer/speedLimitsMode", req)
	if res == "Forbidden" {
		return errors.New(res), 0
	}
	atoi, err := strconv.Atoi(res)
	if err != nil {
		return err, 0
	}

	return nil, atoi
}

// ToggleSpeedLimitsMode Toggle alternative speed limits 切换备用限速方案 int 200正确
func (this *Transfer) ToggleSpeedLimitsMode() (error, int) {
	req := transfer.ApiTransferToggleSpeedLimitsModeReq{}
	res, status := this.client.Get("transfer/toggleSpeedLimitsMode", req)
	if res == "Forbidden" {
		return errors.New(res), status
	}
	return nil, status
}

func (this *Transfer) DownloadLimit() (error, int) {
	req := transfer.ApiTransferDownloadLimitReq{}
	res, _ := this.client.Get("transfer/downloadLimit", req)
	if res == "Forbidden" {
		return errors.New(res), 0
	}
	atoi, err := strconv.Atoi(res)
	if err != nil {
		return err, 0
	}

	return nil, atoi
}

// SetDownloadLimit 设置下载限速
func (this *Transfer) SetDownloadLimit() (error, int) {
	req := transfer.ApiTransferSetDownloadLimitReq{}
	res, status := this.client.Post("transfer/setDownloadLimit", req)
	if res == "Forbidden" {
		return errors.New(res), status
	}
	return nil, status
}

func (this *Transfer) UploadLimit() (error, int) {
	req := transfer.ApiTransferUploadLimitReq{}
	res, _ := this.client.Get("transfer/uploadLimit", req)
	if res == "Forbidden" {
		return errors.New(res), 0
	}
	atoi, err := strconv.Atoi(res)
	if err != nil {
		return err, 0
	}
	return nil, atoi
}

func (this *Transfer) SetUploadLimit() (error, int) {
	req := transfer.ApiTransferSetUploadLimitReq{}
	res, _ := this.client.Get("transfer/setUploadLimit", req)
	if res == "Forbidden" {
		return errors.New(res), 0
	}
	atoi, err := strconv.Atoi(res)
	if err != nil {
		return err, 0
	}
	return nil, atoi
}

func (this *Transfer) BanPeers() (error, int) {
	req := transfer.ApiTransferSetUploadLimitReq{}
	res, status := this.client.Get("transfer/banPeers", req)
	if res == "Forbidden" {
		return errors.New(res), status
	}
	return nil, status
}
