package model

type Parts struct {
	ID                   uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	PartName             string `gorm:"column:part_name;type:varchar(100);NOT NULL" json:"partName"`
	CodePart             string `gorm:"column:code_part;type:varchar(10);NOT NULL" json:"codePart"`
	IDRoad               string `gorm:"column:id_road;type:uuid;NOT NULL" json:"iDRoad"`
	LengthPart           string `gorm:"column:length_part;type:numeric" json:"lengthPart"`
	BoundaryPart         string `gorm:"column:boundary_part;type:geography" json:"boundaryPart"`
	SpotsParking         int    `gorm:"column:spots_parking;type:int4" json:"spotsParking"`
	SpotsParkingDisabled int    `gorm:"column:spots_parking_disabled;type:int4" json:"spotsParkingDisabled"`
	Description          string `gorm:"column:description;type:text" json:"description"`
}
