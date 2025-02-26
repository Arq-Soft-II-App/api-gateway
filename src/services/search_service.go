package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/courses"
	"api-gateway/src/errors"
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

func (s *SearchService) SearchCourses(query string) ([]courses.CourseListDto, error) {
	fmt.Printf("Service: Searching courses using %s URL...\n", s.env.Get("SEARCH_API_URL"))
	searchURL := s.env.Get("SEARCH_API_URL")
	decodedQuery, err := url.QueryUnescape(query)
	if err != nil {
		return nil, errors.NewError("DECODING_ERROR", fmt.Sprintf("Error decodificando query: %v", err), http.StatusBadRequest)
	}
	escapedQuery := url.QueryEscape(decodedQuery)
	fullURL := fmt.Sprintf("%s?q=%s", searchURL, escapedQuery)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, errors.NewError("REQUEST_ERROR", fmt.Sprintf("Error al realizar la solicitud: %v", err), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []courses.CourseListDto{}, errors.NewError("NOT_FOUND", "No se encontraron cursos", http.StatusNotFound)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewError("SERVER_ERROR", fmt.Sprintf("Error del servidor: %d", resp.StatusCode), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewError("READ_ERROR", fmt.Sprintf("Error al leer la respuesta: %v", err), http.StatusInternalServerError)
	}

	type tempCourse struct {
		ID           string  `json:"id"`
		CategoryID   string  `json:"category_id"`
		CourseName   string  `json:"course_name"`
		Description  string  `json:"description"`
		Price        float64 `json:"price"`
		Duration     int     `json:"duration"`
		Capacity     int     `json:"capacity"`
		InitDate     string  `json:"init_date"`
		State        bool    `json:"state"`
		Image        string  `json:"image"`
		CategoryName string  `json:"category_name"`
		RatingAvg    float64 `json:"ratingavg"`
	}

	type tempSearchResponse struct {
		Courses []tempCourse `json:"courses"`
	}

	var tempResponse tempSearchResponse
	err = json.Unmarshal(body, &tempResponse)
	if err != nil {
		return nil, errors.NewError("JSON_PARSE_ERROR", fmt.Sprintf("Error al parsear JSON: %v", err), http.StatusInternalServerError)
	}

	if len(tempResponse.Courses) == 0 {
		return []courses.CourseListDto{}, errors.NewError("NOT_FOUND", "No se encontraron cursos", http.StatusNotFound)
	}

	coursesList := make([]courses.CourseListDto, len(tempResponse.Courses))
	for i, tempCourse := range tempResponse.Courses {
		coursesList[i] = courses.CourseListDto{
			Id:           tempCourse.ID,
			CategoryID:   tempCourse.CategoryID,
			CourseName:   tempCourse.CourseName,
			Description:  tempCourse.Description,
			Price:        tempCourse.Price,
			Duration:     tempCourse.Duration,
			Capacity:     tempCourse.Capacity,
			InitDate:     tempCourse.InitDate,
			State:        tempCourse.State,
			Image:        tempCourse.Image,
			CategoryName: tempCourse.CategoryName,
			RatingAvg:    tempCourse.RatingAvg,
		}
	}

	return coursesList, nil
}
