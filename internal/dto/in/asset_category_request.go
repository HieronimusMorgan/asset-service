package in

type AssetCategoryRequest struct {
	CategoryName string `json:"category_name" validate:"required"`
	Description  string `json:"description" validate:"optional"`
}
