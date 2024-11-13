package courses

type BaseCourseDto struct {
	CourseName        string  `json:"course_name"`
	CourseDescription string  `json:"description"`
	CoursePrice       float64 `json:"price"`
	CourseDuration    int     `json:"duration"`
	CourseCapacity    int     `json:"capacity"`
	CategoryID        string  `json:"category_id"`
	CourseInitDate    string  `json:"init_date"`
	CourseState       *bool   `json:"state"`
	CourseImage       string  `json:"image"`
}

type CourseDTO struct {
	ID string `json:"id"`
	BaseCourseDto
	CourseCategoryName string  `json:"category_name"`
	RatingAvg          float64 `json:"ratingavg"`
}

type CourseBackendDTO struct {
	ID string `json:"_id"`
	BaseCourseDto
	CourseCategoryName string  `json:"category_name"`
	RatingAvg          float64 `json:"ratingavg"`
}

type CourseComment struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

type CourseListDto struct {
	Id           string  `json:"id"`
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

type CourseCommentsResponse []CourseComment
