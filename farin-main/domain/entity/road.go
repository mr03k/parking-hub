package entity

import "database/sql"

type Road struct {
	ID          int64           `json:"ID"`
	Geom        string          `gorm:"type:geometry;column:geom"`
	RoadName    string          `gorm:"type:varchar(100);not null;column:road_name"`
	RoadCode    int64           `gorm:"type:integer;not null;column:road_code"`
	Description sql.NullString  `gorm:"type:varchar(200);column:description"`
	Length      sql.NullFloat64 `gorm:"type:double precision;column:length"`
	SpeedLimit  sql.NullInt64   `gorm:"type:integer;column:speed_limit"`
	RoadType    string          `gorm:"type:varchar(100);not null;column:road_type"`
	RoadGrade   string          `gorm:"type:varchar(1);not null;column:road_grade"`

	Base
}
