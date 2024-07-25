package cache

import (
	"XArr-Rss/util/logsys"
	"github.com/patrickmn/go-cache"
	"time"
)

var GocacheClient *cache.Cache

const (
	CACHE_KEY_SOURCEIDS          = "source_ids"
	CACHE_KEY_GetSourcesItemList = "GetSourcesItemList"
	CACHE_KEY_GetSourceItems     = "GetSourceItems"
	CACHE_KEY_GetSeasonEpisode   = "GetSeasonEpisode"
	CACHE_KEY_GetSeasonAbEpisode = "GetSeasonAbEpisode"
)

// 初始化缓存
func init() {

	GocacheClient = cache.New(5*time.Minute, 10*time.Minute)
	GocacheClient.Set("foo", "bar", cache.DefaultExpiration)

	logsys.Info("缓存系统启动", "缓存")
}
