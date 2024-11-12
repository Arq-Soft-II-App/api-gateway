package categories

type CreateCategoryDto struct {
	Category_Name string `json:"category_name"`
}

type CategoryResponse struct {
	ID            string `json:"category_id"`
	Category_Name string `json:"category_name"`
}
