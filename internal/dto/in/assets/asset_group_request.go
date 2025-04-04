package assets

type AssetGroupRequest struct {
	GroupName   string `json:"asset_group_name" validate:"required"`
	Description string `json:"description" validate:"optional"`
	OwnerUserID uint   `json:"owner_user_id" validate:"required"`
}

type AssetGroupNameRequest struct {
	GroupID     string `json:"asset_group_id" validate:"required"`
	GroupName   string `json:"asset_group_name" validate:"required"`
	Description string `json:"description" validate:"optional"`
	OwnerUserID uint   `json:"owner_user_id" validate:"required"`
}

type AssetGroupPermissionRequest struct {
	PermissionName string `json:"permission_name" validate:"required"`
	Description    string `json:"description" validate:"optional"`
}

type AssetGroupAssetRequest struct {
	AssetID      uint `json:"asset_id" validate:"required"`
	AssetGroupID uint `json:"asset_group_id" validate:"required"`
}

type AssetGroupMemberRequest struct {
	UserID       uint `json:"user_id" validate:"required"`
	AssetGroupID uint `json:"asset_group_id" validate:"required"`
}
