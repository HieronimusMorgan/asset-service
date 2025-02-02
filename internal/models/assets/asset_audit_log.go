package assets

import "time"

type AssetAuditLog struct {
	LogID       uint      `gorm:"primaryKey" json:"log_id"`
	TableName   string    `gorm:"not null" json:"table_name"`
	Action      string    `gorm:"type:varchar(255);not null" json:"action"`
	OldData     *string   `gorm:"type:text" json:"old_data,omitempty"`
	NewData     *string   `gorm:"type:text" json:"new_data,omitempty"`
	PerformedAt time.Time `gorm:"autoCreateTime" json:"performed_at"`
	PerformedBy *string   `gorm:"type:varchar(255)" json:"performed_by,omitempty"`
}
