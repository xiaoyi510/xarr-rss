package transfer

type ApiTransferInfoReq struct {
}

type ApiTransferInfoRes struct {
	ConnectionStatus string `json:"connection_status"` //    Connection status.See possible values here below
	DhtNodes         int    `json:"dht_nodes"`         //    DHT nodes connected to
	DlInfoData       int64  `json:"dl_info_data"`      //    Data downloaded this session (bytes)
	DlInfoSpeed      int    `json:"dl_info_speed"`     //    Global download rate (bytes/s)
	DlRateLimit      int    `json:"dl_rate_limit"`     //    Download rate limit (bytes/s)
	UpInfoData       int64  `json:"up_info_data"`      //    Data uploaded this session (bytes)
	UpInfoSpeed      int    `json:"up_info_speed"`     //    Global upload rate (bytes/s)
	UpRateLimit      int    `json:"up_rate_limit"`     //    Upload rate limit (bytes/s)
}

type ApiTransferSpeedLimitsModeReq struct {
}
type ApiTransferToggleSpeedLimitsModeReq struct {
}

type ApiTransferDownloadLimitReq struct {
}

type ApiTransferUploadLimitReq struct {
}

type ApiTransferSetDownloadLimitReq struct {
	Limit int `json:"limit"` //The global download speed limit to set in bytes/second
}

type ApiTransferSetUploadLimitReq struct {
	Limit int `json:"limit"` //The global upload speed limit to set in bytes/second
}

type ApiTransferBanPeersReq struct {
	Peers int `json:"peers"` //The peer to ban, or multiple peers separated by a pipe |. Each peer is a colon-separated host:port
}
