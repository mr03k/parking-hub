package model

type VehicleCategories struct {
	ID           uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CodeCategory string `gorm:"column:code_category;type:varchar(20);NOT NULL" json:"codeCategory"`
	NameCategory string `gorm:"column:name_category;type:varchar(50);NOT NULL" json:"nameCategory"`
	Description  string `gorm:"column:description;type:text" json:"description"`
	CreatedAt    int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
