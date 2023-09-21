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
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

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

	env := os.Getenv("env")
	if env == "lambda" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run()
	}
}
