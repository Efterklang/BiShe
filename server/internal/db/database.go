package db

import (
	"fmt"
	"log"
	"path/filepath"

	"smartspa-admin/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a shared database handle after initialization.
var DB *gorm.DB

// Init opens the SQLite database and runs schema migrations.
func Init(dbPath string) (*gorm.DB, error) {
	if dbPath == "" {
		dbPath = "spa_management.db"
	}

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("resolve db path: %w", err)
	}

	database, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := migrate(database); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	DB = database
	log.Printf("SQLite initialized at %s", absPath)
	return database, nil
}

func migrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&models.Member{},
		&models.Technician{},
		&models.ServiceItem{},
		&models.Appointment{},
		&models.Schedule{},
		&models.FissionLog{},
	)
}
