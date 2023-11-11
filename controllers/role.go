package controllers

import (
	"fmt"
	"net/http"
	"user-storage/models"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

//  @Summary        Get all Roles
//  @Description    Retrieves a list of Roles 
//  @Tags           roles
//  @Produce        json
//  @Success        200     {array}     models.Role
//  @Failure        500     {object}    models.HTTPError
//  @Router         /roles  [get]
func (t RoleController) GetAllRoles(c *gin.Context) {
	var roles []models.Role
	err := models.DB.Find(&roles)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting data. %v", err.Error.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, roles)
}

//  @Summary        Get Role by Id
//  @Description    Retrieve a Role By RoleID
//  @Tags           roles
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     {object}    models.Role
//  @Failure        400     {object}    models.HTTPError    "RoleId cannot be empy"
//  @Failure        404     {object}    models.HTTPError    "Role not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /roles/{id}   [get]
func (t RoleController) GetRoleByID(c *gin.Context) {
	var roles models.Role
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Role ID cannot be empty",
		})
		return
	}

	if result := models.DB.Find(&roles, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Role ID is not found",
		})
		return
	}

	c.JSON(http.StatusOK, roles)
}

//  @Summary        Add a Role
//  @Description    Add a Role into Database
//  @Tags           roles
//  @Produce        json
//  @Param          role    body        models.Role     true    "Role Details"
//  @Success        200     {object}    models.Role
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON body"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /roles  [post]
func (t RoleController) AddRole(c *gin.Context) {
	var role models.Role
	if err := c.BindJSON(&role); err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Missing required fields. %v", err.Error()),
		})
		return
	}
	if err := validate.Struct(&role); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Unable to create role. %v", err.Error()),
		})
		return
	}
	if result := models.DB.Create(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to create role. %v", result.Error.Error()),
		})
		return
	}
	c.Set("role", role)
	c.JSON(http.StatusOK, role)
}

//  @Summary        Update Role Details by Id
//  @Description    Update a Role By RoleId
//  @Tags           roles
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     {object}    models.Role
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON body"
//  @Failure        404     {object}    models.HTTPError    "Role not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /roles/{id}   [put]
func (t RoleController) UpdateRoleById(c *gin.Context) {
	var role models.Role
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Role ID cannot be empty",
		})
		return
	}
	if err := c.BindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Missing required fields. %v", err.Error()),
		})
		return
	}
	if err := validate.Struct(&role); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Incorrect data type. %v", err.Error()),
		})
		return
	}

	// Check if role exists
	existingRole := models.Role{}
	if err := models.DB.Where("id = ?", id).First(&existingRole).Error; err != nil {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Role not found with given ID",
		})
		return
	}
	c.Set("role", existingRole)

	if result := models.DB.Model(&existingRole).Updates(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to update role. %v", result.Error.Error()),
		})
		return
	}

	c.Set("updatedRole", role)
	c.JSON(http.StatusOK, existingRole)
}

//  @Summary        Delete a Role by Id
//  @Description    Delete a Role By RoleID
//  @Tags           roles
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     "Success"   
//  @Failure        400     {object}    models.HTTPError    "Bad request due to empty string Id"
//  @Failure        404     {object}    models.HTTPError    "Role not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /roles/{id}   [delete]
func (t RoleController) DeleteRoleById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Role ID cannot be empty",
		})
		return
	}
	role := models.Role{}
	if err := models.DB.Where("id = ?", id).First(&role).Error; err != nil {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Role not found with given ID",
		})
		return
	}
	if result := models.DB.Delete(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to delete role. %v", result.Error.Error()),
		})
		return
	}

	c.Set("role", role)
	c.JSON(http.StatusOK, "Success")
}
