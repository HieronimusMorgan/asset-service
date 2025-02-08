package assets

type AssetStatusRequest struct {
	StatusName  string `json:"status_name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
