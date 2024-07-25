package api

import (
	"XArr-Rss/app/model/qbit/sync"
	"XArr-Rss/app/sdk/qbit/client"
	"encoding/json"
	"errors"
	"log"
)

type Sync struct {
	client *client.QbitClient
}

func (this *Sync) SetClient(client *client.QbitClient) *Sync {
	this.client = client
	return this
}

// Maindata GetUri main data
func (this *Sync) Maindata() (error, *sync.ApiSyncMaindataRes) {
	// 调用API进行登录
	req := sync.ApiSyncMaindataReq{}
	res, _ := this.client.Get("sync/maindata", req)
	if res == "Forbidden" {
		return errors.New(res), nil
	}

	var resData sync.ApiSyncMaindataRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		log.Println(err.Error() + res)
		return err, nil
	}

	return nil, &resData
}

// TorrentPeers GetUri torrent peers data
func (this *Sync) TorrentPeers() (error, *sync.ApiSyncTorrentPeersRes) {
	// 调用API进行登录
	req := sync.ApiSyncTorrentPeersReq{}
	res, _ := this.client.Get("sync/torrentPeers", req)
	if res == "Forbidden" {
		return errors.New(res), nil
	}

	var resData sync.ApiSyncTorrentPeersRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}
