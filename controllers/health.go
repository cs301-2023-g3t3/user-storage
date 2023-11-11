package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

//  @Summary        Get Health
//  @Description    Check the health of the service
//  @Tags           health
//  @Produce        json
//  @Success        200     "Sucess"
//  @Router         /health   [get]
func (t HealthController) CheckHealth(c *gin.Context) {
	log.Println("Checking Health")
	c.String(http.StatusOK, "Success")
}
