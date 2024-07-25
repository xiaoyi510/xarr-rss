package http

import (
	"XArr-Rss/app/global/variable"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

func templateLoad(r *gin.Engine, public *embed.FS) {
	// embed files
	tmpl := template.New("")
	tmpl = template.Must(tmpl.ParseFS(public, "web/*.html"))
	r.SetHTMLTemplate(tmpl)

	if variable.IsDebug {

		r.Static("/admin", "./web/admin")
		r.Static("/config", "./web/config")
		r.Static("/data", "./web/data")
		r.Static("/component", "./web/component")
		r.Static("/view", "./web/view")
	} else {

		fp, _ := fs.Sub(public, "web/admin")
		r.StaticFS("/admin", http.FS(fp))
		fp, _ = fs.Sub(public, "web/component")
		r.StaticFS("/component", http.FS(fp))
		fp, _ = fs.Sub(public, "web/data")
		r.StaticFS("/data", http.FS(fp))
		fp, _ = fs.Sub(public, "web/config")
		r.StaticFS("/config", http.FS(fp))
		fp, _ = fs.Sub(public, "web/view")
		r.StaticFS("/view", http.FS(fp))
	}

	//r.LoadHTMLGlob("./app/views/*.html")
	r.SetFuncMap(template.FuncMap{})
}
