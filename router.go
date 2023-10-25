package main

import (
	"context"
	"os"
	"user-storage/controllers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func InitRoutes() {
	router := gin.Default()

    user := new(controllers.UserController)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")

	usersGroup := v1.Group("/users")

	usersGroup.GET("", user.GetAllUsers)
	usersGroup.GET("/:id", user.GetUserByID)

	usersGroup.POST("", user.AddUser)

	usersGroup.PUT("/:id", user.UpdateUserById)

    usersGroup.DELETE("/:id", user.DeleteUserById)

	env := os.Getenv("ENV")
	if env == "lambda" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run()
	}
}
