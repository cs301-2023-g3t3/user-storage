// Copied from https://github.com/cs301-2023-g3t3/points-ledger/blob/main/middlewares/loggerMiddleware.go
package middlewares

import (
	// "fmt"
	"net/http"
	"strings"
	"time"
	"user-storage/models"

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
		latencyTime := endTime.Sub(startTime).Milliseconds()

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
		}

		// Logs for user accounts
		if strings.Contains(reqUri, "/accounts") && !strings.Contains(reqUri, "/accounts/with-roles") {
			if reqMethod == http.MethodPost || reqMethod == http.MethodPut || reqMethod == http.MethodDelete {

				user, _ := ctx.Get("user")
				userValue, _ := user.(models.User)

				data, ok := ctx.Get("userDetails")
				if !ok {
					ctx.JSON(http.StatusInternalServerError, models.HTTPError{
						Code:    http.StatusInternalServerError,
						Message: "Error",
					})
					ctx.Abort()

				}
				userDetailsObj, ok := data.(map[string]interface{})
				if !ok {
					ctx.JSON(http.StatusInternalServerError, models.HTTPError{
						Code:    http.StatusInternalServerError,
						Message: "Error",
					})
					ctx.Abort()
				}
				var action string
				var updatedUserFields log.Fields
				if reqMethod == http.MethodPost {
					action = "add user"
				} else if reqMethod == http.MethodPut {
					action = "update user details"
					newUser, _ := ctx.Get("updatedUser")
					newUserValue, _ := newUser.(models.User)
					updatedUserFields = log.Fields{
						"id": newUserValue.Id,
						// "firstName": newUserValue.FirstName,
						// "lastName":  newUserValue.LastName,
						// "email":     newUserValue.Email,
						"role": newUserValue.Role,
					}
				} else if reqMethod == http.MethodDelete {
					action = "delete user"
				}

				userFields := log.Fields{
					"id": userValue.Id,
					// "firstName": userValue.FirstName,
					// "lastName":  userValue.LastName,
					// "email":     userValue.Email,
					"role": userValue.Role,
				}

				log.WithFields(log.Fields{
					"METHOD":               reqMethod,
					"URI":                  reqUri,
					"STATUS":               statusCode,
					"LATENCY":              latencyTime,
					"ACTOR":                userDetailsObj["user_id"],
					"USER_DETAILS":         userFields,
					"UPDATED_USER_DETAILS": updatedUserFields,
					"ACTION":               action,
					"USER_AGENT":           userAgent,
					"SOURCE_IP":            sourceIP,
				}).Info("USER DETAILS REQUEST")
			}
		}

		// Logs for roles
		if strings.Contains(reqUri, "/roles") {
			role, _ := ctx.Get("role")
			roleValue, _ := role.(models.Role)

			roleFields := log.Fields{
				"id":   roleValue.Id,
				"name": roleValue.Name,
			}

			if reqMethod == http.MethodPost || reqMethod == http.MethodPut || reqMethod == http.MethodDelete {
				var action string
				var updatedRoleFields log.Fields

				if reqMethod == http.MethodPost {
					action = "add role"
				} else if reqMethod == http.MethodPut {
					action = "update role"
					newRole, _ := ctx.Get("updatedRole")
					newRoleValue, _ := newRole.(models.Role)
					updatedRoleFields = log.Fields{
						"id":   newRoleValue.Id,
						"name": newRoleValue.Name,
					}
				} else if reqMethod == http.MethodDelete {
					action = "delete role"
				}

				log.WithFields(log.Fields{
					"METHOD":               reqMethod,
					"URI":                  reqUri,
					"STATUS":               statusCode,
					"LATENCY":              latencyTime,
					"ROLE_DETAILS":         roleFields,
					"UPDATED_ROLE_DETAILS": updatedRoleFields,
					"ACTION":               action,
					"USER_AGENT":           userAgent,
					"SOURCE_IP":            sourceIP,
				}).Info("ROLE REQUEST")
			}
		}

		// if reqMethod == http.MethodGet {
		// 	log.WithFields(log.Fields{
		// 		"METHOD":     reqMethod,
		// 		"URI":        reqUri,
		// 		"STATUS":     statusCode,
		// 		"LATENCY":    latencyTime,
		// 		"USER_AGENT": userAgent,
		// 		"CLIENT_IP":  sourceIP,
		// 	}).Info("HTTP GET REQUEST")
		// }
	}
}
