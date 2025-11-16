package model

type Segments struct {
	ID                   uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	SegmentName          string `gorm:"column:segment_name;type:varchar(100);NOT NULL" json:"segmentName"`
	SegmentCode          string `gorm:"column:segment_code;type:varchar(10);NOT NULL" json:"segmentCode"`
	PartID               string `gorm:"column:part_id;type:uuid;NOT NULL" json:"partID"`
	RoadID               string `gorm:"column:road_id;type:uuid;NOT NULL" json:"roadID"`
	DistrictID           string `gorm:"column:district_id;type:uuid;NOT NULL" json:"districtID"`
	SegmentLength        string `gorm:"column:segment_length;type:numeric" json:"segmentLength"`
	SegmentBoundary      string `gorm:"column:segment_boundary;type:geometry" json:"segmentBoundary"`
	ParkingSpots         int    `gorm:"column:parking_spots;type:int4" json:"parkingSpots"`
	DisabledParkingSpots int    `gorm:"column:disabled_parking_spots;type:int4" json:"disabledParkingSpots"`
	Description          string `gorm:"column:description;type:text" json:"description"`
	CreatedAt            int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
