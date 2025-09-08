package utils

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/payloads"
)

func GetDataFilePath(fileName string) string {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	// Resolve the data path relative to the current working directory
	return filepath.Join(dir, "data", fileName)
}

func LoadCustomers(path string) ([]models.Customer, error) {
	customerFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var customers []models.Customer
	if err := json.Unmarshal(customerFile, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}

func LoadTransactions(path string) ([]models.Transaction, error) {
	transactionFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rawTransactions []payloads.TransactionPayload

	if err := json.Unmarshal(transactionFile, &rawTransactions); err != nil {
		return nil, err
	}
	var transactions []models.Transaction

	for _, rt := range rawTransactions {
		amount, _ := strconv.ParseFloat(rt.Amount, 64)
		clearedBalance, _ := strconv.ParseFloat(rt.ClearedBalance, 64)
		transactionDate, _ := strconv.Atoi(rt.TransactionDate)
		tx := models.Transaction{
			FromAccount:         rt.FromAccount,
			ToAccount:           rt.ToAccount,
			Amount:              amount,
			Remark:              rt.Remark,
			TransactionType:     rt.TransactionType,
			RequestID:           rt.RequestID,
			Reference:           rt.Reference,
			ThirdPartyReference: rt.ThirdPartyReference,
			InstitutionID: func() string {
				if rt.InstitutionID != nil {
					return *rt.InstitutionID
				}
				return ""
			}(),
			ClearedBalance:  clearedBalance,
			TransactionDate: int64(transactionDate),
			BillerID: func() string {
				if rt.BillerID != nil {
					return *rt.BillerID
				}
				return ""
			}(),
		}
		transactions = append(transactions, tx)

	}
	return transactions, nil
}

func LoadSampleCustomers(path string) ([]payloads.SampleCustomer, error) {
	ext := strings.ToLower(filepath.Ext(path))

	filePath := GetDataFilePath(path)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var customers []payloads.SampleCustomer
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &customers); err != nil {
			return nil, err
		}
	case ".csv":
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		for i, row := range records {
			if i == 0 {
				continue
			}
			// if len(row) < 4 {
			// 	continue
			// }
			c := payloads.SampleCustomer{
				CustomerName: row[0],
				AccountNo:    row[1],
				Verified:     false,
			}
			customers = append(customers, c)
		}
	}

	return customers, nil
}

func ValidateCustomers(customerPath string, sampleCustomerPath string) ([]models.ValidationLog, []models.Customer, error) {
	customers, err := LoadCustomers(customerPath)
	if err != nil {
		return nil, nil, err
	}

	sampleCustomers, err := LoadSampleCustomers(sampleCustomerPath)
	if err != nil {
		return nil, nil, err
	}

	logs := []models.ValidationLog{}
	validatedCustomers := []models.Customer{}

	accNumRegex := regexp.MustCompile(`^\d+$`)

	for i, sample := range sampleCustomers {
		log := models.ValidationLog{RecordIndex: i}
		errs := []string{}

		if !accNumRegex.MatchString(sample.AccountNo) {
			errs = append(errs, "Account number must be numeric")
		}

		var matched *models.Customer
		for _, c := range customers {
			if c.AccountNumber == sample.AccountNo {
				matched = &c
				break
			}
		}
		if matched == nil {
			errs = append(errs, "account number not found in canonical list")
		}

		if matched != nil && !strings.EqualFold(strings.TrimSpace(matched.CustomerName), strings.TrimSpace(sample.CustomerName)) {
			errs = append(errs, "name does not match canonical list")
		}

		if len(errs) > 0 {
			log.Verified = false
			log.Errors = errs
		} else {
			sample.Verified = true
			log.Verified = true
			matched.Verified = true
			log.NormalizedRecord = matched
			validatedCustomers = append(validatedCustomers, *matched)
		}
		logs = append(logs, log)

	}
	return logs, validatedCustomers, nil

}

func CalculateRating(transaction []models.Transaction, balance float64) payloads.RatingBreakdown {
	if len(transaction) == 0 {
		return payloads.RatingBreakdown{FinalScore: 1}
	}

	countScore := math.Min(10, float64(len(transaction)))

	total := 0.0
	for _, t := range transaction {
		total += math.Abs(t.Amount)
	}

	volumeScore := math.Min(10, total/1000)

	firstDate := time.UnixMilli(transaction[0].TransactionDate)

	lastDate := time.UnixMilli(transaction[len(transaction)-1].TransactionDate)

	if lastDate.Before(firstDate) {
		firstDate, lastDate = lastDate, firstDate
	}

	days := lastDate.Sub(firstDate).Hours() / 24
	months := days / 30.0
	durationScore := math.Min(10, months)

	if durationScore < 1 && days > 0 {
		durationScore = 1
	}

	stabilityScore := 10.0
	minBalance := transaction[0].ClearedBalance
	for _, t := range transaction {
		if t.ClearedBalance < minBalance {
			minBalance = t.ClearedBalance
		}
	}

	if minBalance < 0 {
		stabilityScore = 2.0
	}

	final := 0.3*countScore + 0.3*volumeScore + 0.2*durationScore + 0.2*stabilityScore
	if final > 10 {
		final = 10
	}

	return payloads.RatingBreakdown{
		CountScore:     math.Round(countScore*100) / 100,
		VolumeScore:    math.Round(volumeScore*100) / 100,
		DurationScore:  durationScore,
		StabilityScore: math.Round(stabilityScore*100) / 100,
		FinalScore:     math.Round(final*100) / 100,
	}
}
