package payloads

type SampleCustomer struct {
	CustomerName string `json:"customerName"`
	AccountNo    string `json:"accountNo"`
	Verified     bool   `json:"verified"`
}
