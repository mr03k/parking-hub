package model

type Modules struct {
	ID          uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	ModuleName  string `gorm:"column:module_name;type:varchar(100);NOT NULL" json:"moduleName"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	CreatedAt   int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
