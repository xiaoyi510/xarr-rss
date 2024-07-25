package model

type DownloadList struct {
	Hash       string `json:"hash"`
	Title      string `json:"title"`
	CreateTime int64  `json:"create_time"`
}
