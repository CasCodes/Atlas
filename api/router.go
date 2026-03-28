package api

import (
	"github.com/gin-gonic/gin"

	"atlas/capture"
)

func SetupRouter(scanner *capture.Scanner) *gin.Engine {
	r := gin.Default()

	r.StaticFile("/", "./static/index.html")

	r.GET("/graph", getGraph(scanner))
	r.POST("/scan", startScan(scanner))

	return r
}
