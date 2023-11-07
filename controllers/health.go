package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (t HealthController) CheckHealth(c *gin.Context) {
	log.Println("Checking Health")
	c.String(http.StatusOK, "Success")
}
