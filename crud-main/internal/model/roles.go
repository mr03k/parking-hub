package model

type Roles struct {
	ID          string `gorm:"column:id;type:uuid;primary_key" json:"id"`
	RoleName    string `gorm:"column:role_name;type:varchar(50);NOT NULL" json:"roleName"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	CreatedAt   int64  `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
