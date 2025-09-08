package models

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model

	ID              int64   `gorm:"primaryKey,autoIncrement" json:"id"`
	CustomerName    string  `json:"customerName"`
	Mobile          string  `json:"mobile"`
	AccountNumber   string  `gorm:"uniqueIndex" json:"accountNo"`
	BranchName      string  `json:"branchName"`
	ProductName     string  `json:"productName"`
	CustomerID      string  `gorm:"uniqueIndex" json:"customerId"`
	BranchCode      string  `json:"branchCode"`
	CustomerBalance float64 `json:"customerBalance"`
	Image           string  `json:"image,omitempty"`
	Verified        bool    `json:"verified"`
}

type SampleCustomer struct {
	Name          string `json:"customerName"`
	AccountNumber string `json:"accountNo"`
	Verified      bool   `json:"verified"`
}

type ValidationLog struct {
	RecordIndex      int       `json:"record_index"`
	Verified         bool      `json:"verified"`
	Errors           []string  `json:"errors"`
	NormalizedRecord *Customer `json:"normalized_record,omitempty"`
}
