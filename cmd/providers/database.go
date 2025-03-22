package providers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DataBase *gorm.DB
}

func NewDatabase(dsn string) (*Database, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{
		DataBase: db,
	}, nil
}
