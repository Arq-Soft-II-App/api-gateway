package controllers

import (
	"api-gateway/src/dto/courses"
	"api-gateway/src/errors"
	"api-gateway/src/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	service services.SearchServiceInterface
}

func NewSearchController(service services.SearchServiceInterface) *SearchController {
	return &SearchController{
		service: service,
	}
}

type SearchControllerInterface interface {
	SearchCourses(c *gin.Context)
}

func (c *SearchController) SearchCourses(ctx *gin.Context) {
	fmt.Println("Controller: Searching courses...")
	query := ctx.Query("q")
	results, err := c.service.SearchCourses(query)
	if err != nil {
		if errors.GetStatusCode(err) == http.StatusNotFound {
			ctx.JSON(http.StatusNotFound, []courses.CourseListDto{})
			return
		}
		ctx.JSON(errors.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, results)
}
