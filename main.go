package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// create scanner
	scanner := NewScanner("en0")
	// create router
	r := gin.Default()

	// endpoints
	r.GET("/graph", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, scanner.packageGraph.ToJSON())
	})

	r.POST("/scan", func(ctx *gin.Context) {
		// parse query args
		durationArg := ctx.Query("duration")
		durationMS, err := strconv.Atoi(durationArg)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid argument",
			})
		}

		scanner.Scan(durationMS, false)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "scan completed",
		})
	})

	// browser interface
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// run the server
	r.Run(":8080")
}
