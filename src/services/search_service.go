package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/courses"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	decodedQuery, err := url.QueryUnescape(query)
	if err != nil {
		return nil, fmt.Errorf("error decodificando query: %v", err)
	}
	escapedQuery := url.QueryEscape(decodedQuery)
	url := fmt.Sprintf("%s?q=%s", searchURL, escapedQuery)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error del servidor: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, fmt.Errorf("error al parsear JSON: %v, body: %s", err, string(body))
	}
	return searchResponse.Courses, nil
}
