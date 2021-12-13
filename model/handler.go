package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dbHandler struct {
	db     *gorm.DB
	cache  map[string]string
	states map[string]map[string]interface{}
	
}

func newSqliteHandler(path string) (DBHandler, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Device{})

	return &dbHandler{
		db:     db,
		cache:  map[string]string{},
		states: map[string]map[string]interface{}{},
	}, nil
}

// func newPostgresqlHandler(path string) (DBHandler, error) {
// 	dsn := "host=localhost user=user password=user_password dbname=godopudb port=5432 sslmode=disable TimeZone=Asia/Seoul"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	db.AutoMigrate(&Device{})

// 	return &dbHandler{db: db, cache: map[string]string{}}, nil
// }

// func (s *dbHandler) GetDevice() *Device {
// 	device := &Device{}
// }
