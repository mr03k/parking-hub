package entity

type Parking struct {
	ID          int64  `json:"ID"`
	ParkCode    int    `gorm:"type:integer" json:"parkCode"`
	ParkType    string `gorm:"type:varchar(1)" json:"parkType"`
	Position    string `gorm:"type:varchar(1)" json:"position"`
	Description string `gorm:"type:varchar(200)" json:"desc"`
	Geom        string `gorm:"type:geometry" json:"geom"`

	Base
}
