package controllers

import (
	"fmt"
	"net/http"
	"user-storage/models"

	"github.com/gin-gonic/gin"
)

type AccessPointController struct {}

func (t AccessPointController) GetAllAccessPoints(c *gin.Context) {
	var accessPoints []models.AccessPoint
	err := models.DB.Find(&accessPoints)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting data. %v", err.Error.Error()),
		})
		return
	}

	fmt.Println(accessPoints[0])

	c.JSON(http.StatusOK, accessPoints)
}

func (t AccessPointController) GetAccessPointByID(c *gin.Context) {
	var accessPoints models.AccessPoint
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: "AccessPoint ID cannot be empty",
		})
		return
	}

	if result := models.DB.Find(&accessPoints, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code: http.StatusNotFound,
			Message: "AccessPoint ID is not found",
		})
		return
	}

	c.JSON(http.StatusOK, accessPoints)
}

func (t AccessPointController) AddAccessPoint(c *gin.Context) {
	var accessPoint models.AccessPoint
	if err := c.BindJSON(&accessPoint); err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to create accessPoint. %v" , err.Error()),
		})
		return
	}

	if err := validate.Struct(accessPoint); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf("Unable to validate accessPoint. %v", err.Error()),
		})
		return
	}

	if result := models.DB.Create(&accessPoint); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to create accessPoint. %v", result.Error.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, accessPoint)
}

func (t AccessPointController) UpdateAccessPoint(c *gin.Context) {
	var accessPoint models.AccessPoint
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: "AccessPoint ID cannot be empty",
		})
		return
	}

	if err := c.BindJSON(&accessPoint); err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to update accessPoint. %v", err.Error()),
		})
		return
	}

	if err := validate.Struct(accessPoint); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf("Unable to validate accessPoint. %v", err.Error()),
		})
		return
	}

	// Check if access point exists
	existingAP := models.AccessPoint{}
	if err := models.DB.Where("id = ?", id).First(&existingAP).Error; err != nil {
        c.JSON(http.StatusNotFound, models.HTTPError{
            Code:    http.StatusNotFound,
            Message: "Access Point not found with given ID",
        })
        return
    }

	if result := models.DB.Model(&existingAP).Updates(&accessPoint); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to update accessPoint. %v", result.Error.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, existingAP)
}

func (t AccessPointController) DeleteAccessPoint(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: "AccessPoint ID cannot be empty",
		})
		return
	}

	accessPoint := models.AccessPoint{}
	if result := models.DB.Where("id = ?", id).First(&accessPoint); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code: http.StatusNotFound,
			Message: "AccessPoint not found with given ID",
		})
		return
	}

	if result := models.DB.Delete(&accessPoint); result.Error != nil {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code: http.StatusNotFound,
			Message: fmt.Sprintf("Unable to delete accessPoint. %v", result.Error.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, "Success")
}

