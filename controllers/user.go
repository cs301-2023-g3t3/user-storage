package controllers

import (
	// "context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	// "user-storage/cache"
	"user-storage/models"
	"user-storage/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"gorm.io/gorm"
)

type UserController struct {
	DB          *gorm.DB
	UserService *services.UserService
}

func NewUserController(db gorm.DB) *UserController {
	return &UserController{
		DB:          &db,
		UserService: services.NewUserService(&db),
	}
}

var validate = validator.New()


//  @Summary        Get all Users
//  @Description    Retrieves a list of users
//  @Tags           users
//  @Produce        json
//  @Success        200     {array}     models.User
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts   [get]
func (t UserController) GetAllUsers(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	role := c.DefaultQuery("role", "0")
	name := c.DefaultQuery("name", "")
	email := c.DefaultQuery("email", "")

	roleInt, err := strconv.Atoi(role)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid page parameter",
        })
        return
    }
	users, code, err := t.UserService.GetAllUsers(roleInt, id, name ,email)
	if err != nil {
		c.JSON(code, models.HTTPError{
			Code:    code,
			Message: fmt.Sprintf("Error getting data. %v", err.Error()),
		})
		return
	}
	c.JSON(code, *users)
}

//  @Summary        Get all Users by Pagination
//  @Description    Retrieves a list of users
//  @Tags           users
//  @Produce        json
//  @Param          page    query   int     true    "page"
//  @Param          size    query   int     true    "size"
//  @Success        200     {array}     models.User
//  @Failure        400     {object}    models.HTTPError    "Invalid parameters"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/paginate   [get]
func (t UserController) GetPaginatedUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("size", "10")

	// Convert page and pageSize to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid page parameter",
		})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid pageSize parameter",
		})
		return
	}

	users, code, err := t.UserService.GetPaginatedUsers(pageInt, pageSizeInt)
	if err != nil {
		c.JSON(code, models.HTTPError{
			Code:    code,
			Message: fmt.Sprintf("Error getting data. %v", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       users,
		"pagination": gin.H{"page": page, "size": pageSize},
	})
}

//  @Summary        Get User by Id
//  @Description    Retrieve a User By UserID
//  @Tags           users
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     {array}     models.User
//  @Failure        400     {object}    models.HTTPError    "UserId cannot be empy"
//  @Failure        404     {object}    models.HTTPError    "User not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/{id}   [get]
func (t UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, code, err := t.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(code, models.HTTPError{
			Code:    code,
			Message: fmt.Sprintf("Failed to retrieve user: %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, *user)
}

//  @Summary        Add a User
//  @Description    Add a User into Database
//  @Tags           users
//  @Produce        json
//  @Param          user    body        models.User     true    "User Details"
//  @Success        200     {object}    models.User
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON body"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts   [post]
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
			Code:    code,
			Message: fmt.Sprintf("Unable to create user. %v", err.Error()),
		})
		return
	}

	c.Set("user", *res)
	c.JSON(code, *res)
}

//  @Summary        Update User Details by Id
//  @Description    Update a User By UserID
//  @Tags           users
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     {object}    models.User
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON body"
//  @Failure        404     {object}    models.HTTPError    "User not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/{id}   [put]
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
			Code:    code,
			Message: fmt.Sprintf("Unable to update user. %v", err.Error()),
		})
		return
	}

	// deletedCount, cacheErr := cache.RedisClient.Del(context.Background(), id).Result()
	// if cacheErr != nil {
	// 	fmt.Println("Error:", cacheErr)
	// } else {
	// 	fmt.Printf("Deleted %d keys\n", deletedCount)
	// }

	// Return the updated user
	c.Set("updatedUser", *res)
	c.JSON(code, *res)
}

//  @Summary        Delete a User by Id
//  @Description    Delete a User By UserID
//  @Tags           users
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     "Success"   
//  @Failure        400     {object}    models.HTTPError    "Bad request due to empty string Id"
//  @Failure        404     {object}    models.HTTPError    "User not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/{id}   [delete]
func (t UserController) DeleteUserById(c *gin.Context) {
	id := c.Param("id")

	res, code, err := t.UserService.DeleteUserById(id)
	if err != nil {
		c.JSON(code, models.HTTPError{
			Code:    code,
			Message: fmt.Sprintf("Unable to delete user. %v", err.Error()),
		})
		return
	}

	// deletedCount, cacheErr := cache.RedisClient.Del(context.Background(), id).Result()
	// if cacheErr != nil {
	// 	fmt.Println("Error:", cacheErr)
	// } else {
	// 	fmt.Printf("Deleted %d keys\n", deletedCount)
	// }

	c.Set("user", *res)
	c.JSON(http.StatusOK, gin.H{
		"data":       "Success",
	})
}


type Input struct {
    Roles []int `json:"roles" validate:"required"`
}

//  @Summary        Get a list of users with roles
//  @Description    Get a list of users with roles
//  @Tags           users
//  @Produce        json
//  @Param          user    body    Input    true    "roles"
//  @Success        200     "Success"   
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON body"
//  @Failure        404     {object}    models.HTTPError    "Cannot find users with given roles"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/with-roles   [post]
func (t UserController) GetUsersWithRole(c *gin.Context) {
    var input Input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
        return
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
