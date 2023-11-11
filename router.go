package main

import (
	"context"
	"os"
	"user-storage/controllers"
	docs "user-storage/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
	"user-storage/middlewares"
	"user-storage/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	metadata := models.RequestMetadata{
		UserAgent: req.RequestContext.Identity.UserAgent,
		SourceIP:  req.RequestContext.Identity.SourceIP,
	}

	// ctx = context.WithValue(ctx, "UserAgent", req.RequestContext.Identity.UserAgent)
	// ctx = context.WithValue(ctx, "SourceIP", req.RequestContext.Identity.SourceIP)
	ctx = context.WithValue(ctx, "RequestMetadata", metadata)

	return ginLambda.ProxyWithContext(ctx, req)
}

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func InitRoutes() {
	router := gin.New()

	// router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(middlewares.LoggingMiddleware())

    health := new(controllers.HealthController)
	user := controllers.NewUserController(*models.DB)

    docs.SwaggerInfo.BasePath = "docs"

	v1 := router.Group("/users")

    healthGroup := v1.Group("/health")
    healthGroup.GET("", health.CheckHealth)

	// Account Routes
	usersGroup := v1.Group("/accounts")

	usersGroup.GET("", user.GetAllUsers)
	usersGroup.GET("/paginate", user.GetPaginatedUsers)
	usersGroup.GET("/:id", user.GetUserByID)

	usersGroup.Use(middlewares.DecodeJWT())
	usersGroup.POST("", user.AddUser)
	usersGroup.POST("/with-roles", user.GetUsersWithRole)

	usersGroup.PUT("/:id", user.UpdateUserById)

	usersGroup.DELETE("/:id", user.DeleteUserById)

	// Role Routes
	role := new(controllers.RoleController)

	rolesGroup := v1.Group("/roles")

	rolesGroup.GET("", role.GetAllRoles)
	rolesGroup.GET("/:id", role.GetRoleByID)

	rolesGroup.Use(middlewares.DecodeJWT())

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

    // Swagger
    router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	env := os.Getenv("ENV")
	if env == "lambda" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run()
	}
}
