package server

import (
	"encoding/json"
	"net/http"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/payloads"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
)

func (s *Server) HandleProcessTransaction(w http.ResponseWriter, r *http.Request) error {

	customers, err := s.DB.GetAllCustomers()
	if err != nil {
		return utils.SendErrorResponse(w, "error retrieving customers", "database_error", http.StatusInternalServerError)
	}

	transactions, err := utils.LoadTransactions("data/transactions.json")
	if err != nil {
		return utils.SendErrorResponse(w, "error loading transactions", "file_error", http.StatusInternalServerError)
	}

	txMap, err := s.DB.ProcessTransaction(customers, transactions)
	if err != nil {
		return utils.SendErrorResponse(w, "error processing transactions", "database_error", http.StatusInternalServerError)
	}

	results := []payloads.Result{}

	for _, c := range customers {
		txs := txMap[uint(c.ID)]
		rating := utils.CalculateRating(txs, c.CustomerBalance)
		results = append(results, payloads.Result{Customer: c, Rating: rating})
	}

	responseData, _ := json.Marshal(payloads.ProcessTransactionsOutput{
		Ratings: results,
	})

	if _, err := w.Write(responseData); err != nil {
		return err
	}

	return nil
}
