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
	r := gin.Default()

	// Use CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{auth.GetFrontendURL(), auth.GetDomain()}
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	r.Use(cors.New(config))

	r.POST("/user-login", controllers.UserLogin)
	r.POST("/user-logout", controllers.UserLogout)
	r.GET("/user-get", controllers.UserGet)

	r.POST("/page-create", controllers.PageCreate)
	r.POST("/page-update", controllers.PageUpdate)
	r.GET("/page-get/:page_uuid", controllers.PageGet)
	r.GET("/page-list", controllers.PageList)
	r.POST("/page-delete", controllers.PageDelete)

	if err := r.Run(); err != nil {
		panic("Router failed to start Gin: " + err.Error())
	}

}
