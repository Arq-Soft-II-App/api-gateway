package routes

import (
	"api-gateway/src/controllers"

	"github.com/gin-gonic/gin"
)

func CommentsRoutes(r *gin.RouterGroup, controller controllers.CommentsControllerInterface) {
	r.POST("/", controller.CreateComment)
	r.GET("/:cid", controller.GetCourseComments)
	r.PUT("/", controller.UpdateComment)
}
