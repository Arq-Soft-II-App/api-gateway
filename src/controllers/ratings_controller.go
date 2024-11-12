package controllers

import (
	"api-gateway/src/services"

	"github.com/gin-gonic/gin"
)

type RatingsController struct {
	service services.RatingsServiceInterface
}

func NewRatingsController(service services.RatingsServiceInterface) *RatingsController {
	return &RatingsController{
		service: service,
	}
}

type RatingsControllerInterface interface {
	NewRating(c *gin.Context)
	UpdateRating(c *gin.Context)
	GetAllRatings(c *gin.Context)
}

func (c *RatingsController) NewRating(ctx *gin.Context) {

}

func (c *RatingsController) UpdateRating(ctx *gin.Context) {

}

func (c *RatingsController) GetAllRatings(ctx *gin.Context) {

}
