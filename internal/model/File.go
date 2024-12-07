package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID   string `gorm:"type:text;primaryKey" json:"id"`
	Name string `gorm:"type:text;not null;" json:"name" form:"name" binding:"required"`

	BaseModel
}

func (f File) TableName() string {
	return "files"
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	f.ID = uuid.NewString()
	return
}
