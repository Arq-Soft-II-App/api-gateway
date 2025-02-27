package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

func CoursesRoutes(r *gin.RouterGroup, controller controllers.CourseControllerInterface, searchController controllers.SearchControllerInterface) {
	r.POST("/create", controller.CreateCourse)
	r.PUT("/update", middlewares.AuthMiddleware(), controller.UpdateCourse)
	r.GET("/", searchController.SearchCourses)
	r.GET("/:id", controller.GetCourseById)
	r.DELETE("/:id", controller.DeleteCourse)
}
