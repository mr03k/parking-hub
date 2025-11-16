package model

type Rings struct {
	ID                   uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	NameRing             string `gorm:"column:name_ring;type:varchar(100);NOT NULL" json:"nameRing"`
	CodeRing             string `gorm:"column:code_ring;type:varchar(10);NOT NULL" json:"codeRing"`
	LengthRing           string `gorm:"column:length_ring;type:numeric;NOT NULL" json:"lengthRing"`
	BoundaryRing         string `gorm:"column:boundary_ring;type:geography" json:"boundaryRing"`
	SpotsParking         int    `gorm:"column:spots_parking;type:int4" json:"spotsParking"`
	SpotsParkingDisabled int    `gorm:"column:spots_parking_disabled;type:int4" json:"spotsParkingDisabled"`
	SignsTraffic         int    `gorm:"column:signs_traffic;type:int4" json:"signsTraffic"`
	SignsTrafficDisabled int    `gorm:"column:signs_traffic_disabled;type:int4" json:"signsTrafficDisabled"`
	PointStart           string `gorm:"column:point_start;type:geography" json:"pointStart"`
	DistanceBuffer       string `gorm:"column:distance_buffer;type:numeric" json:"distanceBuffer"`
	CreatedAt            int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
	Description          string `gorm:"column:description;type:text" json:"description"`
}
