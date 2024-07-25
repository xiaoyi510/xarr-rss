package torrents

// ApiBaseReq 基础请求
type ApiBaseReq struct {
	Hash string `json:"hash"`
}

// ApiTorrentsSetLabelReq 设置标签
type ApiTorrentsSetLabelReq struct {
	ApiBaseReq
	Label string `json:"label"`
}

// ApiTorrentsSetCategoryReq 设置分类
type ApiTorrentsSetCategoryReq struct {
	Hashes   string `json:"hashes"`
	Category string `json:"category"`
}

// ApiTorrentsResumeAllReq 请求 恢复所有种子
type ApiTorrentsResumeAllReq struct {
	ApiBaseReq
}

// ApiTorrentsInfoReq 获取种子信息
type ApiTorrentsInfoReq struct {
	Filter   string `json:"filter,omitempty"`   // Filter torrent list by state. Allowed state filters: all, downloading, seeding, completed, paused, active, inactive, resumed, stalled, stalled_uploading, stalled_downloading, errored
	Category string `json:"category,omitempty"` // GetUri torrents with the given category (empty string means "without category"; no "category" parameter means "any category" <- broken until #11748 is resolved). Remember to Url-encode the category name. For example, My category becomes My%20category
	Tag      string `json:"tag,omitempty"`      // GetUri torrents with the given tag (empty string means "without tag"; no "tag" parameter means "any tag". Remember to Url-encode the category name. For example, My tag becomes My%20tag
	Sort     string `json:"sort,omitempty"`     // Sort torrents by given key. They can be sorted using any field of the response's JSON array (which are documented below) as the sort key.
	Reverse  bool   `json:"reverse,omitempty"`  // Enable reverse sorting. Defaults to false
	Limit    int    `json:"limit,omitempty"`    // Limit the number of torrents returned
	Offset   int    `json:"offset,omitempty"`   // Set offset (if less than 0, offset from end)
	Hashes   string `json:"hashes,omitempty"`   //Filter by hashes. Can contain multiple hashes separated by |
}

