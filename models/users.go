package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID    uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
	City  *string `json:"city"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
