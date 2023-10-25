package controllers

import (
	"fmt"
	"net/http"
	"user-storage/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

func GetUsers(c *gin.Context) {
	var users []models.User
    err := models.DB.Find(&users)
    if err.Error != nil {
        c.JSON(http.StatusInternalServerError, models.HTTPError{
            Code: http.StatusInternalServerError,
            Message: fmt.Sprintf("Error getting data. %v", err.Error.Error()),
        })
        return
    }
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("ID")
    if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
            Code: http.StatusBadRequest,
            Message: "User ID cannot be empty",
        })
		return
    }

	if result := models.DB.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{
            Code: http.StatusNotFound,
            Message: "User ID is not found",
        })
		return
	}

	c.JSON(http.StatusOK, user)
}

func AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusInternalServerError, models.HTTPError{
            Code: http.StatusInternalServerError,
            Message: fmt.Sprintf("Unable to create user. %v" , err.Error()),
        })
		return
	}

    // Validate the inputs
    if err := validate.Struct(&user); err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code: http.StatusBadRequest,
            Message: fmt.Sprintf("Unable to create user. %v" , err.Error()),
        })
        return
    }

    // generate new UUID for user
    user.Id = uuid.NewString()

    err := models.DB.Create(user)
    if err.Error != nil {
        c.JSON(http.StatusInternalServerError, models.HTTPError{
            Code: http.StatusInternalServerError,
            Message: fmt.Sprintf("Unable to create user. %v" , err.Error.Error()),
        })
        return
    }
	c.JSON(http.StatusCreated, user)
}

func EditUser(c *gin.Context) {
    var user models.User
    id := c.Param("ID")
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code:    http.StatusBadRequest,
            Message: fmt.Sprintf("Invalid payload. %v", err.Error()),
        })
        return
    }

    user.Id = id

    // Check if the user exists in the database
    existingUser := models.User{}
    if err := models.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
        c.JSON(http.StatusNotFound, models.HTTPError{
            Code:    http.StatusNotFound,
            Message: "User not found",
        })
        return
    }

    // Update the user's data
    models.DB.Model(&existingUser).Updates(user)

    if models.DB.Error != nil {
        c.JSON(http.StatusInternalServerError, models.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: "Failed to update user",
        })
        return
    }

    // Return the updated user
    c.IndentedJSON(http.StatusOK, existingUser)
}
