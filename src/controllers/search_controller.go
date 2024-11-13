package controllers

import (
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
	fmt.Println("SearchCourses called")
	query := ctx.Query("query")
	results, err := c.service.SearchCourses(query)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, results)
}
