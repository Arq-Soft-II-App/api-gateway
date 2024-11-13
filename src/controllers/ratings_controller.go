package controllers

import (
	"api-gateway/src/dto/ratings"
	"api-gateway/src/services"
	"fmt"
	"net/http"

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
	var rating ratings.RatingDTO
	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdRating, err := c.service.NewRating(rating)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdRating)
}

func (c *RatingsController) UpdateRating(ctx *gin.Context) {
	var rating ratings.RatingDTO
	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedRating, err := c.service.UpdateRating(rating)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedRating)
}

func (c *RatingsController) GetAllRatings(ctx *gin.Context) {
	fmt.Println("GetAllRatings called")
	ratingsList, err := c.service.GetAllRatings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ratingsList)
}
