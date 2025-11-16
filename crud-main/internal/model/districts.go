package model

type Districts struct {
	ID           uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	DistrictName string `gorm:"column:district_name;type:varchar(100);NOT NULL" json:"districtName"`
	CodeDistrict string `gorm:"column:code_district;type:varchar(10)" json:"codeDistrict"`
	IDCity       string `gorm:"column:id_city;type:uuid;NOT NULL" json:"iDCity"`
	BoundaryGeo  string `gorm:"column:boundary_geo;type:geography" json:"boundaryGeo"`
	Population   int64  `gorm:"column:population;type:int8" json:"population"`
	Area         string `gorm:"column:area;type:float8" json:"area"`
	CreatedAt    int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
