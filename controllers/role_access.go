package controllers

import (
	"fmt"
	"net/http"
	"user-storage/models"

	"github.com/gin-gonic/gin"
)

type RoleAccessController struct {}

//  @Summary        Get all Role Accesses
//  @Description    Retrieves a list of Role Access
//  @Tags           role-access
//  @Produce        json
//  @Success        200     {array}     models.RoleAccess
//  @Failure        500     {object}    models.HTTPError
//  @Router         /role-access  [get]
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

//  @Summary        Add a Role Access
//  @Description    Add a Role Access into Database
//  @Tags           role-access
//  @Produce        json
//  @Param          role-access    body        models.RoleAccess     true    "Role Access Details"
//  @Success        200     {object}    models.RoleAccess
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON body"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /role-access [post]
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

//  @Summary        Delete a Role Access
//  @Description    Delete a Role Access
//  @Tags           role-access
//  @Produce        json
//  @Param          role_id    body    string  true    "Role ID"
//  @Param          ap_id      body    string  true    "Access Point ID"
//  @Success        200     "Success"   
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON object"
//  @Failure        404     {object}    models.HTTPError    "Role access is not found given role_id and ap_id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /role-access   [delete]
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

	if result := models.DB.Where("role_id = ? AND ap_id = ?", roleAccess.RoleId, roleAccess.APId).Delete(&existingRoleAccess); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code: http.StatusInternalServerError,
			Message:fmt.Sprintf("Unable to delete role access. %v" , result.Error.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, "Success")
}

