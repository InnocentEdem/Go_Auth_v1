package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/InnocentEdem/Go_Auth_v1/controllers"
	_ "github.com/InnocentEdem/Go_Auth_v1/docs"
	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/middlewares"
	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvironment()
	initializers.ConnectDB()
	// initializers.SyncDatabase()
}

// @title           Gatekeeper Pro API
// @version         1.0
// @description     GateKeeper Pro registers and validates users for your frontend.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      https://go-auth-v1.onrender.com
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())
	clientRoutes := r.Group("/")
	{
		clientRoutes.POST("/create-client", controllers.ClientSignup)
		clientRoutes.POST("/client-login", controllers.ClientLogin)
		clientRoutes.POST("/generate-apn", middlewares.ClientAuthMiddleware(), controllers.GenerateAPN)
		clientRoutes.POST("/invalidate-apn", middlewares.ClientAuthMiddleware(), controllers.InvalidateAPN)
		clientRoutes.GET("/all-users", middlewares.ClientAuthMiddleware(), controllers.GetUsersByClient)
		clientRoutes.GET("/client-apn", middlewares.ClientAuthMiddleware(), controllers.GetClientAPN)
		clientRoutes.POST("/delete-user", middlewares.ClientAuthMiddleware(), controllers.DeleteUserByClient)
		clientRoutes.POST("/feature-request", middlewares.ClientAuthMiddleware(), controllers.HandleFeatureRequest)
		clientRoutes.GET("/config", middlewares.ClientAuthMiddleware(), controllers.GetClientAdvancedConfig)
		clientRoutes.PUT("/config", middlewares.ClientAuthMiddleware(), controllers.UpdateClientAdvancedConfigHandler)
		clientRoutes.GET("/client", middlewares.ClientAuthMiddleware(), controllers.GetClient)
		clientRoutes.POST("/update-password", middlewares.ClientAuthMiddleware(), controllers.ClientUpdatePassword)
		clientRoutes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	userRoutes := r.Group("/user")
	userRoutes.Use(middlewares.APNAuthMiddleware())
	{
		userRoutes.POST("/signup", controllers.UserSignup)
		userRoutes.POST("/login", controllers.UserLogin)
		userRoutes.GET("/validate", middlewares.UserAuthMiddleware(), controllers.ValidateUser)
		userRoutes.POST("/update-password", middlewares.UserAuthMiddleware(), controllers.UserUpdatePassword)
		userRoutes.GET("/profile", middlewares.UserAuthMiddleware(), controllers.GetUserProfile)
		userRoutes.PUT("/profile", middlewares.UserAuthMiddleware(), controllers.UpdateUserProfile)
		userRoutes.POST("/refresh-token", middlewares.UserAuthMiddleware(), controllers.RefreshToken)
		userRoutes.POST("/validate-email-code", middlewares.UserAuthMiddleware(), controllers.ValidateConfirmationCode)
	}

	r.Run()
}