package main

import (
	"log"
	"net/http"

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
	config.AllowOrigins = auth.GetFrontendURLs()
	if len(config.AllowOrigins) == 0 {
		log.Fatal("FRONTEND_URL must be set so the frontend is allowed by CORS")
	}
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowHeaders = append(config.AllowHeaders, "Access-Control-Allow-Origin")
	r.Use(cors.New(config))

	// Lets the frontend wake the server up (free hosting sleeps when idle)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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
