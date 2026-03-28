package api

import (
	"github.com/gin-gonic/gin"

	"atlas/capture"
)

func SetupRouter(scanner *capture.Scanner) *gin.Engine {
	r := gin.Default()

	r.StaticFile("/", "./static/index.html")

	r.GET("/graph", getGraph(scanner))             // fetch current graph in memory
	r.POST("/scanDuration", scanDuration(scanner)) // scan for fixed duration

	r.POST("/start", startScan(scanner))   // start scan without timeout
	r.POST("/stop", stopScan(scanner))     // stop any running scan
	r.GET("/stream", streamGraph(scanner)) // SSE connection to stream graph

	return r
}
