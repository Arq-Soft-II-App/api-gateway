package routes

import (
	"api-gateway/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, controller *controllers.Controller) {

	authGroup := engine.Group("/auth")
	{
		AuthRoutes(authGroup, controller.Auth)
	}

	userGroup := engine.Group("/users")
	{
		UserRoutes(userGroup, controller.Users)
	}

	/* 	courseGroup := engine.Group("/courses")
	   	{
	   	     CourseRoutes(courseGroup, controller.Course)
	   	} */

}
