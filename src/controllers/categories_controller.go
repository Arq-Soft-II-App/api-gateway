package controllers

import (
	"api-gateway/src/dto/categories"
	"api-gateway/src/errors"
	"api-gateway/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoriesController struct {
	service services.CategoriesServiceInterface
}

func NewCategoriesController(service services.CategoriesServiceInterface) *CategoriesController {
	return &CategoriesController{
		service: service,
	}
}

type CategoriesControllerInterface interface {
	CreateCategory(c *gin.Context)
	GetCategories(c *gin.Context)
}

func (c *CategoriesController) CreateCategory(ctx *gin.Context) {
	var categoryData categories.CreateCategoryDto
	if err := ctx.ShouldBindJSON(&categoryData); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("Datos inválidos"))
		return
	}

	err := c.service.CreateCategory(categoryData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("Error al crear la categoría: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Categoría creada exitosamente"})
}

func (c *CategoriesController) GetCategories(ctx *gin.Context) {
	categories, err := c.service.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("Error al obtener las categorías: "+err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, categories)
}
