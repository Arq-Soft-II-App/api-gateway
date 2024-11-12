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

	courseGroup := engine.Group("/courses")
	{
		CoursesRoutes(courseGroup, controller.Courses)
	}
	commentsGroup := engine.Group("/comment")
	{
		CommentsRoutes(commentsGroup, controller.Comments)
	}

	engine.POST("/category/create", controller.Categories.CreateCategory)
	engine.GET("/categories", controller.Categories.GetCategories)
}
