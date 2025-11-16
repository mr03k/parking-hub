package model

type Roads struct {
	ID                   uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	RoadName             string `gorm:"column:road_name;type:varchar(100);NOT NULL" json:"roadName"`
	CodeRoad             string `gorm:"column:code_road;type:varchar(10);NOT NULL" json:"codeRoad"`
	TypeRoad             string `gorm:"column:type_road;type:varchar(50)" json:"typeRoad"`
	GradeRoad            string `gorm:"column:grade_road;type:varchar(1)" json:"gradeRoad"`
	LengthRoad           string `gorm:"column:length_road;type:numeric" json:"lengthRoad"`
	WidthRoad            string `gorm:"column:width_road;type:numeric" json:"widthRoad"`
	LimitSpeed           int    `gorm:"column:limit_speed;type:int4" json:"limitSpeed"`
	BoundaryRoad         string `gorm:"column:boundary_road;type:geography" json:"boundaryRoad"`
	SpotsParking         int    `gorm:"column:spots_parking;type:int4" json:"spotsParking"`
	SpotsParkingDisabled int    `gorm:"column:spots_parking_disabled;type:int4" json:"spotsParkingDisabled"`
	CreatedAt            int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
	Description          string `gorm:"column:description;type:text" json:"description"`
}
