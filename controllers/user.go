package controllers

import (
	"context"
	"fmt"
	"net/http"
	"user-storage/cache"
	"user-storage/models"
	"user-storage/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"gorm.io/gorm"
)
type UserController struct{
    DB  *gorm.DB
    UserService *services.UserService
}

func NewUserController(db gorm.DB) *UserController {
    return &UserController{
        DB: &db,
        UserService: services.NewUserService(&db),
    }
}

var validate = validator.New()

func (t UserController) GetAllUsers(c *gin.Context) {
	var users []models.User

    code, err := t.UserService.GetAllUsers(&users)
	if err != nil {
		c.JSON(code, models.HTTPError{
			Code:    code,
			Message: fmt.Sprintf("Error getting data. %v", err.Error()),
		})
		return
	}
	c.JSON(code, users)
}

func (t UserController) GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")

    code, err := t.UserService.GetUserByID(&user, id)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: fmt.Sprintf("Failed to retrieve user: %v", err.Error()),
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

    code, err := t.UserService.AddUser(&user)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: fmt.Sprintf("Unable to create user. %v", err.Error()),
        })
        return
    }
  
    c.Set("user", user)
	c.JSON(code, user)
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

    code, err := t.UserService.UpdateUserById(&user, id)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: fmt.Sprintf("Unable to create user. %v", err.Error()),
        })
        return
    }

    deletedCount, cacheErr := cache.RedisClient.Del(context.Background(), id).Result()
	if cacheErr != nil {
		fmt.Println("Error:", cacheErr)
	} else {
		fmt.Printf("Deleted %d keys\n", deletedCount)
	}
    
    // Return the updated user
    c.Set("updatedUser", user)
    c.JSON(code, user)
}

func (t UserController) DeleteUserById(c *gin.Context) {
    id := c.Param("id")

    existingUser := models.User{}
    code, err := t.UserService.DeleteUserById(&existingUser, id)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: fmt.Sprintf("Unable to delete user. %v", err.Error()),
        })
        return
    }

    deletedCount, cacheErr := cache.RedisClient.Del(context.Background(), id).Result()
	if cacheErr != nil {
		fmt.Println("Error:", cacheErr)
	} else {
		fmt.Printf("Deleted %d keys\n", deletedCount)
	}

    c.Set("user", existingUser)
    fmt.Print(existingUser)
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
	err := t.DB.Where("role IN ?", roles).Find(&users)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting data. %v", err.Error.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}
