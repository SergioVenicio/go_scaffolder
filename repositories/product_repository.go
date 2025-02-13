package repositories

import (
	"github.com/SergioVenicio/go_scaffolder/database"
	"github.com/SergioVenicio/go_scaffolder/models"
	"github.com/sirupsen/logrus"
)

type ProductRepository struct {
	logger *logrus.Logger
	db     database.Database
}

func NewProductRepository(db database.Database, logger *logrus.Logger) Repository[models.Product] {
	return &ProductRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ProductRepository) Insert(p models.Product) error {
	r.logger.Infof("inserting new product %s", p.ID.String())
	tx := r.db.GetDb().Create(&p)
	if tx.Error != nil {
		r.logger.WithError(tx.Error).Error("product insert error")
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}
