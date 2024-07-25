package sync

import "XArr-Rss/app/model/qbit/torrents"

type ApiSyncMaindataReq struct {
	Rid int `json:"rid,omitempty"` //Response Id. If not provided, rid=0 will be assumed. If the given rid is different from the one of last server reply, full_update will be true (see the server reply details for more info)
}

type ApiSyncMaindataRes struct {
	Rid               int                                 `json:"rid,omitempty"`                //	Response Id
	FullUpdate        bool                                `json:"full_update,omitempty"`        //	Whether the response contains all the data or partial data
	Torrents          map[string]ApiSyncMaindataTorrents  `json:"torrents,omitempty"`           //	Property: torrent hash, value: same as torrent list
	TorrentsRemoved   []string                            `json:"torrents_removed,omitempty"`   //	List of hashes of torrents removed since last request
	Categories        map[string]ApiSyncMaindataCategorie `json:"categories,omitempty"`         //	Info for categories added since last request
	CategoriesRemoved map[string]ApiSyncMaindataCategorie `json:"categories_removed,omitempty"` //	List of categories removed since last request
	Tags              interface{}                         `json:"tags,omitempty"`               //	List of tags added since last request
	TagsRemoved       map[int]string                      `json:"tags_removed,omitempty"`       //	List of tags removed since last request
	ServerState       ApiSyncMaindataServerState          `json:"server_state,omitempty"`       //	Global transfer info
	Trackers          map[string]ApiSyncMainataTrackers   `json:"trackers"`
}

type ApiSyncMainataTrackers []string

type ApiSyncMaindataTorrents struct {
	AddedOn           int     `json:"added_on"`
	AmountLeft        int     `json:"amount_left"`
	AutoTmm           bool    `json:"auto_tmm"`
	Availability      float32 `json:"availability"`
	Category          string  `json:"category"`
	Completed         int     `json:"completed"`
	CompletionOn      int     `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	DlLimit           int     `json:"dl_limit"`
	Dlspeed           int     `json:"dlspeed"`
	Downloaded        int     `json:"downloaded"`
	DownloadedSession int     `json:"downloaded_session"`
	Eta               int     `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	LastActivity      int     `json:"last_activity"`
	MagnetUri         string  `json:"magnet_uri"`
	MaxRatio          float32 `json:"max_ratio"`
	MaxSeedingTime    int     `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int     `json:"num_complete"`
	NumIncomplete     int     `json:"num_incomplete"`
	NumLeechs         int     `json:"num_leechs"`
	NumSeeds          int     `json:"num_seeds"`
	Priority          int     `json:"priority"`
	Progress          float64 `json:"progress"`
	Ratio             float64 `json:"ratio"`
	RatioLimit        float32 `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SeedingTime       int     `json:"seeding_time"`
	SeedingTimeLimit  int     `json:"seeding_time_limit"`
	SeenComplete      int     `json:"seen_complete"`
	SeqDl             bool    `json:"seq_dl"`
	Size              int     `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int     `json:"time_active"`
	TotalSize         int     `json:"total_size"`
	Tracker           string  `json:"tracker"`
	TrackersCount     int     `json:"trackers_count"`
	UpLimit           int     `json:"up_limit"`
	Uploaded          int     `json:"uploaded"`
	UploadedSession   int     `json:"uploaded_session"`
	Upspeed           int     `json:"upspeed"`
}

type ApiSyncMaindataServerState struct {
	AlltimeDl            int64  `json:"alltime_dl"`
	AlltimeUl            int64  `json:"alltime_ul"`
	AverageTimeQueue     int    `json:"average_time_queue"`
	ConnectionStatus     string `json:"connection_status"`
	DhtNodes             int    `json:"dht_nodes"`
	DlInfoData           int64  `json:"dl_info_data"`
	DlInfoSpeed          int    `json:"dl_info_speed"`
	DlRateLimit          int    `json:"dl_rate_limit"`
	FreeSpaceOnDisk      int64  `json:"free_space_on_disk"`
	GlobalRatio          string `json:"global_ratio"`
	QueuedIoJobs         int    `json:"queued_io_jobs"`
	Queueing             bool   `json:"queueing"`
	ReadCacheHits        string `json:"read_cache_hits"`
	ReadCacheOverload    string `json:"read_cache_overload"`
	RefreshInterval      int    `json:"refresh_interval"`
	TotalBuffersSize     int    `json:"total_buffers_size"`
	TotalPeerConnections int    `json:"total_peer_connections"`
	TotalQueuedSize      int    `json:"total_queued_size"`
	TotalWastedSession   int    `json:"total_wasted_session"`
	UpInfoData           int    `json:"up_info_data"`
	UpInfoSpeed          int    `json:"up_info_speed"`
	UpRateLimit          int    `json:"up_rate_limit"`
	UseAltSpeedLimits    bool   `json:"use_alt_speed_limits"`
	WriteCacheOverload   string `json:"write_cache_overload"`
}

type ApiSyncMaindataCategorie struct {
	Name     string `json:"name,omitempty"`
	SavePath string `json:"save_path,omitempty"`
}

////////////////////////

type ApiSyncTorrentPeersReq struct {
	torrents.ApiBaseReq
	Rid string `json:"rid"` // Response Id. If not provided, rid=0 will be assumed. If the given rid is different from the one of last server reply, full_update will be true (see the server reply details for more info)
}
type ApiSyncTorrentPeersRes struct {
	FullUpdate bool                           `json:"full_update"`
	Peers      map[string]ApiSyncTorrentPeers `json:"peers"`
	Rid        int                            `json:"rid"`
	ShowFlags  bool                           `json:"show_flags"`
}

type ApiSyncTorrentPeers struct {
	Client      string `json:"client"`
	Connection  string `json:"connection"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	DlSpeed     int    `json:"dl_speed"`
	Downloaded  int    `json:"downloaded"`
	Files       string `json:"files"`
	Flags       string `json:"flags"`
	FlagsDesc   string `json:"flags_desc"`
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	Progress    int    `json:"progress"`
	Relevance   int    `json:"relevance"`
	UpSpeed     int    `json:"up_speed"`
	Uploaded    int    `json:"uploaded"`
}

////////////////////////
