package dbmodel

type DownloadList struct {
	Model
	Hash          string `json:"hash"`
	Title         string `json:"title"`
	Status        int    `json:"status"`
	Process       int    `json:"process"`
	OriginalTitle string `json:"original_title"`
	DownClient    string `json:"down_client"`
}

const (
	DownloadListStatusWait        = iota + 1 // 1. 等待
	DownloadListStatusRename                 // 2. 已修改名称
	DownloadListStatusDownSuccess            // 3. 已下载完成
)
