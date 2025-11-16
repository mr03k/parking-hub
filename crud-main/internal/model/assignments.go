package model

type Assignments struct {
	ID            uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	IDUser        string `gorm:"column:id_user;type:uuid;NOT NULL" json:"iDUser"`
	IDRole        string `gorm:"column:id_role;type:uuid;NOT NULL" json:"iDRole"`
	IDModule      string `gorm:"column:id_module;type:uuid;NOT NULL" json:"iDModule"`
	IDForm        string `gorm:"column:id_form;type:uuid;NOT NULL" json:"iDForm"`
	CreatedAt     int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
	AccessEndDate int    `gorm:"column:access_end_date;type:int4;NOT NULL" json:"accessEndDate"`
}
