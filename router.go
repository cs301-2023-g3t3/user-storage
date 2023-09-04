package main

import (
	"github.com/gin-gonic/gin"
	"user-storage/controllers"
)

func InitRoutes() {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")

	usersGroup := v1.Group("/users")

	usersGroup.GET("", controllers.GetUsers)
	usersGroup.POST("", controllers.AddUser)
	usersGroup.GET("/:ID", controllers.GetUserByID)
	usersGroup.PUT("/:ID", controllers.EditUser)
	router.Run()
}