type ApiTorrentInfo struct {
	AddedOn           int     `json:"added_on"`           // 	Time (Unix Epoch) when the torrent was added to the client
	AmountLeft        int64   `json:"amount_left"`        // 	Amount of data left to download (bytes)
	AutoTmm           bool    `json:"auto_tmm"`           // 	Whether this torrent is managed by Automatic Torrent Management
	Availability      float64 `json:"availability"`       // 	Percentage of file pieces currently available
	Category          string  `json:"category"`           // 	Category of the torrent
	Completed         int64   `json:"completed"`          // 	Amount of transfer data completed (bytes)
	CompletionOn      int     `json:"completion_on"`      // 	Time (Unix Epoch) when the torrent completed
	ContentPath       string  `json:"content_path"`       // 	Absolute path of torrent content (root path for multifile torrents, absolute file path for singlefile torrents)
	DlLimit           int     `json:"dl_limit"`           // 	Torrent download speed limit (bytes/s). -1 if ulimited.
	Dlspeed           int     `json:"dlspeed"`            // 	Torrent download speed (bytes/s)
	Downloaded        int64   `json:"downloaded"`         // 	Amount of data downloaded
	DownloadedSession int64   `json:"downloaded_session"` // 	Amount of data downloaded this session
	Eta               int     `json:"eta"`                // 	Torrent ETA (seconds)
	FLPiecePrio       bool    `json:"f_l_piece_prio"`     // 	True if first last piece are prioritized
	ForceStart        bool    `json:"force_start"`        // 	True if force start is enabled for this torrent
	Hash              string  `json:"hash"`               // 	Torrent hash
	LastActivity      int     `json:"last_activity"`      // 	Last time (Unix Epoch) when a chunk was downloaded/uploaded
	MagnetUri         string  `json:"magnet_uri"`         // 	Magnet URI corresponding to this torrent
	MaxRatio          float64 `json:"max_ratio"`          // 	Maximum share ratio until torrent is stopped from seeding/uploading
	MaxSeedingTime    int     `json:"max_seeding_time"`   // 	Maximum seeding time (seconds) until torrent is stopped from seeding
	Name              string  `json:"name"`               // 	Torrent name
	NumComplete       int     `json:"num_complete"`       // 	Number of seeds in the swarm
	NumIncomplete     int     `json:"num_incomplete"`     // 	Number of leechers in the swarm
	NumLeechs         int     `json:"num_leechs"`         // 	Number of leechers connected to
	NumSeeds          int     `json:"num_seeds"`          // 	Number of seeds connected to
	Priority          int     `json:"priority"`           // 	Torrent priority. Returns -1 if queuing is disabled or torrent is in seed mode
	Progress          float64 `json:"progress"`           // 	Torrent progress (percentage/100)
	Ratio             float64 `json:"ratio"`              // 	Torrent share ratio. Max ratio value: 9999.
	RatioLimit        float64 `json:"ratio_limit"`        // 	TODO (what is different from max_ratio?)
	SavePath          string  `json:"save_path"`          // 	Path where this torrent's data is stored
	SeedingTime       int     `json:"seeding_time"`       // 	Torrent elapsed time while complete (seconds)
	SeedingTimeLimit  int     `json:"seeding_time_limit"` // 	TODO (what is different from max_seeding_time?) seeding_time_limit is a per torrent setting, when Automatic Torrent Management is disabled, furthermore then max_seeding_time is set to seeding_time_limit for this torrent. If Automatic Torrent Management is enabled, the value is -2. And if max_seeding_time is unset it have a default value -1.
	SeenComplete      int     `json:"seen_complete"`      // 	Time (Unix Epoch) when this torrent was last seen complete
	SeqDl             bool    `json:"seq_dl"`             // 	True if sequential download is enabled
	Size              int64   `json:"size"`               // 	Total size (bytes) of files selected for download
	State             string  `json:"state"`              // 	Torrent state. See table here below for the possible values
	SuperSeeding      bool    `json:"super_seeding"`      // 	True if super seeding is enabled
	Tags              string  `json:"tags"`               // 	Comma-concatenated tag list of the torrent
	TimeActive        int     `json:"time_active"`        // 	Total active time (seconds)
	TotalSize         int64   `json:"total_size"`         // 	Total size (bytes) of all file in this torrent (including unselected ones)
	Tracker           string  `json:"tracker"`            // 	The first tracker with working status. Returns empty string if no tracker is working.
	TrackersCount     int     `json:"trackers_count"`     // 	Torrent upload speed limit (bytes/s). -1 if ulimited.
	UpLimit           int     `json:"up_limit"`           // 	Amount of data uploaded
	Uploaded          int64   `json:"uploaded"`           // 	Amount of data uploaded this session
	UploadedSession   int     `json:"uploaded_session"`   // 	Torrent upload speed (bytes/s)
	Upspeed           int     `json:"upspeed"`            //
}
type ApiTorrentInfoRes []ApiTorrentInfo

type ApiTorrentPropertiesReq struct {
	Hash string `json:"hash,omitempty"` //The hash of the torrent you want to get the generic properties of
}

type ApiTorrentPropertiesRes struct {
	AdditionDate           int     `json:"addition_date"`
	Comment                string  `json:"comment"`
	CompletionDate         int     `json:"completion_date"`
	CreatedBy              string  `json:"created_by"`
	CreationDate           int     `json:"creation_date"`
	DlLimit                int     `json:"dl_limit"`
	DlSpeed                int     `json:"dl_speed"`
	DlSpeedAvg             int     `json:"dl_speed_avg"`
	Eta                    int     `json:"eta"`
	LastSeen               int     `json:"last_seen"`
	NbConnections          int     `json:"nb_connections"`
	NbConnectionsLimit     int     `json:"nb_connections_limit"`
	Peers                  int     `json:"peers"`
	PeersTotal             int     `json:"peers_total"`
	PieceSize              int     `json:"piece_size"`
	PiecesHave             int     `json:"pieces_have"`
	PiecesNum              int     `json:"pieces_num"`
	Reannounce             int     `json:"reannounce"`
	SavePath               string  `json:"save_path"`
	SeedingTime            int     `json:"seeding_time"`
	Seeds                  int     `json:"seeds"`
	SeedsTotal             int     `json:"seeds_total"`
	ShareRatio             float64 `json:"share_ratio"`
	TimeElapsed            int     `json:"time_elapsed"`
	TotalDownloaded        int64   `json:"total_downloaded"`
	TotalDownloadedSession int     `json:"total_downloaded_session"`
	TotalSize              int64   `json:"total_size"`
	TotalUploaded          int64   `json:"total_uploaded"`
	TotalUploadedSession   int     `json:"total_uploaded_session"`
	TotalWasted            int     `json:"total_wasted"`
	UpLimit                int     `json:"up_limit"`
	UpSpeed                int     `json:"up_speed"`
	UpSpeedAvg             int     `json:"up_speed_avg"`
}

