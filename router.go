package main

import (
	"context"
	"os"
	"user-storage/controllers"
	"user-storage/middlewares"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func InitRoutes() {
	router := gin.Default()

	router.Use(cors.Default())

	user := new(controllers.UserController)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware())

	v1 := router.Group("/users")

	// Account Routes
	usersGroup := v1.Group("/accounts")

	usersGroup.GET("", user.GetAllUsers)
	usersGroup.GET("/:id", user.GetUserByID)

	usersGroup.POST("", user.AddUser)
	usersGroup.POST("/with-roles", user.GetUsersWithRole)

	usersGroup.PUT("/:id", user.UpdateUserById)

	usersGroup.DELETE("/:id", user.DeleteUserById)


	// Role Routes
	role := new(controllers.RoleController)

	rolesGroup := v1.Group("/roles")

	rolesGroup.GET("", role.GetAllRoles)
	rolesGroup.GET("/:id", role.GetRoleByID)

	rolesGroup.POST("", role.AddRole)

	rolesGroup.PUT("/:id", role.UpdateRoleById)

	rolesGroup.DELETE("/:id", role.DeleteRoleById)

	// Access Points
	accessPoint := new(controllers.AccessPointController)

	accessPointsGroup := v1.Group("/access-points")

	accessPointsGroup.GET("", accessPoint.GetAllAccessPoints)
	accessPointsGroup.GET("/:id", accessPoint.GetAccessPointByID)

	accessPointsGroup.POST("", accessPoint.AddAccessPoint)

	accessPointsGroup.PUT("/:id", accessPoint.UpdateAccessPoint)

	accessPointsGroup.DELETE("/:id", accessPoint.DeleteAccessPoint)

	// Role Access
	roleAccess := new(controllers.RoleAccessController)

	roleAccessesGroup := v1.Group("/role-access")

	roleAccessesGroup.GET("", roleAccess.GetAllRoleAccesses)

	roleAccessesGroup.POST("", roleAccess.AddRoleAccess)

	roleAccessesGroup.DELETE("", roleAccess.DeleteRoleAccess)

	env := os.Getenv("ENV")
	if env == "lambda" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run()
	}
}
