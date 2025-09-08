package db

import (
	"math/rand"
	"time"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	"github.com/google/uuid"
)

func (db *Database) ProcessTransaction(customers []models.Customer, transactions []models.Transaction) (map[uint][]models.Transaction, error) {
	customerTxs := make(map[uint][]models.Transaction)

	for _, c := range customers {
		var txs []models.Transaction
		for _, t := range transactions {
			// match customer account number
			if t.FromAccount == c.AccountNumber {
				t.CustomerID = c.ID
				db.DB.Create(&t)
				txs = append(txs, t)
			}
		}

		if len(txs) == 0 {
			for i := 0; i < 3; i++ {
				tx := models.Transaction{
					CustomerID:          c.ID,
					FromAccount:         c.AccountNumber,
					ToAccount:           "ACCT-" + uuid.New().String(),
					Amount:              rand.Float64() * 500,
					Remark:              "Synthetic transaction",
					TransactionType:     "Derash Bill Payment",
					RequestID:           uuid.New().String(),
					Reference:           "REF-" + uuid.New().String(),
					ThirdPartyReference: "TPR-" + uuid.New().String(),
					InstitutionID:       "INST-" + uuid.New().String(),
					ClearedBalance:      c.CustomerBalance + (rand.Float64() * 500),
					TransactionDate:     time.Now().AddDate(0, 0, -(i * 10)).Unix(),
					BillerID:            "BILLER-" + uuid.New().String(),
				}
				db.DB.Create(&tx)
				txs = append(txs, tx)
			}
		}
		customerTxs[uint(c.ID)] = txs
	}
	return customerTxs, nil
}
