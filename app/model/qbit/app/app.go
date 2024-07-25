package app

// ApiAppPreferencesReq 首选项请求
type ApiAppPreferencesReq struct {
}

// ApiAppPreferencesRes 首选项返回
type ApiAppPreferencesRes struct {
	locale                                 string         //	Currently selected language (e.g. en_GB for English)
	create_subfolder_enabled               bool           //	True if a subfolder should be created when adding a torrent
	start_paused_enabled                   bool           //	True if torrents should be added in a Paused state
	auto_delete_mode                       int            //	TODO
	preallocate_all                        bool           //	True if disk space should be pre-allocated for all files
	incomplete_files_ext                   bool           //	True if ".!qB" should be appended to incomplete files
	auto_tmm_enabled                       bool           //	True if Automatic Torrent Management is enabled by default
	torrent_changed_tmm_enabled            bool           //	True if torrent should be relocated when its Category changes
	save_path_changed_tmm_enabled          bool           //	True if torrent should be relocated when the default save path changes
	category_changed_tmm_enabled           bool           //	True if torrent should be relocated when its Category's save path changes
	save_path                              string         //	Default save path for torrents, separated by slashes
	temp_path_enabled                      bool           //	True if folder for incomplete torrents is enabled
	temp_path                              string         //	Path for incomplete torrents, separated by slashes
	scan_dirs                              map[string]int //	Property: directory to watch for torrent files, value: where torrents loaded from this directory should be downloaded to (see list of possible values below). Slashes are used as path separators; multiple key/value pairs can be specified
	export_dir                             string         //	Path to directory to copy .torrent files to. Slashes are used as path separators
	export_dir_fin                         string         //	Path to directory to copy .torrent files of completed downloads to. Slashes are used as path separators
	mail_notification_enabled              bool           //	True if e-mail notification should be enabled
	mail_notification_sender               string         //	e-mail where notifications should originate from
	mail_notification_email                string         //	e-mail to send notifications to
	mail_notification_smtp                 string         //	smtp server for e-mail notifications
	mail_notification_ssl_enabled          bool           //	True if smtp server requires SSL connection
	mail_notification_auth_enabled         bool           //	True if smtp server requires authentication
	mail_notification_username             string         //	Username for smtp authentication
	mail_notification_password             string         //	Password for smtp authentication
	autorun_enabled                        bool           //	True if external program should be run after torrent has finished downloading
	autorun_program                        string         //	Program path/name/arguments to run if autorun_enabled is enabled; path is separated by slashes; you can use %f and %n arguments, which will be expanded by qBittorent as path_to_torrent_file and torrent_name (from the GUI; not the .torrent file name) respectively
	queueing_enabled                       bool           //	True if torrent queuing is enabled
	max_active_downloads                   int            //	Maximum number of active simultaneous downloads
	max_active_torrents                    int            //	Maximum number of active simultaneous downloads and uploads
	max_active_uploads                     int            //	Maximum number of active simultaneous uploads
	dont_count_slow_torrents               bool           //	If true torrents w/o any activity (stalled ones) will not be counted towards max_active_* limits; see dont_count_slow_torrents for more information
	slow_torrent_dl_rate_threshold         int            //	Download rate in KiB/s for a torrent to be considered "slow"
	slow_torrent_ul_rate_threshold         int            //	Upload rate in KiB/s for a torrent to be considered "slow"
	slow_torrent_inactive_timer            int            //	Seconds a torrent should be inactive before considered "slow"
	max_ratio_enabled                      bool           //	True if share ratio limit is enabled
	max_ratio                              float32        //	GetUri the global share ratio limit
	max_ratio_act                          int            //	Action performed when a torrent reaches the maximum share ratio. See list of possible values here below.
	listen_port                            int            //	Port for incoming connections
	upnp                                   bool           //	True if UPnP/NAT-PMP is enabled
	random_port                            bool           //	True if the port is randomly selected
	dl_limit                               int            //	Global download speed limit in KiB/s; -1 means no limit is applied
	up_limit                               int            //	Global upload speed limit in KiB/s; -1 means no limit is applied
	max_connec                             int            //	Maximum global number of simultaneous connections
	max_connec_per_torrent                 int            //	Maximum number of simultaneous connections per torrent
	max_uploads                            int            //	Maximum number of upload slots
	max_uploads_per_torrent                int            //	Maximum number of upload slots per torrent
	stop_tracker_timeout                   int            //	Timeout in seconds for a stopped announce request to trackers
	enable_piece_extent_affinity           bool           //	True if the advanced libtorrent option piece_extent_affinity is enabled
	bittorrent_protocol                    int            //	Bittorrent Protocol to use (see list of possible values below)
	limit_utp_rate                         bool           //	True if [du]l_limit should be applied to uTP connections; this option is only available in qBittorent built against libtorrent version 0.16.X and higher
	limit_tcp_overhead                     bool           //	True if [du]l_limit should be applied to estimated TCP overhead (service data: e.g. packet headers)
	limit_lan_peers                        bool           //	True if [du]l_limit should be applied to peers on the LAN
	alt_dl_limit                           int            //	Alternative global download speed limit in KiB/s
	alt_up_limit                           int            //	Alternative global upload speed limit in KiB/s
	scheduler_enabled                      bool           //	True if alternative limits should be applied according to schedule
	schedule_from_hour                     int            //	Scheduler starting hour
	schedule_from_min                      int            //	Scheduler starting minute
	schedule_to_hour                       int            //	Scheduler ending hour
	schedule_to_min                        int            //	Scheduler ending minute
	scheduler_days                         int            //	Scheduler days. See possible values here below
	dht                                    bool           //	True if DHT is enabled
	pex                                    bool           //	True if PeX is enabled
	lsd                                    bool           //	True if LSD is enabled
	encryption                             int            //	See list of possible values here below
	anonymous_mode                         bool           //	If true anonymous mode will be enabled; read more here; this option is only available in qBittorent built against libtorrent version 0.16.X and higher
	proxy_type                             int            //	See list of possible values here below
	proxy_ip                               string         //	HttpProxy IP address or domain name
	proxy_port                             int            //	HttpProxy port
	proxy_peer_connections                 bool           //	True if peer and web seed connections should be proxified; this option will have any effect only in qBittorent built against libtorrent version 0.16.X and higher
	proxy_auth_enabled                     bool           //	True proxy requires authentication; doesn't apply to SOCKS4 proxies
	proxy_username                         string         //	Username for proxy authentication
	proxy_password                         string         //	Password for proxy authentication
	proxy_torrents_only                    bool           //	True if proxy is only used for torrents
	ip_filter_enabled                      bool           //	True if external IP filter should be enabled
	ip_filter_path                         string         //	Path to IP filter file (.dat, .p2p, .p2b files are supported); path is separated by slashes
	ip_filter_trackers                     bool           //	True if IP filters are applied to trackers
	web_ui_domain_list                     string         //	Comma-separated list of domains to accept when performing Host header validation
	web_ui_address                         string         //	IP address to use for the WebUI
	web_ui_port                            int            //	WebUI port
	web_ui_upnp                            bool           //	True if UPnP is used for the WebUI port
	web_ui_username                        string         //	WebUI username
	web_ui_password                        string         //	For API ≥ v2.3.0: Plaintext WebUI password, not readable, write-only. For API < v2.3.0: MD5 hash of WebUI password, hash is generated from the following string: username:Web UI Access:plain_text_web_ui_password
	web_ui_csrf_protection_enabled         bool           //	True if WebUI CSRF protection is enabled
	web_ui_clickjacking_protection_enabled bool           //	True if WebUI clickjacking protection is enabled
	web_ui_secure_cookie_enabled           bool           //	True if WebUI cookie Secure flag is enabled
	web_ui_max_auth_fail_count             int            //	Maximum number of authentication failures before WebUI access ban
	web_ui_ban_duration                    int            //	WebUI access ban duration in seconds
	web_ui_session_timeout                 int            //	Seconds until WebUI is automatically signed off
	web_ui_host_header_validation_enabled  bool           //	True if WebUI host header validation is enabled
	bypass_local_auth                      bool           //	True if authentication challenge for loopback address (127.0.0.1) should be disabled
	bypass_auth_subnet_whitelist_enabled   bool           //	True if webui authentication should be bypassed for clients whose ip resides within (at least) one of the subnets on the whitelist
	bypass_auth_subnet_whitelist           string         //	(White)list of ipv4/ipv6 subnets for which webui authentication should be bypassed; list entries are separated by commas
	alternative_webui_enabled              bool           //	True if an alternative WebUI should be used
	alternative_webui_path                 string         //	File path to the alternative WebUI
	use_https                              bool           //	True if WebUI HTTPS access is enabled
	ssl_key                                string         //	For API < v2.0.1: SSL keyfile contents (this is a not a path)
	ssl_cert                               string         //	For API < v2.0.1: SSL certificate contents (this is a not a path)
	web_ui_https_key_path                  string         //	For API ≥ v2.0.1: Path to SSL keyfile
	web_ui_https_cert_path                 string         //	For API ≥ v2.0.1: Path to SSL certificate
	dyndns_enabled                         bool           //	True if server DNS should be updated dynamically
	dyndns_service                         int            //	See list of possible values here below
	dyndns_username                        string         //	Username for DDNS service
	dyndns_password                        string         //	Password for DDNS service
	dyndns_domain                          string         //	Your DDNS domain name
	rss_refresh_interval                   int            //	RSS refresh interval
	rss_max_articles_per_feed              int            //	Max stored articles per RSS feed
	rss_processing_enabled                 bool           //	Enable processing of RSS feeds
	rss_auto_downloading_enabled           bool           //	Enable auto-downloading of torrents from the RSS feeds
	rss_download_repack_proper_episodes    bool           //	For API ≥ v2.5.1: Enable downloading of repack/proper Episodes
	rss_smart_episode_filters              string         //	For API ≥ v2.5.1: List of RSS Smart Episode Filters
	add_trackers_enabled                   bool           //	Enable automatic adding of trackers to new torrents
	add_trackers                           string         //	List of trackers to add to new torrent
	web_ui_use_custom_http_headers_enabled bool           //	For API ≥ v2.5.1: Enable custom http headers
	web_ui_custom_http_headers             string         //	For API ≥ v2.5.1: List of custom http headers
	max_seeding_time_enabled               bool           //	True enables max seeding time
	max_seeding_time                       int            //	Number of minutes to seed a torrent
	announce_ip                            string         //	TODO
	announce_to_all_tiers                  bool           //	True always announce to all tiers
	announce_to_all_trackers               bool           //	True always announce to all trackers in a tier
	async_io_threads                       int            //	Number of asynchronous I/O threads
	banned_IPs                             string         //	List of banned IPs
	checking_memory_use                    int            //	Outstanding memory when checking torrents in MiB
	current_interface_address              string         //	IP Address to bind to. Empty String means All addresses
	current_network_interface              string         //	Network Interface used
	disk_cache                             int            //	Disk cache used in MiB
	disk_cache_ttl                         int            //	Disk cache expiry interval in seconds
	embedded_tracker_port                  int            //	Port used for embedded tracker
	enable_coalesce_read_write             bool           //	True enables coalesce reads & writes
	enable_embedded_tracker                bool           //	True enables embedded tracker
	enable_multi_connections_from_same_ip  bool           //	True allows multiple connections from the same IP address
	enable_os_cache                        bool           //	True enables os cache
	enable_upload_suggestions              bool           //	True enables sending of upload piece suggestions
	file_pool_size                         int            //	File pool size
	outgoing_ports_max                     int            //	Maximal outgoing port (0: Disabled)
	outgoing_ports_min                     int            //	Minimal outgoing port (0: Disabled)
	recheck_completed_torrents             bool           //	True rechecks torrents on completion
	resolve_peer_countries                 bool           //	True resolves peer countries
	save_resume_data_interval              int            //	Save resume data interval in min
	send_buffer_low_watermark              int            //	Send buffer low watermark in KiB
	send_buffer_watermark                  int            //	Send buffer watermark in KiB
	send_buffer_watermark_factor           int            //	Send buffer watermark factor in percent
	socket_backlog_size                    int            //	Socket backlog size
	upload_choking_algorithm               int            //	Upload choking algorithm used (see list of possible values below)
	upload_slots_behavior                  int            //	Upload slots behavior used (see list of possible values below)
	upnp_lease_duration                    int            //	UPnP lease duration (0: Permanent lease)
	utp_tcp_mixed_mode                     int            //	μTP-TCP mixed mode algorithm (see list of possible values below)
}

type ApiAppVersionReq struct {
	Test string `json:"test,omitempty"`
}

type ApiAppWebApiVersionReq struct {
}

//GetUri build info
type ApiAppBuildInfoReq struct {
}

type ApiAppBuildInfoRes struct {
	qt         string //QT version
	libtorrent string //libtorrent version
	boost      string //Boost //version
	openssl    string //OpenSSL version
	bitness    int    //Application bitness (e.g. 64-bit)
}

type ApiAppShutdownReq struct {
}