type ApiTorrentTrackersReq struct {
	Hash string `json:"hash,omitempty"` //The hash of the torrent you want to get the generic properties of
}

type ApiTorrentTrackers struct {
	Msg           string      `json:"msg"`
	NumDownloaded int         `json:"num_downloaded"`
	NumLeeches    int         `json:"num_leeches"`
	NumPeers      int         `json:"num_peers"`
	NumSeeds      int         `json:"num_seeds"`
	Status        int         `json:"status"`
	Tier          interface{} `json:"tier"`
	Url           string      `json:"url"`
}

type ApiTorrentTrackersRes []ApiTorrentTrackers

//////////////////////////////
type ApiTorrentWebseedsReq struct {
	Hash string `json:"hash,omitempty"` //The hash of the torrent you want to get the generic properties of
}
type ApiTorrentWebseeds struct {
	Url string `json:"url"`
}

type ApiTorrentWebseedsRes []ApiTorrentWebseeds

//////////////////////////////
type ApiTorrentFilesReq struct {
	Hash    string `json:"hash"`              //The hash of the torrent you want to get the generic properties of
	Indexes string `json:"indexes,omitempty"` // The indexes of the files you want to retrieve. indexes can contain multiple values separated by |.
}

type ApiTorrentFiles struct {
	Availability int     `json:"availability"`
	Index        int     `json:"index"`       // File index
	IsSeed       bool    `json:"is_seed"`     // True if file is seeding/complete
	Name         string  `json:"name"`        // File name (including relative path)
	PieceRange   []int   `json:"piece_range"` // The first number is the starting piece index and the second number is the ending piece index (inclusive)
	Priority     int     `json:"priority"`    // File priority. See possible values here below
	Progress     float64 `json:"progress"`    // File progress (percentage/100)
	Size         int64   `json:"size"`        // 	File size (bytes)
}
type ApiTorrentFilesRes []ApiTorrentFiles

//////////////////////////////

type ApiTorrentPieceStatesReq struct {
	Hash string `json:"hash,omitempty"` //The hash of the torrent you want to get the generic properties of
}

//////////////////////////////
type ApiTorrentPauseReq struct {
	Hashes string `json:"hashes"` // The hashes of the torrents you want to pause. hashes can contain multiple hashes separated by |, to pause multiple torrents, or set to all, to pause all torrents.
}
type ApiTorrentResumeReq struct {
	Hashes string `json:"hashes"` // The hashes of the torrents you want to pause. hashes can contain multiple hashes separated by |, to pause multiple torrents, or set to all, to pause all torrents.
}
type ApiTorrentDeleteReq struct {
	Hashes      string `json:"hashes"`      // The hashes of the torrents you want to pause. hashes can contain multiple hashes separated by |, to pause multiple torrents, or set to all, to pause all torrents.
	DeleteFiles bool   `json:"deleteFiles"` // If set to true, the downloaded data will also be deleted, otherwise has no effect.
}
type ApiTorrentRecheckReq struct {
	Hashes string `json:"hashes"` // The hashes of the torrents you want to pause. hashes can contain multiple hashes separated by |, to pause multiple torrents, or set to all, to pause all torrents.
}
type ApiTorrentReannounceReq struct {
	Hashes string `json:"hashes"` // The hashes of the torrents you want to pause. hashes can contain multiple hashes separated by |, to pause multiple torrents, or set to all, to pause all torrents.
}

//////////////////////////////

