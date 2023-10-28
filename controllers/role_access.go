package controllers

import (
	"fmt"
	"net/http"
	"user-storage/models"

	"github.com/gin-gonic/gin"
)

type RoleAccessController struct {}

func (t RoleAccessController) GetAllRoleAccesses(c *gin.Context) {
	var roleAccesses []models.RoleAccess
	err := models.DB.Find(&roleAccesses)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting data. %v", err.Error.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, roleAccesses)
}

func (t RoleAccessController) AddRoleAccess(c *gin.Context) {
	var roleAccess models.RoleAccess

	if err := c.BindJSON(&roleAccess); err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to create role access. %v" , err.Error()),
		})
		return
	}
	if err := validate.Struct(&roleAccess); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf("Unable to create role access. %v" , err.Error()),
		})
		return
	}
	if result := models.DB.Create(&roleAccess); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to create role access. %v" , result.Error.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, roleAccess)
}

func (t RoleAccessController) DeleteRoleAccess(c *gin.Context) {
	var roleAccess models.RoleAccess
	
	if err := c.BindJSON(&roleAccess); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf("Missing required fields. %v" , err.Error()),
		})
		return
	}

	if err := validate.Struct(&roleAccess); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf("Incorrect data type. %v" , err.Error()),
		})
		return
	}

	existingRoleAccess := models.RoleAccess{}
	if result := models.DB.Where("role_id = ? AND ap_id = ?", roleAccess.RoleId, roleAccess.APId).First(&existingRoleAccess); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code: http.StatusNotFound,
			Message: "Role access is not found",
		})
		return
	}

	fmt.Println(existingRoleAccess)

	if result := models.DB.Where("role_id = ? AND ap_id = ?", roleAccess.RoleId, roleAccess.APId).Delete(&existingRoleAccess); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message:fmt.Sprintf("Unable to delete role access. %v" , result.Error.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, "Success")
}

