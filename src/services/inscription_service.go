package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/courses"
	"api-gateway/src/dto/inscriptions"
	"api-gateway/src/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type InscriptionService struct {
	InscriptionInterface InscriptionServiceInterface
	env                  envs.Envs
	courseService        CourseServiceInterface
	usersService         UsersServiceInterface
}

type InscriptionServiceInterface interface {
	CreateInscription(data inscriptions.EnrollRequestResponseDto) (inscriptions.EnrollRequestResponseDto, error)
	GetMyCourses(userId string) ([]courses.CourseListDto, error)
	GetCourseStudents(courseId string) (inscriptions.StudentsInCourse, error)
	IsEnrolled(courseId string, userId string) (bool, error)
}

func NewInscriptionsService(env envs.Envs, courseService CourseServiceInterface, usersService UsersServiceInterface) *InscriptionService {
	return &InscriptionService{
		env:           env,
		courseService: courseService,
		usersService:  usersService,
	}
}

func (s *InscriptionService) CreateInscription(data inscriptions.EnrollRequestResponseDto) (inscriptions.EnrollRequestResponseDto, error) {
	inscriptionsURL := fmt.Sprintf("%senroll", s.env.Get("INSCRIPTIONS_URL"))

	jsonData, err := json.Marshal(data)
	if err != nil {
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar la inscripción", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", inscriptionsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("REQUEST_ERROR", "Error al crear la inscripción", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var createdInscription inscriptions.EnrollRequestResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&createdInscription); err != nil {
		fmt.Printf("Error al decodificar la respuesta: %v\n", err)
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	return createdInscription, nil
}

func (s *InscriptionService) GetMyCourses(userId string) ([]courses.CourseListDto, error) {
	fmt.Println("userId", userId)
	inscriptionsURL := fmt.Sprintf("%smyCourses?userId=%s", s.env.Get("INSCRIPTIONS_URL"), userId)

	resp, err := http.Get(inscriptionsURL)
	fmt.Println("resp inscriptions", resp)
	if err != nil {
		return nil, errors.NewError("REQUEST_ERROR", "Error al obtener las inscripciones", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// Manejar el caso de 404
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NewError("NOT_FOUND", "No hay inscripciones para este usuario", http.StatusNotFound)
	}

	// Decodificar la respuesta
	var inscriptionsData []inscriptions.Course
	if err := json.NewDecoder(resp.Body).Decode(&inscriptionsData); err != nil {
		fmt.Printf("Error al decodificar la respuesta: %v\n", err)
		return nil, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	// Extraer los IDs de los cursos
	var courseIDs []string
	for _, inscription := range inscriptionsData {
		courseIDs = append(courseIDs, inscription.CourseId)
	}

	// Llamar al servicio de cursos para obtener la información detalladaError al obtener los cursos
	coursesList, err := s.courseService.GetCoursesList(courseIDs)
	if err != nil {
		fmt.Printf("Error al obtener los cursos: %v\n", err)
		return nil, errors.NewError("COURSE_SERVICE_ERROR", "Error al obtener los cursos", http.StatusInternalServerError)
	}

	return coursesList, nil
}

func (s *InscriptionService) GetCourseStudents(courseId string) (inscriptions.StudentsInCourse, error) {
	inscriptionsURL := fmt.Sprintf("%sstudentsInThisCourse/%s", s.env.Get("INSCRIPTIONS_URL"), courseId)

	req, err := http.NewRequest("GET", inscriptionsURL, nil)
	if err != nil {
		return nil, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud HTTP", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewError("REQUEST_ERROR", "Error al obtener los estudiantes", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var students inscriptions.StudentsInCourse
	if err := json.NewDecoder(resp.Body).Decode(&students); err != nil {
		return nil, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	userIds := make([]string, 0)
	for _, student := range students {
		userIds = append(userIds, student.UserId)
	}

	// Obtener información de usuarios
	users, err := s.usersService.GetUsersList(userIds)
	if err != nil {
		return nil, err
	}

	// Construir respuesta final
	response := make(inscriptions.StudentsInCourse, len(users))
	for i, user := range users {
		response[i] = inscriptions.Student{
			UserId:   user.ID,
			UserName: user.Name + " " + user.Lastname,
			Avatar:   user.Avatar,
		}
	}

	return response, nil
}

func (s *InscriptionService) IsEnrolled(courseId string, userId string) (bool, error) {
	inscriptionsURL := fmt.Sprintf("%s/isEnrolled/%s/%s", s.env.Get("INSCRIPTIONS_URL"), courseId, userId)

	req, err := http.NewRequest("GET", inscriptionsURL, nil)
	if err != nil {
		return false, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.NewError("REQUEST_ERROR", "Error al obtener la inscripción", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var response struct {
		Enrolled bool `json:"enrolled"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error al decodificar la respuesta:", err)
		return false, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	return response.Enrolled, nil
}
