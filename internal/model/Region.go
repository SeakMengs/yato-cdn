package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	ID     string `gorm:"type:text;primaryKey" json:"id"`
	Name   string `gorm:"type:text;not null;" json:"name" form:"name" binding:"required"`
	Domain string `gorm:"type:text;not null;" json:"domain" form:"domain" binding:"required"`

	BaseModel
}

func (r Region) TableName() string {
	return "regions"
}

func (r *Region) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	r.ID = uuid.NewString()
	return
}
