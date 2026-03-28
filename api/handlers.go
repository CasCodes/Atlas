package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"atlas/capture"
)

// returns the current state of graph in memory
func getGraph(s *capture.Scanner) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, s.GetGraphJSON())
	}
}

func startScan(s *capture.Scanner) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// parse query args
		durationArg := ctx.Query("duration")
		durationMS, err := strconv.Atoi(durationArg)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid argument",
			})
		}

		s.Scan(durationMS, false)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "scan completed",
		})
	}
}
