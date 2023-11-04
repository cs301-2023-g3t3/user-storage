package controllers

import (
	"fmt"
	"net/http"
	"user-storage/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserController struct{}

var validate = validator.New()

func (t UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	err := models.DB.Find(&users)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting data. %v", err.Error.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (t UserController) GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "User ID cannot be empty",
		})
		return
	}

	if result := models.DB.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: "User ID is not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (t UserController) AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to create user. %v", err.Error()),
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
    
    c.Set("user", user)
	  c.JSON(http.StatusCreated, user)
}

func (t UserController) UpdateUserById(c *gin.Context) {
    var user models.User
    id := c.Param("id")
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code:    http.StatusBadRequest,
            Message: fmt.Sprintf("Invalid payload. %v", err.Error()),
        })
        return
    }
    c.Set("user", user)

    user.Id = id

    // Check if the user exists in the database
    existingUser := models.User{}
    if err := models.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
        c.JSON(http.StatusNotFound, models.HTTPError{
            Code:    http.StatusNotFound,
            Message: "User not found with given ID",
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
    c.Set("updatedUser", user)
    c.JSON(http.StatusOK, existingUser)
}

func (t UserController) DeleteUserById(c *gin.Context) {
    id := fmt.Sprint(c.Param("id"))
    if id == "" {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code: http.StatusBadRequest,
            Message: "User ID cannot be empty",
        })
    }

    existingUser := models.User{}
    if err := models.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
        c.JSON(http.StatusNotFound, models.HTTPError{
            Code:    http.StatusNotFound,
            Message: "User not found with given ID",
        })
        return
    }

    err := models.DB.Where("id = ?", id).Delete(&existingUser)
    if err.Error != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code: http.StatusBadRequest,
            Message: fmt.Sprintf("Unable to delete data. %v", err.Error.Error()),
        })
        return
    }
    c.Set("user", existingUser)
    c.JSON(http.StatusOK, "Success")
}

func (t UserController) GetUsersWithRole(c *gin.Context) {
	var users []models.User
	var input struct {
		Roles []int `json:"roles"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roles := input.Roles
	err := models.DB.Where("role IN ?", roles).Find(&users)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting data. %v", err.Error.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}
