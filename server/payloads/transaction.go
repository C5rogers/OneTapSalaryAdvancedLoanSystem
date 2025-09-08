package payloads

type TransactionPayload struct {
	ID                  string  `json:"id"`
	FromAccount         string  `json:"fromAccount"`
	ToAccount           string  `json:"toAccount"`
	Amount              string  `json:"amount"` // string in JSON → will parse
	Remark              string  `json:"remark"`
	TransactionType     string  `json:"transactionType"`
	RequestID           string  `json:"requestId"`
	Reference           string  `json:"reference"`
	ThirdPartyReference string  `json:"thirdPartyReference"`
	InstitutionID       *string `json:"institutionId"`
	ClearedBalance      string  `json:"clearedBalance"`  // also string → parse
	TransactionDate     string  `json:"transactionDate"` // ms timestamp as string
	BillerID            *string `json:"billerId"`
}
