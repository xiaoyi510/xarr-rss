package http

import (
	"XArr-Rss/app/global/conf/appconf"
	"XArr-Rss/util/logsys"
	"embed"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

// Run 启动一个http服 作为总服务端监控通知
func Run(public *embed.FS) {
	// 记录到文件。
	f, _ := os.Create("http.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.Use(Cors())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// 推送服务
	routers(r)
	templateLoad(r, public)

	port := "8086"
	port = strconv.Itoa(appconf.AppConf.System.Port)
	logsys.Info("监听地址: 0.0.0.0:"+port, "系统")
	go func() {
		log.Println("pprof", http.ListenAndServe(":8087", nil))

	}()
	_ = r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
