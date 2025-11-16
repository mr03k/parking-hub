package model

type Forms struct {
	ID          uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	FormName    string `gorm:"column:form_name;type:varchar(100);NOT NULL" json:"formName"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	CreatedAt   int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
