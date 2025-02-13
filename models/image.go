package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid" json:"id" validate:"uuid"`
	Product uuid.UUID `gorm:"type:uuid" json:"product" validate:"uuid"`
	MD5     string    `json:"md5" validate:"required"`
	URL     string    `json:"url" validate:"required"`
}
