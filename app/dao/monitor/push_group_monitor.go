package monitor

import (
	"XArr-Rss/app/global/variable"
	"XArr-Rss/app/model"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/json"
	"os"
	"time"
)

type PushGroupMonitor struct {
}

func (this *PushGroupMonitor) SyncPushGroup() {

	if variable.ServerState.IsVip {
		logsys.Info("同步发布组监听中 30分钟同步一次 赞助会员方可自动开启", "发布组")
	} else {
		logsys.Info("同步发布组监听中 30分钟同步一次", "发布组")
	}
	for {
		//if variable.ServerState.IsVip {
		PushGroupSyncToLocal()
		//}
		time.Sleep(30 * time.Minute)
	}

}

func PushGroupSyncToLocal() {
	//https://xarr.52nyg.com/api/v1/push-group/list
	err, response := helper.CurlHttpHelper{}.GetProxyResult(variable.XArrApiHost+"/api/v1/push-group/list", "")
	if err != nil {
		if variable.XArrApiHost == "https://xarr.co1o.com" {
			variable.XArrApiHost = "https://xarr.52nyg.com"
		} else {
			variable.XArrApiHost = "https://xarr.co1o.com"
		}
		err, response = helper.CurlHttpHelper{}.GetProxyResult(variable.XArrApiHost+"/api/v1/push-group/list", "")
		if err != nil {
			logsys.Error("同步发布组获取失败:%s", "发布组", err.Error())
		}

	}
	if err == nil {
		// 解析返回数据
		resp := &model.ApiPushGroupListRes{}
		err := json.Unmarshal(response, resp)
		if err != nil {
			// 获取数据异常
			logsys.Error("同步发布组解析失败:%s", "发布组", err.Error())
		} else if resp.Code == 200 {
			// 写出文件
			marshal, _ := json.Marshal(resp.Data)
			err = os.WriteFile("./conf/push_group.json", marshal, 0777)
			if err != nil {
				logsys.Error("同步发布组写出文件失败:%s", "发布组", err.Error())
			}
		}

	}
}
