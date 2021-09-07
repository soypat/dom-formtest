// +build !js

package main

import (
	"net/http"

	"github.com/soypat/dom-formtest/js/frame"

	"github.com/LIA-Aerospace/ctl/platforms/vtvleatt"
	"github.com/gin-gonic/gin"
)

func main() {
	params := vtvleatt.DefaultParameters
	r := gin.Default()
	grp := r.Group("/api")

	api, err := frame.New("", "", &params)
	must(err)

	api.Route(grp)
	sv := http.Server{
		Addr:    ":8085",
		Handler: r,
	}
	sv.ListenAndServe()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func Route(a *APIer, g *gin.RouterGroup) {
	g.POST("post", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Allow javascript requests from localhost
		err := c.BindJSON(a.dataPtr)
		if err != nil {
			status.WriteBadRequest(c, "unable to bind json", err)
			return
		}
		c.JSON(200, 200)
	})
	g.GET("get", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Allow javascript requests from localhost
		c.JSON(200, a.dataPtr)
	})
}
