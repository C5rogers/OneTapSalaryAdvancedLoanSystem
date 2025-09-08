package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	ID                  uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FromAccount         string    `json:"fromAccount"`
	ToAccount           string    `json:"toAccount"`
	Amount              float64   `json:"amount"`
	Remark              string    `json:"remark"`
	TransactionType     string    `json:"transactionType"`
	RequestID           string    `json:"requestId"`
	Reference           string    `json:"reference"`
	ThirdPartyReference string    `json:"thirdPartyReference"`
	InstitutionID       string    `json:"institutionId"`
	ClearedBalance      float64   `json:"clearedBalance"`
	TransactionDate     int64     `json:"transactionDate"`
	BillerID            string    `json:"billerId"`
	CustomerID          int64     `json:"customerId"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
