package model

type BaseRates struct {
	ID                uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	IDCategoryVehicle string `gorm:"column:id_category_vehicle;type:uuid;NOT NULL" json:"iDCategoryVehicle"`
	FromMinutes       int    `gorm:"column:from_minutes;type:int4;NOT NULL" json:"fromMinutes"`
	ToMinutes         int    `gorm:"column:to_minutes;type:int4;NOT NULL" json:"toMinutes"`
	BaseRate          string `gorm:"column:base_rate;type:numeric;NOT NULL" json:"baseRate"`
	Description       string `gorm:"column:description;type:text" json:"description"`
	CreatedAt         int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
