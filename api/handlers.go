package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"atlas/capture"
)

// returns the current state of graph in memory
func getGraph(s *capture.Scanner) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, s.GetGraphJSON())
	}
}

// scan for fixed duration, returns success message
func scanDuration(s *capture.Scanner) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// parse query args
		durationArg := ctx.Query("duration")
		durationMS, err := strconv.Atoi(durationArg)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid argument",
			})
		}

		s.ScanDuration(durationMS, false)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "scan completed",
		})
	}
}

// start continuous scan
func startScan(s *capture.Scanner) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if s.IsScanning() {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "already scanning"})
		}

		go s.Start(false) // start continuous scan in new goroutine
		ctx.JSON(http.StatusOK, gin.H{"message": "started"})
	}
}

// stop continuous scan
func stopScan(s *capture.Scanner) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.Stop()
		ctx.JSON(http.StatusOK, gin.H{"message": "stopped"})
	}
}

// lets client open SSE connection, allows backend to push regular updates (stream data)
func streamGraph(s *capture.Scanner) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")

		flusher, ok := ctx.Writer.(http.Flusher)
		if !ok {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// push updates every 2 seconds
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				data, err := json.Marshal(s.GetGraphJSON())
				if err != nil {
					return // json parse error (close)
				}
				// write json to stream
				fmt.Fprintf(ctx.Writer, "data: %s\n\n", data)
				flusher.Flush()
			case <-ctx.Request.Context().Done():
				return // close
			}
		}
	}
}
