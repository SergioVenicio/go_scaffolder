package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var validate validator.Validate

func init() {
	validate = *validator.New()
}

type Product struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid" json:"id" validate:"uuid"`
	Description string    `json:"description" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Status      string    `json:"status"`
	Stock       int       `json:"stock" validate:"min=0"`
	Images      []Image   `gorm:"foreignKey:Product" json:"images" validate:"required"`
}

func (p *Product) NewID() {
	p.ID = uuid.New()
}

func (p *Product) Validate() error {
	return validate.Struct(p)
}
