package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-storage/models"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	var user []models.User
	id := c.Param("ID")
	if result := models.DB.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		return
	}
	models.DB.Create(user)
	c.IndentedJSON(http.StatusCreated, user)
}

func EditUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}
	models.DB.Save(user)
	c.IndentedJSON(http.StatusOK, user)
}
