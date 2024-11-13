package ratings

type RatingDTO struct {
	UserId      string `json:"user_id"`
	CourseId    string `json:"course_id"`
	RatingValue int    `json:"rating"`
}
