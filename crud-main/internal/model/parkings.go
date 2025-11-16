package model

type Parkings struct {
	ID                 string `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CodeParking        string `gorm:"column:code_parking;type:varchar(50);NOT NULL" json:"codeParking"`
	IDSegment          string `gorm:"column:id_segment;type:uuid;NOT NULL" json:"iDSegment"`
	TypeParking        string `gorm:"column:type_parking;type:varchar(10);NOT NULL" json:"typeParking"`
	BoundaryParking    string `gorm:"column:boundary_parking;type:geography" json:"boundaryParking"`
	Position           string `gorm:"column:position;type:varchar(1)" json:"position"`
	StatusAvailability string `gorm:"column:status_availability;type:varchar(20)" json:"statusAvailability"`
	Description        string `gorm:"column:description;type:text" json:"description"`
	CreatedAt          int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
