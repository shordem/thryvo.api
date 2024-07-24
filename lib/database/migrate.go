package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var MigrationDir = "migrations"

type MigrationRecord struct {
	ID        uint `gorm:"primaryKey"`
	Filename  string
	AppliedAt time.Time
}

func Migrate(database DatabaseInterface) {
	// Read files in directory
	files, err := os.ReadDir(MigrationDir)

	// create migrations directory if not exists
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(MigrationDir, 0755); err != nil {
				fmt.Println("Failed to create migrations directory:", err)
			}
		} else {
			fmt.Println("Failed to read directory:", err)
			return
		}
	}

	// create migrations table if not exists
	if err := database.Connection().AutoMigrate(&MigrationRecord{}); err != nil {
		fmt.Println("Failed to create migrations_record table:", err)
		return
	}

	// Run migrations
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Get migration filename
		filename := file.Name()

		// Check if migration has already been applied
		var count int64
		database.Connection().Model(&MigrationRecord{}).Where("filename = ?", filename).Count(&count)
		if count > 0 {
			fmt.Println("Migration", filename, "has already been applied.")
			continue
		}

		// Read migration file content
		content, err := os.ReadFile(filepath.Join(MigrationDir, filename))
		if err != nil {
			fmt.Println("Failed to read migration file:", filename, err)
			return
		}

		// Split migration file content into individual SQL statements
		statements := strings.Split(string(content), ";")

		// Execute each SQL statement
		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}

			// Execute migration SQL
			if err := database.Connection().Exec(statement).Error; err != nil {
				fmt.Println("Failed to execute migration:", filename, err)
				return
			}
		}

		// Store migration record
		migrationRecord := MigrationRecord{
			Filename:  filename,
			AppliedAt: time.Now(),
		}

		if err := database.Connection().Create(&migrationRecord).Error; err != nil {
			fmt.Println("Failed to store migration record:", filename, err)
			return
		}

		fmt.Println("Migration", filename, "has been applied successfully.")
	}

	fmt.Println("All migrations have been applied.")

}
