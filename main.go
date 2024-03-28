package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opalescencelabs/backend/controllers"
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
	gin.Use(CORSMiddleware())

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
