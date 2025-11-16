package model

type RoadCategories struct {
	ID               uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CodeCategoryRoad string `gorm:"column:code_category_road;type:varchar(10);NOT NULL" json:"codeCategoryRoad"`
	NameCategoryRoad string `gorm:"column:name_category_road;type:varchar(50);NOT NULL" json:"nameCategoryRoad"`
	Description      string `gorm:"column:description;type:text" json:"description"`
	CreatedAt        int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
