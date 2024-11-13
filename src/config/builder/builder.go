package builder

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/controllers"
	"api-gateway/src/routes"
	"api-gateway/src/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Build(env envs.Envs) *gin.Engine {

	/* 	corsConfig := cors.DefaultConfig()
	   	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	   	corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"}
	   	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	   	corsConfig.ExposeHeaders = []string{"Content-Length"}
	   	corsConfig.AllowCredentials = true
	   	corsConfig.MaxAge = 12 * time.Hour */

	engine := gin.Default()
	engine.Use(CORSMiddleware())

	service := services.NewService(env)
	controller := controllers.NewController(service)
	routes.SetupRoutes(engine, controller)

	return engine
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("CORS Middleware triggered for %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
