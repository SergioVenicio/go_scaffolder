package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/SergioVenicio/go_scaffolder/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresql struct{}

var (
	once sync.Once
	db   *gorm.DB
)

func NewPostgresql(conf *config.Config) Database {
	var dbInstance postgresql
	once.Do(func() {
		if db == nil {
			dsn := fmt.Sprintf(conf.DatabaseURI)
			dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				panic("failed to connect database")
			}

			sqlDB, err := dbConn.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(50)
			sqlDB.SetConnMaxLifetime(time.Hour)
			db = dbConn
		}
		dbInstance = postgresql{}
	})

	return &dbInstance
}

func (p *postgresql) GetDb() *gorm.DB {
	return db
}
