package model

type ApiReq struct {
	Method    string      `json:"method,omitempty"`
	Arguments interface{} `json:"arguments"`
	Tag       string      `json:"tag"`
}

type ApiStatusRes struct {
	Arguments struct {
		ActiveTorrentCount int `json:"activeTorrentCount"`
		CumulativeStats    struct {
			DownloadedBytes int `json:"downloadedBytes"`
			FilesAdded      int `json:"filesAdded"`
			SecondsActive   int `json:"secondsActive"`
			SessionCount    int `json:"sessionCount"`
			UploadedBytes   int `json:"uploadedBytes"`
		} `json:"cumulative-stats"`
		CurrentStats struct {
			DownloadedBytes int `json:"downloadedBytes"`
			FilesAdded      int `json:"filesAdded"`
			SecondsActive   int `json:"secondsActive"`
			SessionCount    int `json:"sessionCount"`
			UploadedBytes   int `json:"uploadedBytes"`
		} `json:"current-stats"`
		DownloadSpeed      int `json:"downloadSpeed"`
		PausedTorrentCount int `json:"pausedTorrentCount"`
		TorrentCount       int `json:"torrentCount"`
		UploadSpeed        int `json:"uploadSpeed"`
	} `json:"arguments"`
	Result string `json:"result"`
}

type ApiSessionGetVersionRes struct {
	Arguments struct {
		Version string `json:"version"`
	} `json:"arguments"`
	Result string `json:"result"`
	Tag    int    `json:"tag"`
}

type ApiTorrentGetRes struct {
	Arguments struct {
		Torrents []struct {
			ActivityDate       int     `json:"activityDate"`
			AddedDate          int     `json:"addedDate"`
			DoneDate           int     `json:"doneDate"`
			DownloadDir        string  `json:"downloadDir"`
			DownloadedEver     int     `json:"downloadedEver"`
			Error              int     `json:"error"`
			ErrorString        string  `json:"errorString"`
			HashString         string  `json:"hashString"`
			Id                 int     `json:"id"`
			LeftUntilDone      int     `json:"leftUntilDone"`
			Name               string  `json:"name"`
			PeersGettingFromUs int     `json:"peersGettingFromUs"`
			PeersSendingToUs   int     `json:"peersSendingToUs"`
			PercentDone        float64 `json:"percentDone"`
			QueuePosition      int     `json:"queuePosition"`
			RateDownload       int     `json:"rateDownload"`
			RateUpload         int     `json:"rateUpload"`
			RecheckProgress    int     `json:"recheckProgress"`
			Status             int     `json:"status"`
			TotalSize          int     `json:"totalSize"`
			UploadRatio        float64 `json:"uploadRatio"`
			UploadedEver       int     `json:"uploadedEver"`
		} `json:"torrents"`
	} `json:"arguments"`
	Result string `json:"result"`
}

type ApiTorrentRenamePathRes struct {
	Arguments struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"arguments"`
	Result string `json:"result"`
}
