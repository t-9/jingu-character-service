package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

// OpenGorm open the Character DB By Gorm.
func OpenGorm() (*gorm.DB, error) {
	return gorm.Open(
		os.Getenv("DB_DRIVER"),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
}
