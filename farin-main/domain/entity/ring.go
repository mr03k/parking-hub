package entity

type Ring struct {
	ID        int64   `json:"ID"`
	RingCode  string  `gorm:"type:varchar(100)" json:"ringCode"`
	Length    float64 `gorm:"type:double precision" json:"length"`
	RingName  string  `gorm:"type:varchar(120)" json:"ringName"`
	Geom      string  `gorm:"type:geometry" json:"geom"`
	ObjectId  int64   `json:"objectId"`
	ShapeLeng int64   `json:"shapeLeng"`
	ShapeArea int64   `json:"shapeArea"`

	Base
}
