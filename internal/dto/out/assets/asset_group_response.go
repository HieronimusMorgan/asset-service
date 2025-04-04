package assets

type AssetGroupResponse struct {
	AssetGroupID uint   `gorm:"primaryKey;column:asset_group_id"  json:"asset_group_id,omitempty"`
	GroupName    string `gorm:"type:varchar(100);not null"  json:"group_name,omitempty"`
	Description  string `gorm:"type:text" json:"description,omitempty"`
}

type AssetGroupDetailResponse struct {
	AssetGroupID uint                       `gorm:"primaryKey;column:asset_group_id"  json:"asset_group_id,omitempty"`
	GroupName    string                     `gorm:"type:varchar(100);not null"  json:"group_name,omitempty"`
	Description  string                     `gorm:"type:text" json:"description,omitempty"`
	Member       []AssetGroupMemberResponse `json:"member,omitempty"`
}

type AssetGroupMemberResponse struct {
	UserID         uint                                     `json:"user_id"`
	Username       string                                   `json:"username"`
	FullName       string                                   `json:"full_name"`
	ProfilePicture string                                   `json:"profile_picture"`
	Permission     []AssetGroupMemberWithPermissionResponse `json:"permission"`
}

type AssetGroupMemberWithPermissionResponse struct {
	PermissionID   *uint   `json:"permission_id"`
	PermissionName *string `json:"permission_name"`
}
