package http

import (
	v12 "XArr-Rss/app/dao/api/v1"
	"XArr-Rss/app/global/variable"
	"github.com/gin-gonic/gin"
)

func routers(r *gin.Engine) {
	apiV1Group(r)
	r.Any("/", v12.IndexDao{}.Index)
	r.Any("/login", v12.IndexDao{}.Login)
	r.Any("/login.html", func(context *gin.Context) {
		// 跳转到 Oauth中心
		context.Redirect(302, variable.XArrApiHost+"/login/oauth")
	})

	r.Any("/rss/group/group_:groupId.xml", v12.RssDao{}.Result)
	r.Any("/torznab/:groupId", v12.ApiTorznabNewDao{}.Api)
	r.Any("/torznab/:groupId/api", v12.ApiTorznabNewDao{}.Api)

}

func apiV1Group(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// 系统设置
		v1.GET("/system/get", v12.Apiv1SystemDao{}.Get)
		v1.GET("/system/words-rule", v12.Apiv1SystemDao{}.GetWordsRule)
		v1.POST("/system/edit-words-rule", v12.Apiv1SystemDao{}.EditWordsRule)
		v1.POST("/system/save", v12.Apiv1SystemDao{}.Save)
		v1.GET("/system/export", v12.Apiv1SystemDao{}.Export)
		//
		v1.GET("/service/get", v12.Apiv1ServiceDao{}.Get)
		v1.POST("/service/save", v12.Apiv1ServiceDao{}.Save)
		v1.GET("/service/sonarr-reload", v12.Apiv1ServiceDao{}.SonarrReload)
		v1.GET("/service/sonarr-reload-real", v12.Apiv1ServiceDao{}.SonarrReloadReal)
		v1.GET("/service/sonarr-test", v12.Apiv1ServiceDao{}.SonarrTest)
		v1.GET("/service/tmdb-test", v12.Apiv1ServiceDao{}.TmdbTest)
		v1.GET("/service/tmdb-reload", v12.Apiv1ServiceDao{}.TmdbReload)
		v1.GET("/service/tmdb-reload-real", v12.Apiv1ServiceDao{}.TmdbReloadRefresh)
		v1.GET("/service/qbit-test", v12.Apiv1ServiceDao{}.QbitTest)
		v1.GET("/service/transmision-test", v12.Apiv1ServiceDao{}.TransmissionTest)

		v1.GET("/medias/get", v12.ApiV1Medias{}.Get)
		v1.GET("/medias/info", v12.ApiV1Medias{}.Info)
		v1.POST("/medias/rename", v12.ApiV1Medias{}.Rename)
		v1.GET("/menu/get", v12.ApiV1Menu{}.Get)

		// 获取所有分组 可选groupId
		v1.GET("/group_push", v12.ApiV1Groups{}.GetPushGroup)
		v1.GET("/groups/get/*groupId", v12.ApiV1Groups{}.Get) // 已做
		//v1.GET("/groups/get/*groupId/*mediaId", v12.ApiV1Groups{}.Get) // 已做
		v1.POST("/groups/add", v12.ApiV1Groups{}.Add)                  // 已做
		v1.DELETE("/groups/delete/:groupId", v12.ApiV1Groups{}.Delete) // 已做
		v1.POST("/groups/edit/:groupId", v12.ApiV1Groups{}.Edit)       // 已做
		v1.GET("/groups/refresh", v12.ApiV1Groups{}.Refresh)
		v1.POST("/groups/set-sonarr-index", v12.ApiV1Groups{}.SetSonarrIndex)

		// 分组模板列表
		v1.GET("/groups/template/list", v12.Apiv1GroupTemplateDao{}.GetList)
		v1.POST("/groups/template/add", v12.Apiv1GroupTemplateDao{}.Add)
		v1.POST("/groups/template/edit", v12.Apiv1GroupTemplateDao{}.Edit)
		v1.POST("/groups/template/delete", v12.Apiv1GroupTemplateDao{}.Delete)
		v1.GET("/groups/template/info", v12.Apiv1GroupTemplateDao{}.Info)
		v1.POST("/groups/template/batch_use", v12.Apiv1GroupTemplateDao{}.BatchUse)

		// 获取分组下面使用的所有媒体源
		v1.GET("/groups/medias/:groupId", v12.ApiV1Groups{}.GetGroupMedias)
		v1.GET("/groups/medias/:groupId/:mediaId", v12.Apiv1GroupMediasDao{}.GetMediasInfo)
		v1.POST("/groups/medias/:groupId/add", v12.Apiv1GroupMediasDao{}.AddMedias)
		v1.POST("/groups/medias/:groupId/edit/:mediaId", v12.Apiv1GroupMediasDao{}.EditMedias)
		v1.DELETE("/groups/medias/:groupId/delete/:mediaId", v12.Apiv1GroupMediasDao{}.DeleteMedias)
		v1.DELETE("/groups/medias/:groupId/batchRemove", v12.Apiv1GroupMediasDao{}.BatchRemove)

		v1.GET("/sources/get", v12.ApiV1Sources{}.Get)                                 // 已做
		v1.GET("/sources/get/*sourceId", v12.ApiV1Sources{}.Get)                       // 已做
		v1.GET("/sources/:sourceId/get-medias", v12.ApiV1Sources{}.GetMedias)          // 获取数据源下面的媒体信息
		v1.GET("/sources/:sourceId/get-medias-json", v12.ApiV1Sources{}.GetMediasJson) // 获取数据源下面的媒体信息
		//v1.GET("/sources/get-site", v12.ApiV1Sources{}.GetSite)               // 已做
		v1.POST("/sources/add", v12.ApiV1Sources{}.Add)                             // 已做
		v1.DELETE("/sources/delete/:sourceId", v12.ApiV1Sources{}.Delete)           // 已做
		v1.POST("/sources/edit/:sourceId", v12.ApiV1Sources{}.Edit)                 // 已做
		v1.GET("/sources/refresh/:sourceId", v12.ApiV1Sources{}.Refresh)            // 已做
		v1.GET("/sources/refresh-parse/:sourceId", v12.ApiV1Sources{}.RefreshParse) // 已做
		v1.GET("/sonarr/medias", v12.Apiv1SonarrDao{}.Medias)                       // 已做
		v1.GET("/sonarr/image/:sonarrId/:image", v12.Apiv1SonarrDao{}.Image)        // 已做

		v1.POST("/group-medias/auto-gen-reg", v12.Apiv1GroupMediasDao{}.AutoGenReg)
		v1.POST("/medias/test-reg", v12.ApiV1Medias{}.TestReg)
		v1.Any("/down", v12.Apiv1DownDao{}.Down)
		v1.Any("/login", v12.Apiv1LoginDao{}.Login)
		v1.GET("/user/info", v12.Apiv1UserDao{}.Info)
		v1.GET("/static/info", v12.Apiv1StaticsDao{}.Info)
		v1.GET("/one-said", v12.ApiV1OneSaidDao{}.Get)
	}

}
