package transmission

import (
	"XArr-Rss/app/sdk/transmission/model"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type Transmission struct {
	Username  string
	Password  string
	Host      string
	SessionId string
	Header    []helper.CurlHeader
}

func (this Transmission) GetClient(host, username, password string) *Transmission {
	return &Transmission{
		Host:     strings.TrimRight(host, "/"),
		Username: username,
		Password: password,
	}
}

func (this *Transmission) setHeader(name string, value string) {
	for k, v := range this.Header {
		if v.Name == name {
			v.Value = value
			this.Header[k] = v
			return
		}
	}
	this.Header = append(this.Header, helper.CurlHeader{Name: name, Value: value})
}
func (this *Transmission) Init() error {
	if this.Host == "" {
		return errors.New("请输入正确的 Transmission 地址")
	}
	//this.Header = []helper.CurlHeader{}
	// 初始化
	if len(this.Username) > 0 && len(this.Password) > 0 {
		this.setHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(this.Username+":"+this.Password)))
	}
	this.setHeader("Content-Type", "application/json")

	// 初始化一波
	rst := this.RpcPost("", 1)
	if this.SessionId == "" {
		return errors.New("请检查服务器是否正常:" + rst)
	}
	return nil
}

func (this *Transmission) TorrentGet(ids []string) (error, *model.ApiTorrentGetRes) {
	arg := make(map[string][]string)
	arg["fields"] = []string{
		"id",
		"name",
		"status",
		"hashString",
		"totalSize",
		"percentDone",
		"addedDate",
		//"trackerStats",
		"leftUntilDone",
		"rateDownload",
		"rateUpload",
		"recheckProgress",
		"peersGettingFromUs",
		"peersSendingToUs",
		"uploadRatio",
		"uploadedEver",
		"downloadedEver",
		"downloadDir",
		"error",
		"errorString",
		"doneDate",
		"queuePosition",
		"activityDate",
	}

	if len(ids) > 0 {
		arg["ids"] = ids
	}

	//{"arguments":{"torrents":[{"activityDate":0,"doneDate":0,"downloadedEver":0,"error":0,"errorString":"","id":1,"leftUntilDone":241930465,"peersGettingFromUs":0,"peersSendingToUs":0,"percentDone":0,"queuePosition":0,"rateDownload":0,"rateDownload":0,"rateUpload":0,"rateUpload":0,"status":0,"uploadRatio":-1,"uploadedEver":0}]},"result":"success"}
	//recently-active

	err, data := this.Exec("torrent-get", arg, "")
	if err != nil {
		return err, nil
	}
	res := &model.ApiTorrentGetRes{}
	err = json.Unmarshal([]byte(data), res)
	if err != nil {
		// 获取状态失败
		logsys.Error("检测状态异常:%s", "Transmission", err.Error())
		return err, nil
	}
	//log.Println("SessionGet", data)

	if res.Result == "success" {
		// 成功
		return nil, res
	} else {
		return errors.New("请求失败:" + res.Result), nil
	}
}

// {"method":"torrent-rename-path","arguments":{"ids":[1],"path":"[RPG Fudousan][12][BIG5][1080P].mp4","name":"[RPG Fudousan][12][BIG5][1080P]1.mp4"},"tag":""}:
// 修改名称 path 原始文件名 name 新的名字
func (this *Transmission) TorrentRenamePath(id int, path, name string) error {
	arg := make(map[string]interface{})
	arg["path"] = path
	arg["name"] = name
	arg["ids"] = []int{id}

	//{"arguments":{"torrents":[{"activityDate":0,"doneDate":0,"downloadedEver":0,"error":0,"errorString":"","id":1,"leftUntilDone":241930465,"peersGettingFromUs":0,"peersSendingToUs":0,"percentDone":0,"queuePosition":0,"rateDownload":0,"rateDownload":0,"rateUpload":0,"rateUpload":0,"status":0,"uploadRatio":-1,"uploadedEver":0}]},"result":"success"}
	//recently-active

	err, data := this.Exec("torrent-rename-path", arg, "")
	if err != nil {
		return err
	}
	res := &model.ApiTorrentRenamePathRes{}
	err = json.Unmarshal([]byte(data), res)
	if err != nil {
		// 获取状态失败
		logsys.Error("检测状态异常:%s", "Transmission", err.Error())
		return err
	}

	if res.Result == "success" {
		// 成功
		return nil
	} else {
		return errors.New("请求失败:" + res.Result)
	}
}

func (this *Transmission) Status() error {
	err, data := this.Exec("session-stats", "", "")
	if err != nil {
		return err
	}
	res := &model.ApiStatusRes{}
	err = json.Unmarshal([]byte(data), res)
	if err != nil {
		// 获取状态失败
		logsys.Error("检测状态异常:%s", "Transmission", err.Error())
		return err
	}
	if res.Result == "success" {
		// 成功
		return nil
	}
	log.Println("Status", data)

	return errors.New(res.Result)
}

func (this *Transmission) SessionGetVersion() (error, string) {
	arg := make(map[string][]string)
	arg["fields"] = []string{
		"version",
	}
	err, data := this.Exec("session-get", arg, "")
	if err != nil {
		return err, ""
	}
	res := &model.ApiSessionGetVersionRes{}
	err = json.Unmarshal([]byte(data), res)
	if err != nil {
		// 获取状态失败
		logsys.Error("检测状态异常:%s", "Transmission", err.Error())
		return err, ""
	}
	log.Println("SessionGet", data)

	if res.Result == "success" {
		// 成功
		return nil, res.Arguments.Version
	} else {
		return errors.New("请求失败:" + res.Result), ""
	}
}

func (this *Transmission) Exec(method string, arguments interface{}, tag string) (error, string) {
	if this.SessionId == "" {
		return errors.New("请先初始化"), ""
	}
	val := &model.ApiReq{
		Method:    method,
		Arguments: arguments,
		Tag:       tag,
	}
	valStr, _ := json.Marshal(val)
	//log.Println("请求", method, string(valStr))
	res := this.RpcPost(string(valStr), 1)
	//log.Println(res)
	return nil, res
}

func (this *Transmission) RpcPost(data string, errCount int) string {

	err, response, retHeader, statusCode := helper.CurlHelper{}.PostStringReturnHeader(this.Host+"/transmission/rpc", data, this.Header)
	if err != nil {
		return ""
	}
	//log.Println("rpc", string(response), "data", data, this.Header)
	sessionId := retHeader.Get("X-Transmission-Session-Id")
	if statusCode == 409 && len(sessionId) > 0 {
		// 初始化
		this.SessionId = sessionId

		this.setHeader("X-Transmission-Session-Id", sessionId)

		if errCount > 3 {
			return ""
		}
		return this.RpcPost(data, errCount+1)
	}

	// 返回内容
	//response
	return string(response)
}
