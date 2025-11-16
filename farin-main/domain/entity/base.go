package entity

import "gorm.io/plugin/soft_delete"

type Base struct {
	ID        string                `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt int64                 `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt int64                 `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"-"`
}
