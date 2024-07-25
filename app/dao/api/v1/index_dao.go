package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexDao struct {
}

func (this IndexDao) Index(c *gin.Context) {

	err := checkAuth(c)
	if err != nil {
		c.Redirect(302, "/login")
		return
	}
	//c.Redirect(http.StatusOK, "/")
	c.HTML(http.StatusOK, "index.html", gin.H{})

}

func (this IndexDao) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
