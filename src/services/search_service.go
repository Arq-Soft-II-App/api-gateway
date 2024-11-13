package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/courses"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SearchService struct {
	env envs.Envs
}

func NewSearchService(env envs.Envs) *SearchService {
	return &SearchService{
		env: env,
	}
}

type SearchServiceInterface interface {
	SearchCourses(query string) ([]courses.CourseListDto, error)
}

type SearchResponse struct {
	Courses []courses.CourseListDto `json:"courses"`
}

func (s *SearchService) SearchCourses(query string) ([]courses.CourseListDto, error) {
	searchURL := s.env.Get("SEARCH_API_URL")
	url := fmt.Sprintf("%s?q=%s", searchURL, query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, err
	}
	return searchResponse.Courses, nil
}
