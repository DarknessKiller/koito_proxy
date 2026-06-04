package database

import (
	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
	"log"

	"github.com/libtnb/sqlite"
	"gorm.io/gorm"
)

func InitiateDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := OpenSQLiteDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
		return nil, err
	}

	if err := db.AutoMigrate(
		model.Rule{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
		return nil, err
	}

	return db, nil
}

func OpenSQLiteDatabase(cfg *config.Config) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
		return nil, err
	}

	log.Println("SQLite connected! File:", cfg.DBPath)
	return db, nil
}
