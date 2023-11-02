// Copied from https://github.com/cs301-2023-g3t3/points-ledger/blob/main/middlewares/loggerMiddleware.go
package middlewares

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Starting time request
		startTime := time.Now()

		// Process the request
		ctx.Next()

		// End Time request
		endTime := time.Now()

		// Execution time
		latencyTime := endTime.Sub(startTime)

		// Request data
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		// userAgent := ctx.GetHeader("User-Agent")
		metadata, ok := ctx.Request.Context().Value("RequestMetadata").(models.RequestMetadata)
		var userAgent string
		var sourceIP string
		if ok {
			// Access the UserAgent and SourceIP
			userAgent = metadata.UserAgent
			sourceIP = metadata.SourceIP

			fmt.Println("UserAgent in middleware: ", userAgent)
			fmt.Println("SourceIP in middleware: ", sourceIP)
		}

		// Request IP
		if reqMethod == http.MethodPost {
			log.WithFields(log.Fields{
				"METHOD":     reqMethod,
				"URI":        reqUri,
				"STATUS":     statusCode,
				"LATENCY":    latencyTime,
				"USER_AGENT": userAgent,
				"SOURCE_IP":  sourceIP,
			}).Info("HTTP POST REQUEST")
		}

		if reqMethod == http.MethodGet {
			log.WithFields(log.Fields{
				"METHOD":     reqMethod,
				"URI":        reqUri,
				"STATUS":     statusCode,
				"LATENCY":    latencyTime,
				"USER_AGENT": userAgent,
				"CLIENT_IP":  sourceIP,
			}).Info("HTTP GET REQUEST")
		}
	}
}