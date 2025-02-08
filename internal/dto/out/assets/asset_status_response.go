package assets

type AssetStatusResponse struct {
	AssetStatusID uint   `json:"asset_status_id,omitempty"`
	StatusName    string `json:"status_name,omitempty"`
	Description   string `json:"description,omitempty"`
}
