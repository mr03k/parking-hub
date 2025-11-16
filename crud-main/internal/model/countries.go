package model

type Countries struct {
	ID          uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CountryName string `gorm:"column:country_name;type:varchar(100);NOT NULL" json:"countryName"`
	CountryCode string `gorm:"column:country_code;type:varchar(3)" json:"countryCode"`
	IsoCode     string `gorm:"column:iso_code;type:varchar(2)" json:"isoCode"`
	Region      string `gorm:"column:region;type:varchar(50)" json:"region"`
	Capital     string `gorm:"column:capital;type:varchar(100)" json:"capital"`
	PhoneCode   string `gorm:"column:phone_code;type:varchar(10)" json:"phoneCode"`
	Currency    string `gorm:"column:currency;type:varchar(50)" json:"currency"`
	Population  int64  `gorm:"column:population;type:int8" json:"population"`
	Area        string `gorm:"column:area;type:float8" json:"area"`
	GeoBoundary string `gorm:"column:geo_boundary;type:geometry" json:"geoBoundary"`
	CreatedAt   int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
