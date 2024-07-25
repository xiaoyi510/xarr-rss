package api

import (
	"XArr-Rss/app/model/qbit/torrents"
	"XArr-Rss/app/sdk/qbit/client"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Torrents struct {
	client *client.QbitClient
}

func (t *Torrents) SetClient(client *client.QbitClient) *Torrents {
	t.client = client
	return t
}

func (t *Torrents) SetLabel(hashes []string, label string) {
	//command/setLabel
	// 调用API进行登录
	req := torrents.ApiTorrentsSetLabelReq{
		Label: label,
	}
	req.Hash = t.client.GenHashs(hashes)

	res, _ := t.client.Post("command/setLabel", req)
	fmt.Println("设置结果", res)
}

// ResumeAll 恢复所有
func (t *Torrents) ResumeAll(hashes []string) {
	//command/setLabel
	// 调用API进行登录
	req := torrents.ApiTorrentsResumeAllReq{}
	req.Hash = t.client.GenHashs(hashes)

	res, _ := t.client.Post("command/resumeAll", req)
	fmt.Println("设置结果", res)
}

// Info 获取种子列表
func (t *Torrents) Info(hashes []string, req torrents.ApiTorrentsInfoReq) (error, *torrents.ApiTorrentInfoRes) {
	//req := torrents.ApiTorrentsInfoReq{}
	req.Hashes = t.client.GenHashs(hashes)
	res, status := t.client.Post("torrents/info", req)
	if status != 200 {
		return errors.New(res), nil
	}
	var resData torrents.ApiTorrentInfoRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}

// Properties 获取torrent泛型属性
func (t *Torrents) Properties(hash string) (error, *torrents.ApiTorrentPropertiesRes) {
	req := torrents.ApiTorrentPropertiesReq{}
	req.Hash = hash
	res, status := t.client.Post("torrents/properties", req)
	if status != 200 {
		return errors.New(res), nil
	}
	var resData torrents.ApiTorrentPropertiesRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}

// Trackers 获取torrent trackers
func (t *Torrents) Trackers(hash string) (error, *torrents.ApiTorrentTrackersRes) {
	req := torrents.ApiTorrentPropertiesReq{}
	req.Hash = hash
	res, status := t.client.Post("torrents/trackers", req)
	if status != 200 {
		return errors.New(res), nil
	}
	var resData torrents.ApiTorrentTrackersRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}

func (t *Torrents) Webseeds(hash string) (error, *torrents.ApiTorrentWebseedsRes) {
	req := torrents.ApiTorrentWebseedsReq{}
	req.Hash = hash
	res, status := t.client.Post("torrents/webseeds", req)
	if status != 200 {
		return errors.New(res), nil
	}
	var resData torrents.ApiTorrentWebseedsRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}

// Files GetUri torrent contents
func (t *Torrents) Files(req torrents.ApiTorrentFilesReq) (error, *torrents.ApiTorrentFilesRes) {
	res, status := t.client.Post("torrents/files", req)
	if status != 200 {
		return errors.New(res), nil
	}
	var resData torrents.ApiTorrentFilesRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}

func (t *Torrents) Pause(hashes []string) (error, int) {
	req := torrents.ApiTorrentPauseReq{}
	req.Hashes = t.client.GenHashs(hashes)
	res, status := t.client.Post("torrents/pause", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Resume(hashes []string) (error, int) {
	req := torrents.ApiTorrentResumeReq{}
	req.Hashes = t.client.GenHashs(hashes)
	res, status := t.client.Post("torrents/resume", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Delete(hashes []string, deleteFiles bool) (error, int) {
	req := torrents.ApiTorrentDeleteReq{}
	req.Hashes = t.client.GenHashs(hashes)
	req.DeleteFiles = deleteFiles
	res, status := t.client.Post("torrents/delete", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Recheck(hashes []string) (error, int) {
	req := torrents.ApiTorrentRecheckReq{}
	req.Hashes = t.client.GenHashs(hashes)
	res, status := t.client.Post("torrents/recheck", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Reannounce(hashes []string) (error, int) {
	req := torrents.ApiTorrentRecheckReq{}
	req.Hashes = t.client.GenHashs(hashes)
	res, status := t.client.Post("torrents/reannounce", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Add(req torrents.ApiTorrentAddReq) (error, int) {
	res, status := t.client.Post("torrents/add", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}
func (t *Torrents) AddTrackers(req torrents.ApiTorrentAddTrackersReq) (error, int) {
	res, status := t.client.Post("torrents/addTrackers", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) EditTracker(req torrents.ApiTorrentEditTrackersReq) (error, int) {
	res, status := t.client.Post("torrents/editTracker", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status

	/*
		Returns:

		HTTP Status Code	Scenario
		400	newUrl is not a valid Url
		404	Torrent hash was not found
		409	newUrl already exists for the torrent
		409	origUrl was not found
		200	All other scenarios
	*/
}

func (t *Torrents) RemoveTrackers(req torrents.ApiTorrentRemoveTrackersReq) (error, int) {
	res, status := t.client.Post("torrents/removeTrackers", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status

	/*
			Returns:
			HTTP Status Code	Scenario
		404	Torrent hash was not found
		409	All urls were not found
		200	All other scenarios
	*/
}

// IncreasePrio 提高种子优先级
func (t *Torrents) IncreasePrio(req torrents.ApiTorrentIncreasePrioReq) (error, int) {
	res, status := t.client.Post("torrents/increasePrio", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
	/*
			Returns:
			HTTP Status Code	Scenario
		409	Torrent queueing is not enabled
		200	All other scenarios
	*/
}

// DecreasePrio 降低种子优先级
func (t *Torrents) DecreasePrio(req torrents.ApiTorrentDecreasePrioReq) (error, int) {
	res, status := t.client.Post("torrents/decreasePrio", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
	/*
			Returns:
			HTTP Status Code	Scenario
		409	Torrent queueing is not enabled
		200	All other scenarios
	*/
}

// DownloadLimit 获取种子下载限制
func (t *Torrents) DownloadLimit(req torrents.ApiTorrentDownloadLimitReq) (error, *torrents.ApiTorrentDownloadLimitRes) {
	res, status := t.client.Post("torrents/downloadLimit", req)
	if status != 200 {
		return errors.New(res), nil
	}
	var resData torrents.ApiTorrentDownloadLimitRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}

	return nil, &resData
}

// SetDownloadLimit 设置种子下载限制
func (t *Torrents) SetDownloadLimit(req torrents.ApiTorrentSetDownloadLimitReq) (error, int) {
	res, status := t.client.Post("torrents/setDownloadLimit", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

// SetShareLimits  设置种子分享限制
func (t *Torrents) SetShareLimits(hashes []string, req torrents.ApiTorrentSetShareLimitsReq) (error, int) {
	req.Hashes = t.client.GenHashs(hashes)
	res, status := t.client.Post("torrents/setShareLimits", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) UploadLimit(req torrents.ApiTorrentUploadLimitReq) (error, *torrents.ApiTorrentUploadLimitRes) {
	res, status := t.client.Post("torrents/uploadLimit", req)
	if status != 200 {
		return errors.New(res), nil
	}

	var resData torrents.ApiTorrentUploadLimitRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}
	return nil, &resData
}

func (t *Torrents) SetUploadLimit(req torrents.ApiTorrentSetUploadLimitReq) (error, int) {
	res, status := t.client.Post("torrents/setUploadLimit", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) SetLocation(req torrents.ApiTorrentSetLocationReq) (error, int) {
	res, status := t.client.Post("torrents/setLocation", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Rename(req torrents.ApiTorrentRenameReq) (error, int) {
	res, status := t.client.Post("torrents/rename", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) RenameFile(req torrents.ApiTorrentRenameFileReq) (error, int) {
	res, status := t.client.Post("torrents/renameFile", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) SetCategory(hashes []string, category string) (error, int) {
	req := torrents.ApiTorrentsSetCategoryReq{
		Category: category,
	}
	req.Hashes = t.client.GenHashs(hashes)

	res, status := t.client.Post("torrents/setCategory", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Categories() (error, *torrents.ApiTorrentCategoriesRes) {
	var req interface{}
	res, status := t.client.Post("torrents/categories", req)
	if status != 200 {
		return errors.New(res), nil
	}

	var resData torrents.ApiTorrentCategoriesRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}
	return nil, &resData
}

func (t *Torrents) CreateCategory(categoryName string, savePath string) (error, int) {
	req := torrents.ApiTorrentCreateCategoriesReq{
		Category: categoryName,
		SavePath: savePath,
	}

	res, status := t.client.Post("torrents/createCategory", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}
func (t *Torrents) EditCategory(categoryName string, savePath string) (error, int) {
	req := torrents.ApiTorrentEditCategoriesReq{
		Category: categoryName,
		SavePath: savePath,
	}

	res, status := t.client.Post("torrents/editCategory", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) RemoveCategory(categoryName []string) (error, int) {
	req := torrents.ApiTorrentRemoveCategoriesReq{
		Category: strings.Join(categoryName, "\n"),
	}

	res, status := t.client.Post("torrents/removeCategories", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) AddTags(hashes []string, tag []string) (error, int) {
	req := torrents.ApiTorrentAddTagsReq{
		Tags:   strings.Join(tag, ","),
		Hashes: t.client.GenHashs(hashes),
	}

	res, status := t.client.Post("torrents/addTags", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) RemoveTags(hashes []string, tag []string) (error, int) {
	req := torrents.ApiTorrentRemoveTagsReq{
		Tags:   strings.Join(tag, ","),
		Hashes: t.client.GenHashs(hashes),
	}

	res, status := t.client.Post("torrents/removeTags", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) Tags() (error, *torrents.ApiTorrentTagsRes) {
	var req interface{}
	res, status := t.client.Post("torrents/tags", req)
	if status != 200 {
		return errors.New(res), nil
	}

	var resData torrents.ApiTorrentTagsRes
	err := json.Unmarshal([]byte(res), &resData)
	if err != nil {
		return err, nil
	}
	return nil, &resData
}

func (t *Torrents) CreateTags(req torrents.ApiTorrentCreateTagsReq) (error, int) {
	res, status := t.client.Post("torrents/createTags", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}

func (t *Torrents) DeleteTags(req torrents.ApiTorrentDeleteTagsReq) (error, int) {
	res, status := t.client.Post("torrents/deleteTags", req)
	if status != 200 {
		return errors.New(res), status
	}
	return nil, status
}
