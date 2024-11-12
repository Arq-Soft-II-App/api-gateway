package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/ratings"
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
	GetAllRatings(ids []string) ([]ratings.RatingDTO, error)
}

func (s *RatingsService) NewRating(rating ratings.RatingDTO) (ratings.RatingDTO, error) {

}

func (s *RatingsService) UpdateRating(rating ratings.RatingDTO) (ratings.RatingDTO, error) {

}

func (s *RatingsService) GetAllRatings(ids []string) ([]ratings.RatingDTO, error) {

}

func (s *RatingsService) GetRatingByCourseId(courseId string) ([]ratings.RatingDTO, error) {

}
