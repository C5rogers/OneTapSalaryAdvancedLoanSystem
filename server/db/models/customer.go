package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name          string  `json:"name"`
	AccountNumber string  `gorm:"uniqueIndex" json:"account_number"`
	AccountType   string  `json:"account_type"`
	Balance       float64 `json:"balance"`
	Image         string  `json:"image,omitempty"`
	Verified      bool    `json:"verified"`
}
