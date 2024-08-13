package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction model
type Transaction struct {
	ID        uint      `gorm:"primaryKey"`
	Amount    float64   `gorm:"not null"`
	Token     string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	Sender    string    `gorm:"not null"`
	Receiver  string    `gorm:"not null"`
	Type      string    `gorm:"not null"` // "deposit" or "withdrawal"
	TxHash    string    `gorm:"unique;not null"`
	AccountID uint      `gorm:"not null"`
}

// Save Create or update the transaction in the database
func (transaction *Transaction) Save(db *gorm.DB) error {
	return db.Save(transaction).Error
}
