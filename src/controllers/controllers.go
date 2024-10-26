package controllers

import "api-gateway/src/services"

type Controller struct {
	Auth  AuthControllerInterface
	Users UsersControllerInterface
	/* 	Courses      CoursesControllerInterface
	   	Inscriptions InscriptionsControllerInterface
	   	Search       SearchControllerInterface */
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		Auth:  NewAuthController(service.Auth),
		Users: NewUsersController(service.Users),
		/* 		Courses:      NewCoursesController(service.Courses),
		   		Inscriptions: NewInscriptionsController(service.Inscriptions),
		   		Search:       NewSearchController(service.Search), */
	}
}
