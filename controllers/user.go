package controllers

import (
	"context"
	"encoding/json"
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
    users, code, err := t.UserService.GetAllUsers()
	if err != nil {
		c.JSON(code, models.HTTPError{
			Code:    code,
			Message: fmt.Sprintf("Error getting data. %v", err.Error()),
		})
		return
	}
	c.JSON(code, *users)
}

func (t UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

    user, code, err := t.UserService.GetUserByID(id)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: fmt.Sprintf("Failed to retrieve user: %v", err.Error()),
        })
        return
    }

	c.JSON(http.StatusOK, *user)
}

func (t UserController) AddUser(c *gin.Context) {
    var user models.User
    decoder := json.NewDecoder(c.Request.Body)
    if err := decoder.Decode(&user); err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code:    http.StatusBadRequest,
            Message: fmt.Sprintf("Invalid JSON request: %v", err.Error()),
        })
        return
    }

    res, code, err := t.UserService.AddUser(&user)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: fmt.Sprintf("Unable to create user. %v", err.Error()),
        })
        return
    }

    c.Set("user", *res)
	c.JSON(code, *res)
}

func (t UserController) UpdateUserById(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    decoder := json.NewDecoder(c.Request.Body)
    if err := decoder.Decode(&user); err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code:    http.StatusBadRequest,
            Message: fmt.Sprintf("Invalid JSON request: %v", err.Error()),
        })
        return
    }

    c.Set("user", user)

    res, code, err := t.UserService.UpdateUserById(&user, id)
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
    c.Set("updatedUser", *res)
    c.JSON(code, *res)
}

func (t UserController) DeleteUserById(c *gin.Context) {
    id := c.Param("id")

    res, code, err := t.UserService.DeleteUserById(id)
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

    c.Set("user", *res)
    c.JSON(http.StatusOK, "Success")
}

func (t UserController) GetUsersWithRole(c *gin.Context) {
	var input struct {
        Roles []int `json:"roles" validate:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    if err := validate.Struct(input); err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code: http.StatusBadRequest,
            Message: err.Error(),
        })
    }

	roles := input.Roles
    res, code, err := t.UserService.GetUsersWithRole(roles)
	if err != nil {
		c.JSON(code, models.HTTPError{
			Code:    code,
			Message: fmt.Sprintf("Error getting data. %v", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, *res)
}
