package out

type AssetCategoryResponse struct {
	AssetCategoryID uint   `json:"asset_category_id"`
	CategoryName    string `json:"category_name"`
	Description     string `json:"description"`
}
