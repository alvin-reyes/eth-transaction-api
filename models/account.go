package models

import (
	"gorm.io/gorm"
)

// Account model
type Account struct {
	ID           uint          `gorm:"primaryKey"`
	AccountUuid  string        `gorm:"unique;not null"`
	Address      string        `gorm:"unique;not null"`
	Transactions []Transaction `gorm:"foreignKey:AccountID"`
}

// Save Create or update the account in the database
func (account *Account) Save(db *gorm.DB) error {
	return db.Save(account).Error
}
