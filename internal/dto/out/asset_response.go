package out

import "time"

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if time.Time(d).IsZero() {
		return []byte(`null`), nil
	}
	formatted := time.Time(d).Format(`"2006-01-02"`)
	return []byte(formatted), nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = DateOnly(time.Time{})
		return nil
	}
	parsed, err := time.Parse(`"2006-01-02"`, string(data))
	if err != nil {
		return err
	}
	*d = DateOnly(parsed)
	return nil
}

type AssetResponse struct {
	ID              uint     `json:"asset_id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	CategoryName    string   `json:"category_name,omitempty"`
	StatusName      string   `json:"status_name,omitempty"`
	PurchaseDate    DateOnly `json:"purchase_date,omitempty"`
	Value           float64  `json:"value,omitempty"`
	MaintenanceDate DateOnly `json:"maintenance_date,omitempty"`
	MaintenanceCost float64  `json:"maintenance_cost,omitempty"`
}
