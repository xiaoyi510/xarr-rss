package variable

import (
	"XArr-Rss/util/logsys"
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ServerVersion   = "2.0.0.7"
	IsDebug         = false
	ConfDir         = ""
	ProgramFileTime = ""
	HttpTokenAfter  = ""
	XArrApiHost     = "https://xarr.52nyg.com"
	LoginToken      = sync.Map{}

	ServerState serverState
	MainDb      *sql.DB

	SonarrTags = make(map[int]string)
)

// 服务器授权信息
type serverState struct {
	IsVip       bool
	ExpireTime  int64
	Username    string
	Email       string
	LicenseKey  string
	RefreshTime int64 // 间隔时间不能超出5s 超出拉黑
}

func CheckTokenExpires() {
	logsys.Info("开启监控登录信息", "系统")

	for {
		time.Sleep(1 * time.Second)
		LoginToken.Range(func(key, value interface{}) bool {
			LoginToken.Store(key, value.(int)-1)

			if value.(int)-1 < 0 {
				LoginToken.Delete(key)
				logsys.Error("Cookie过期:%v", "系统", key)
			}
			return true
		})
	}
}

func CheckToken(token string) bool {
	if token == "" {
		return false
	}
	has := false
	LoginToken.Range(func(key, value interface{}) bool {
		if token == key {
			has = true
			return false
		}
		return true
	})
	return has
}
