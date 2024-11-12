package ratings

type RatingDTO struct {
	Id          string `json:"id"`
	UserId      string `json:"user_id"`
	CourseId    string `json:"course_id"`
	RatingValue int    `json:"rating_value"`
}