type ApiTorrentAddReq struct {
	Urls               string  `json:"urls,omitempty"`     // URLs separated with newlines
	Savepath           string  `json:"savepath,omitempty"` //Download folder
	Cookie             string  `json:"cookie,omitempty"`
	Category           string  `json:"category,omitempty"`           // Category for the torrent
	SkipChecking       bool    `json:"skipChecking,omitempty"`       // Skip hash checking. Possible values are true, false (default)
	Paused             bool    `json:"paused,omitempty"`             // Add torrents in the paused state. Possible values are true, false (default)
	RootFolder         bool    `json:"rootFolder,omitempty"`         // CreateSource the root folder. Possible values are true, false, unset (default)
	Tags               string  `json:"tags,omitempty"`               //Tags for the torrent, split by ','
	Rename             string  `json:"rename,omitempty"`             //Rename torrent
	UpLimit            int     `json:"upLimit,omitempty"`            //	Set torrent upload speed limit. Unit in bytes/second
	DlLimit            int     `json:"dlLimit,omitempty"`            // Set torrent download speed limit. Unit in bytes/second
	RatioLimit         float32 `json:"ratioLimit,omitempty"`         //Set torrent share ratio limit
	SeedingTimeLimit   int     `json:"seedingTimeLimit,omitempty"`   //Set torrent seeding time limit. Unit in seconds
	AutoTMM            bool    `json:"autoTMM,omitempty"`            //Whether Automatic Torrent Management should be used
	SequentialDownload string  `json:"sequentialDownload,omitempty"` //Enable sequential download. Possible values are true, false (default)
	FirstLastPiecePrio string  `json:"firstLastPiecePrio,omitempty"` //Prioritize download first last piece. Possible values are true, false (default)
}

//////////////////////////////
type ApiTorrentAddTrackersReq struct {
	Hash string `json:"hash,omitempty"`
	Urls string `json:"urls,omitempty"` // 多个使用 %0A分割
}

//////////////////////////////

type ApiTorrentEditTrackersReq struct {
	Hash    string `json:"hash,omitempty"`
	OrigUrl string `json:"origUrl,omitempty"` // The tracker Url you want to edit
	NewUrl  string `json:"newUrl,omitempty"`  // The new Url to replace the origUrl
}

//////////////////////////////
type ApiTorrentRemoveTrackersReq struct {
	Hash string `json:"hash,omitempty"`
	Urls string `json:"urls,omitempty"` // 多个使用 %0A分割
}

//////////////////////////////
type ApiTorrentIncreasePrioReq struct {
	Hashes string `json:"hashes,omitempty"`
}

//////////////////////////////
type ApiTorrentDecreasePrioReq struct {
	Hashes string `json:"hashes,omitempty"`
}

//////////////////////////////
type ApiTorrentDownloadLimitReq struct {
	Hashes string `json:"hashes,omitempty"`
}
type ApiTorrentDownloadLimitRes map[string]int

//////////////////////////////
type ApiTorrentSetDownloadLimitReq struct {
	Hashes string `json:"hashes,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}
type ApiTorrentSetShareLimitsReq struct {
	Hashes           string  `json:"hashes,omitempty"`
	SeedingTimeLimit int     `json:"seedingTimeLimit,omitempty"`
	RatioLimit       float64 `json:"ratioLimit,omitempty"`
}

//////////////////////////////
type ApiTorrentUploadLimitReq struct {
	Hashes string `json:"hashes,omitempty"`
}

type ApiTorrentUploadLimitRes map[string]int

//////////////////////////////
type ApiTorrentSetUploadLimitReq struct {
	Hashes string `json:"hashes,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

//////////////////////////////
type ApiTorrentSetLocationReq struct {
	Hashes   string `json:"hashes,omitempty"`
	Location string `json:"location,omitempty"`
}

//////////////////////////////
type ApiTorrentRenameReq struct {
	Hash string `json:"hash,omitempty"`
	Name string `json:"name,omitempty"`
}
type ApiTorrentRenameFileReq struct {
	Hash    string `json:"hash,omitempty"`
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
}

type ApiTorrentCategories struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}
type ApiTorrentCategoriesRes map[string]ApiTorrentCategories

type ApiTorrentCreateCategoriesReq struct {
	Category string `json:"category"`
	SavePath string `json:"savePath"`
}

type ApiTorrentEditCategoriesReq struct {
	Category string `json:"category"`
	SavePath string `json:"savePath"`
}

type ApiTorrentRemoveCategoriesReq struct {
	Category string `json:"category"`
}

//////////////////////////////
type ApiTorrentAddTagsReq struct {
	Hashes string `json:"hashes"`
	Tags   string `json:"tags"`
}
type ApiTorrentRemoveTagsReq struct {
	Hashes string `json:"hashes"`
	Tags   string `json:"tags"`
}
type ApiTorrentTagsRes []string

//////////////////////////////

type ApiTorrentCreateTagsReq struct {
	Tags string `json:"tags"`
}

type ApiTorrentDeleteTagsReq struct {
	Tags string `json:"tags"`
}
