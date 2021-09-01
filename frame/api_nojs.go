// +build !js

package frame

import (
	"github.com/LIA-Aerospace/olord/gin/status"
	"github.com/gin-gonic/gin"
)

func (a *APIer) Route(g *gin.RouterGroup) {
	g.POST("post", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Allow javascript requests from localhost
		err := c.BindJSON(a)
		if err != nil {
			status.WriteBadRequest(c, "unable to bind json", err)
			return
		}
		c.JSON(200, 200)
	})
	g.GET("get", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Allow javascript requests from localhost
		c.JSON(200, a)
	})
}
