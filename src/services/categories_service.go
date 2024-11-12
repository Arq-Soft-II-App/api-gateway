package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/categories"
	"api-gateway/src/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CategoriesService struct {
	CategoriesInterface CategoriesServiceInterface
	env                 envs.Envs
}

type CategoriesServiceInterface interface {
	CreateCategory(data categories.CreateCategoryDto) error
	GetCategories() ([]categories.CategoryResponse, error)
}

func NewCategoriesService(env envs.Envs) *CategoriesService {
	return &CategoriesService{
		env: env,
	}
}

func (s *CategoriesService) CreateCategory(data categories.CreateCategoryDto) error {
	url := s.env.Get("CATEGORIES_URL")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return errors.NewInternalServerError("Error al convertir los datos a JSON: " + err.Error())
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return errors.NewInternalServerError("Error al llamar al servicio de categorías: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return errors.NewInternalServerError("Error al crear la categoría, código de respuesta: " + fmt.Sprint(resp.StatusCode))
	}
	return nil
}

func (s *CategoriesService) GetCategories() ([]categories.CategoryResponse, error) {
	url := s.env.Get("CATEGORIES_URL")
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.NewInternalServerError("Error al llamar al servicio de categorías: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewInternalServerError("Error al obtener las categorías, código de respuesta: " + fmt.Sprint(resp.StatusCode))
	}
	var categories []categories.CategoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, errors.NewInternalServerError("Error al decodificar la respuesta: " + err.Error())
	}
	return categories, nil
}
