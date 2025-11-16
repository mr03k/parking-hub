package model

type Cities struct {
	ID        uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CityName  string `gorm:"column:city_name;type:varchar(100);NOT NULL" json:"cityName"`
	CodeCity  string `gorm:"column:code_city;type:varchar(3)" json:"codeCity"`
	IDCountry string `gorm:"column:id_country;type:uuid;NOT NULL" json:"iDCountry"`
	Boundary  string `gorm:"column:boundary;type:geography" json:"boundary"`
	CreatedAt int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
