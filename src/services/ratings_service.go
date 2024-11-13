package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/ratings"
	"encoding/json"

	"api-gateway/src/errors"
	"bytes"
	"net/http"
)

type RatingsService struct {
	env envs.Envs
}

func NewRatingsService(env envs.Envs) *RatingsService {
	return &RatingsService{
		env: env,
	}
}

type RatingsServiceInterface interface {
	NewRating(rating ratings.RatingDTO) (ratings.RatingDTO, error)
	UpdateRating(rating ratings.RatingDTO) (ratings.RatingDTO, error)
	GetAllRatings() ([]ratings.RatingDTO, error)
}

func (s *RatingsService) NewRating(rating ratings.RatingDTO) (ratings.RatingDTO, error) {
	ratingsURL := s.env.Get("RATINGS_URL")

	jsonData, err := json.Marshal(rating)
	if err != nil {
		return ratings.RatingDTO{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar la valoracion", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", ratingsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return ratings.RatingDTO{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ratings.RatingDTO{}, errors.NewError("REQUEST_ERROR", "Error al crear la valoracion", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var createdRating ratings.RatingDTO
	if err := json.NewDecoder(resp.Body).Decode(&createdRating); err != nil {
		return ratings.RatingDTO{}, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	return createdRating, nil

}

func (s *RatingsService) UpdateRating(rating ratings.RatingDTO) (ratings.RatingDTO, error) {
	ratingsURL := s.env.Get("RATINGS_URL")

	jsonData, err := json.Marshal(rating)
	if err != nil {
		return ratings.RatingDTO{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar la valoracion", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("PUT", ratingsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return ratings.RatingDTO{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ratings.RatingDTO{}, errors.NewError("REQUEST_ERROR", "Error al actualizar la valoracion", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var updatedRating ratings.RatingDTO
	if err := json.NewDecoder(resp.Body).Decode(&updatedRating); err != nil {
		return ratings.RatingDTO{}, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	return updatedRating, nil

}

func (s *RatingsService) GetAllRatings() ([]ratings.RatingDTO, error) {
	ratingsURL := s.env.Get("RATINGS_URL")

	req, err := http.NewRequest("GET", ratingsURL, nil)
	if err != nil {
		return []ratings.RatingDTO{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []ratings.RatingDTO{}, errors.NewError("REQUEST_ERROR", "Error al obtener las valoraciones", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var ratingsList []ratings.RatingDTO
	if err := json.NewDecoder(resp.Body).Decode(&ratingsList); err != nil {
		return []ratings.RatingDTO{}, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	return ratingsList, nil
}
