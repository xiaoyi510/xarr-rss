package v1

import (
	"XArr-Rss/app/model/apiv1/group"
	"XArr-Rss/app/model/dbmodel"
	"XArr-Rss/app/service/groups"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Apiv1GroupTemplateDao struct {
}

// 获取模板列表
func (this Apiv1GroupTemplateDao) GetList(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	groupTemplateList := groups.GroupTemplateService{}.GetList()
	c.JSON(200, gin.H{
		"code":    0,
		"data":    groupTemplateList,
		"message": "success",
	})
}

func (this Apiv1GroupTemplateDao) Info(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	id := cast.ToInt32(c.Query("id"))
	gid := cast.ToInt(c.Query("gid"))
	if gid > 0 {
		// 查询groupInfo
		groupInfo, _ := groups.GroupsService{}.GroupInfo(gid)
		if groupInfo != nil && groupInfo.GroupTemplateId > 0 {
			id = groupInfo.GroupTemplateId
		}
	}

	if id <= 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}

	info := groups.GroupTemplateService{}.FindOne(id)
	if info == nil {

		c.JSON(200, gin.H{
			"code":    404,
			"data":    nil,
			"message": "没有找到信息",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    0,
		"data":    info,
		"message": "success",
	})
}

func (d Apiv1GroupTemplateDao) Add(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	var req group.Apiv1GroupTemplateAdd
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}

	if req.Name == "" {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请输入模板名称",
		})
		return
	}

	item := &dbmodel.GroupTemplate{
		Name:            req.Name,
		Language:        req.Language,
		Quality:         req.Quality,
		Regex:           req.Regex,
		UseSource:       req.UseSource,
		FilterPushGroup: req.FilterPushGroup,
		EchoTitleAnime:  req.EchoTitleAnime,
		EchoTitleTv:     req.EchoTitleTv,
	}

	err = groups.GroupTemplateService{}.Add(item)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "添加失败:" + err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "添加成功",
	})
}

func (this Apiv1GroupTemplateDao) Edit(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}

	var req group.Apiv1GroupTemplateEdit
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}

	if req.Id < 1 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请选择要编辑的数据",
		})
		return
	}

	info := groups.GroupTemplateService{}.FindOne(req.Id)
	if info == nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "未找到数据",
		})
		return
	}

	info.Name = req.Name
	info.Language = req.Language
	info.Quality = req.Quality
	info.Regex = req.Regex
	info.UseSource = req.UseSource
	info.FilterPushGroup = req.FilterPushGroup
	info.EchoTitleAnime = req.EchoTitleAnime
	info.EchoTitleTv = req.EchoTitleTv

	err = groups.GroupTemplateService{}.Edit(info)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "编辑失败:" + err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "编辑成功",
	})
}

func (this Apiv1GroupTemplateDao) Delete(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	var req group.Apiv1GroupTemplateDelete
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}
	if req.Id <= 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求参数错误",
		})
		return
	}

	delErr := groups.GroupTemplateService{}.Delete(req.Id)
	if delErr != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "删除失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "删除成功",
	})
}

func (this Apiv1GroupTemplateDao) BatchUse(c *gin.Context) {
	err := checkAuth(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    401,
			"data":    nil,
			"message": "登录异常:" + err.Error(),
		})
		return
	}
	var req group.Apiv1GroupTemplateBatchUse
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请求数据异常",
		})
		return
	}
	if req.Id <= 0 || len(req.GroupMediaIds) == 0 {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "请选择媒体或者您请求参数有问题",
		})
		return
	}

	delErr := groups.GroupTemplateService{}.BatchUse(req.Id, req.GroupMediaIds)
	if delErr != nil {
		c.JSON(200, gin.H{
			"code":    502,
			"data":    nil,
			"message": "批量设置失败:" + delErr.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"data":    nil,
		"message": "批量设置成功",
	})
}
