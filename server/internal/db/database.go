package db

import (
	"fmt"
	"log"
	"path/filepath"

	"smartspa-admin/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a shared database handle after initialization.
var DB *gorm.DB

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

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

	// Create default admin user if not exists
	if err := createDefaultAdmin(database); err != nil {
		return nil, fmt.Errorf("create default admin: %w", err)
	}

	DB = database
	log.Printf("SQLite initialized at %s", absPath)
	return database, nil
}

func migrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&models.User{},
		&models.Member{},
		&models.Technician{},
		&models.ServiceItem{},
		&models.Appointment{},
		&models.Schedule{},
		&models.FissionLog{},
		&models.PhysicalProduct{},
		&models.InventoryLog{},
		&models.Order{},
	)
}

// createDefaultAdmin creates a default admin user if none exists
func createDefaultAdmin(database *gorm.DB) error {
	var count int64
	if err := database.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Only create default admin if no users exist
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash default password: %w", err)
		}

		admin := models.User{
			Username:     "admin",
			PasswordHash: string(hashedPassword),
			Role:         "manager",
			IsActive:     true,
		}

		if err := database.Create(&admin).Error; err != nil {
			return fmt.Errorf("create admin user: %w", err)
		}

		log.Println("Default admin user created (username: admin, password: Admin123!)")
	}

	return nil
}
