package controllers

import "api-gateway/src/services"

type Controller struct {
	Auth         AuthControllerInterface
	Users        UsersControllerInterface
	Courses      CourseControllerInterface
	Comments     CommentsControllerInterface
	Ratings      RatingsControllerInterface
	Categories   CategoriesControllerInterface
	Inscriptions InscriptionControllerInterface
	Search       SearchControllerInterface
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		Auth:         NewAuthController(service.Auth),
		Users:        NewUsersController(service.Users),
		Courses:      NewCourseController(service.Courses),
		Comments:     NewCommentsController(service.Comments),
		Ratings:      NewRatingsController(service.Ratings),
		Categories:   NewCategoriesController(service.Categories),
		Inscriptions: NewInscriptionsController(service.Inscriptions),
		Search:       NewSearchController(service.Search),
	}
}
