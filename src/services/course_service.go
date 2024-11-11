package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/courses"
	"api-gateway/src/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CourseService struct {
	CourseInterface CourseServiceInterface
	env             envs.Envs
}

type CourseServiceInterface interface {
	CreateCourse(data courses.CourseDTO) (courses.CourseDTO, error)
	UpdateCourse(data courses.CourseDTO) (courses.CourseDTO, error)
	GetCoursesList() ([]courses.CourseDTO, error)
	GetCourseById(id string) (courses.CourseDTO, error)
}

func NewCourseService(env envs.Envs) *CourseService {
	return &CourseService{
		env: env,
	}
}

func (s *CourseService) CreateCourse(data courses.CourseDTO) (courses.CourseDTO, error) {
	backendData := courses.CourseBackendDTO{
		BaseCourseDto: data.BaseCourseDto,
	}

	jsonData, err := json.Marshal(backendData)
	if err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al procesar el curso")
	}

	resp, err := http.Post(s.env.Get("COURSES_API_URL"), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al crear el curso")
	}
	defer resp.Body.Close()

	var backendResponse courses.CourseBackendDTO
	if err := json.NewDecoder(resp.Body).Decode(&backendResponse); err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al procesar la respuesta")
	}

	return courses.CourseDTO{
		ID:                 backendResponse.ID,
		BaseCourseDto:      backendResponse.BaseCourseDto,
		CourseCategoryName: backendResponse.CourseCategoryName,
		RatingAvg:          backendResponse.RatingAvg,
	}, nil
}

func (s *CourseService) UpdateCourse(data courses.CourseDTO) (courses.CourseDTO, error) {
	backendData := courses.CourseBackendDTO{
		ID:            data.ID,
		BaseCourseDto: data.BaseCourseDto,
	}

	jsonData, err := json.Marshal(backendData)
	if err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al procesar el curso")
	}

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/%s", s.env.Get("COURSES_API_URL"), data.ID),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al actualizar el curso")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al actualizar el curso")
	}
	defer resp.Body.Close()

	var backendResponse courses.CourseBackendDTO
	if err := json.NewDecoder(resp.Body).Decode(&backendResponse); err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al procesar la respuesta")
	}

	return courses.CourseDTO{
		ID:                 backendResponse.ID,
		BaseCourseDto:      backendResponse.BaseCourseDto,
		CourseCategoryName: backendResponse.CourseCategoryName,
		RatingAvg:          backendResponse.RatingAvg,
	}, nil
}

func (s *CourseService) GetCoursesList() ([]courses.CourseDTO, error) {
	resp, err := http.Get(s.env.Get("COURSES_API_URL"))
	if err != nil {
		return nil, errors.NewInternalServerError("Error al obtener los cursos")
	}
	defer resp.Body.Close()

	var backendCourses []courses.CourseBackendDTO
	if err := json.NewDecoder(resp.Body).Decode(&backendCourses); err != nil {
		return nil, errors.NewInternalServerError("Error al procesar la respuesta")
	}

	frontendCourses := make([]courses.CourseDTO, len(backendCourses))
	for i, course := range backendCourses {
		frontendCourses[i] = courses.CourseDTO{
			ID:                 course.ID,
			BaseCourseDto:      course.BaseCourseDto,
			CourseCategoryName: course.CourseCategoryName,
			RatingAvg:          course.RatingAvg,
		}
	}

	return frontendCourses, nil
}

func (s *CourseService) GetCourseById(id string) (courses.CourseDTO, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", s.env.Get("COURSES_API_URL"), id))
	if err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al obtener el curso")
	}
	defer resp.Body.Close()

	var backendCourse courses.CourseBackendDTO
	if err := json.NewDecoder(resp.Body).Decode(&backendCourse); err != nil {
		return courses.CourseDTO{}, errors.NewInternalServerError("Error al procesar la respuesta")
	}

	return courses.CourseDTO{
		ID:                 backendCourse.ID,
		BaseCourseDto:      backendCourse.BaseCourseDto,
		CourseCategoryName: backendCourse.CourseCategoryName,
		RatingAvg:          backendCourse.RatingAvg,
	}, nil
}
