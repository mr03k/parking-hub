package entity

type Segment struct {
	ID          int64   `json:"ID"`
	SegCode     string  `gorm:"type:varchar(100)" json:"segCode"`
	SegLength   float64 `gorm:"type:double precision" json:"segLength"`
	SegName     string  `gorm:"type:varchar(200)" json:"segName"`
	Description string  `gorm:"type:varchar(200)" json:"desc"`
	Geom        string  `gorm:"type:geometry" json:"geom"`
	ObjectID    int64   `json:"objectID"`
	Junction    int8    `json:"junction"`
	ShapeLeng   float64 `json:"shapeLeng"`
	ShapeArea   float64 `json:"shapeArea"`

	Base
}
