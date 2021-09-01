// +build !js

package main

import (
	"formtest/frame"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	grp := r.Group("/api")

	api := frame.APIer{1337, "l33t sp33k is my g4m3"}
	api.Route(grp)
	sv := http.Server{
		Addr:    ":8081",
		Handler: r,
	}
	sv.ListenAndServe()
}
