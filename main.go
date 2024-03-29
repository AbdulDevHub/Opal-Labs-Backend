package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/opalescencelabs/backend/controllers"
	"github.com/opalescencelabs/backend/controllers/auth"
	"github.com/opalescencelabs/backend/initializers"
)

// Initialize environment variables, connections to database and Redis
func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.InitializeRedis()
}

// Start application
func main() {
	gin := gin.Default()

	// Use CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{auth.GetFrontendURL()}
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	gin.Use(cors.New(config))

	gin.POST("/user-login", controllers.UserLogin)
	gin.POST("/user-logout", controllers.UserLogout)
	gin.GET("/user-get", controllers.UserGet)

	gin.POST("/page-create", controllers.PageCreate)
	gin.POST("/page-update", controllers.PageUpdate)
	gin.GET("/page-get/:page_uuid", controllers.PageGet)
	gin.GET("/page-list", controllers.PageList)
	gin.POST("/page-delete", controllers.PageDelete)

	if err := gin.Run(); err != nil {
		panic("Router failed to start Gin: " + err.Error())
	}

}
