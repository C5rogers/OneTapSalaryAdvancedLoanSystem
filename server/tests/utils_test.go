package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/payloads"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
)

func TestErrorResponse(t *testing.T) {
	err := utils.CoalesceError()
	if err != nil {
		t.Error("Expected nil, got error:", err)
	}
}

func TestGetDataFilePath(t *testing.T) {
	filePath := utils.GetDataFilePath("data/customers.json")
	if filePath == "" {
		t.Error("Expected non-empty file path")
	}
}

func TestLoadCustomers(t *testing.T) {
	// Sample customer data
	sampleCustomers := []models.Customer{
		{
			CustomerName:    "Alice Johnson",
			Mobile:          "0912345678",
			AccountNumber:   "123456789",
			BranchName:      "Main Branch",
			ProductName:     "Savings",
			CustomerID:      "C001",
			BranchCode:      "01",
			CustomerBalance: 5000.0,
			Verified:        true,
		},
		{
			CustomerName:    "Bob Smith",
			Mobile:          "0922334455",
			AccountNumber:   "987654321",
			BranchName:      "Sub Branch",
			ProductName:     "Current",
			CustomerID:      "C002",
			BranchCode:      "02",
			CustomerBalance: 1200.0,
			Verified:        false,
		},
	}

	// Create a temporary JSON file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "customers.json")

	fileContent, err := json.Marshal(sampleCustomers)
	if err != nil {
		t.Fatalf("failed to marshal sample customers: %v", err)
	}

	if err := os.WriteFile(tmpFile, fileContent, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Call the function under test
	customers, err := utils.LoadCustomers(tmpFile)
	if err != nil {
		t.Fatalf("LoadCustomers failed: %v", err)
	}

	// Assertions
	if len(customers) != len(sampleCustomers) {
		t.Errorf("expected %d customers, got %d", len(sampleCustomers), len(customers))
	}

	for i, c := range customers {
		if c.CustomerName != sampleCustomers[i].CustomerName {
			t.Errorf("expected name %s, got %s", sampleCustomers[i].CustomerName, c.CustomerName)
		}
		if c.AccountNumber != sampleCustomers[i].AccountNumber {
			t.Errorf("expected account %s, got %s", sampleCustomers[i].AccountNumber, c.AccountNumber)
		}
		if c.CustomerBalance != sampleCustomers[i].CustomerBalance {
			t.Errorf("expected balance %.2f, got %.2f", sampleCustomers[i].CustomerBalance, c.CustomerBalance)
		}
	}
}

func TestLoadTransactions(t *testing.T) {
	// Sample transaction payloads
	sampleTransactions := []payloads.TransactionPayload{
		{
			FromAccount:     "111111",
			ToAccount:       "222222",
			Amount:          "200.50",
			Remark:          "Test payment",
			TransactionType: "mpesa Transaction",
			RequestID:       "REQ123",
			Reference:       "REF123",
			TransactionDate: "1731867679364", // ms since epoch as string
			ClearedBalance:  "102248.58",
		},
		{
			FromAccount:     "333333",
			ToAccount:       "444444",
			Amount:          "500.00",
			Remark:          "Deposit",
			TransactionType: "bank Transfer",
			RequestID:       "REQ456",
			Reference:       "REF456",
			TransactionDate: "1731867679365",
			ClearedBalance:  "50000.00",
		},
	}

	// Create a temporary JSON file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "transactions.json")

	fileContent, err := json.Marshal(sampleTransactions)
	if err != nil {
		t.Fatalf("failed to marshal sample transactions: %v", err)
	}

	if err := os.WriteFile(tmpFile, fileContent, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Call the function under test
	transactions, err := utils.LoadTransactions(tmpFile)
	if err != nil {
		t.Fatalf("LoadTransactions failed: %v", err)
	}

	// Assertions
	if len(transactions) != len(sampleTransactions) {
		t.Errorf("expected %d transactions, got %d", len(sampleTransactions), len(transactions))
	}

	// Verify first transaction fields
	first := transactions[0]
	if first.FromAccount != "111111" {
		t.Errorf("expected FromAccount %s, got %s", "111111", first.FromAccount)
	}
	if first.ToAccount != "222222" {
		t.Errorf("expected ToAccount %s, got %s", "222222", first.ToAccount)
	}
	if first.Amount != 200.50 {
		t.Errorf("expected Amount %.2f, got %.2f", 200.50, first.Amount)
	}
	if first.ClearedBalance != 102248.58 {
		t.Errorf("expected ClearedBalance %.2f, got %.2f", 102248.58, first.ClearedBalance)
	}
	if first.Reference != "REF123" {
		t.Errorf("expected Reference %s, got %s", "REF123", first.Reference)
	}
}
