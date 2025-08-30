package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
)

func initPayment(form *models.PaymentRequest, config config.Config) (*models.PaymentResponse, error) {

	body := models.PaymentRequest{
		Amount:             form.Amount,
		Currency:           "ETB",
		Email:              form.Email,
		FirstName:          form.FirstName,
		PhoneNumber:        form.PhoneNumber,
		TxRef:              "payer-" + form.TxRef,
		CallbackURL:        form.CallbackURL,
		ReturnURL:          form.ReturnURL,
		CustomizationTitle: form.CustomizationTitle,
		CustomizationDesc:  form.CustomizationDesc,
	}

	header := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", config.Chapa.SecretKey),
		"Content-Type":  "application/json",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return &models.PaymentResponse{Status: false, Message: "Failed to create payment request"}, err
	}

	req, err := http.NewRequest("POST", config.Chapa.PaymentEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return &models.PaymentResponse{Status: false, Message: "Failed to create payment request"}, err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &models.PaymentResponse{Status: false, Message: err.Error()}, err
	}

	defer resp.Body.Close()

	// Read the response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &models.PaymentResponse{Status: false, Message: "Failed to read response"}, err
	}

	// Parse the response
	var chapaResponse map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &chapaResponse); err != nil {
		return &models.PaymentResponse{Status: false, Message: "Failed to parse response"}, err
	}

	return &models.PaymentResponse{Status: true, ChapaResponse: chapaResponse}, nil

}

func (s *Server) HandlePaymentInChapa(w http.ResponseWriter, r *http.Request) error {
	// here the user must pass the book_id to create the payment in chapage
	// extract the room information from the book like payment

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}
	payload := models.ChapaPaymentPayload{}

	err = json.Unmarshal(reqBody, &payload)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	_, ok := r.Context().Value("user").(*models.UserClaims)
	if !ok {
		return utils.SendErrorResponse(w, "Unauthorized", "unauthorized", http.StatusUnauthorized)
	}

	book, err := s.Graph.GetBookById(payload.BookID)
	if err != nil {
		return utils.SendErrorResponse(w, "error in retriving book information", "database_error", http.StatusInternalServerError)
	}
	if book == nil {
		return utils.SendErrorResponse(w, "book not found", "book_not_found", http.StatusNotFound)
	}

	var chapaForm models.PaymentRequest
	price, err := strconv.ParseFloat(book.EnterpriseRoom.Charge, 64)
	if err != nil {
		return utils.SendErrorResponse(w, "error retriving book", "database_error", http.StatusInternalServerError)
	}

	formattedURL := fmt.Sprintf("%s/%x", s.Config.Chapa.CallbackUrl, book.ID)

	chapaForm.Amount = price
	chapaForm.Email = book.User.Email
	chapaForm.FirstName = book.User.FullName
	chapaForm.PhoneNumber = book.User.PhoneNumber
	txRef := uuid.New().String()
	chapaForm.TxRef = txRef
	chapaForm.CallbackURL = formattedURL
	chapaForm.ReturnURL = s.Config.Chapa.ReturnURL
	chapaForm.CustomizationTitle = "Book Payment"
	chapaForm.CustomizationDesc = "Book Payment"

	paymentResponse, err := initPayment(&chapaForm, *s.Config)
	if err != nil {
		return utils.SendErrorResponse(w, "error generating payment", "database_error", http.StatusInternalServerError)
	}
	if data, ok := paymentResponse.ChapaResponse["data"].(map[string]interface{}); ok {
		if checkoutUrl, ok := data["checkout_url"].(string); ok {

			insertingPayment := &models.InsertPayment{
				BookID:          book.ID,
				BasePrice:       book.EnterpriseRoom.Charge,
				PaymentMethod:   "chapa",
				NetTotal:        book.EnterpriseRoom.Charge,
				PhoneNumber:     book.User.PhoneNumber,
				Status:          "pending",
				PaymentProvider: "chapa",
				PaymentVia:      "payment_getway",
			}

			insertedPayment, err := s.Graph.CreatePayment(insertingPayment)
			if err != nil {
				return utils.SendErrorResponse(w, "error generating payment", "database_error", http.StatusInternalServerError)
			}

			// here we will save the payment information to the database
			responseData, _ := json.Marshal(models.ChapaPaymentResponse{
				PaymentURL: checkoutUrl,
				PaymentID:  insertedPayment.ID,
			})
			if _, err := w.Write(responseData); err != nil {
				return err
			}
		} else {
			return utils.SendErrorResponse(w, "error generating payment", "database_error", http.StatusInternalServerError)
		}
	} else {
		return utils.SendErrorResponse(w, "error generating payment", "database_error", http.StatusInternalServerError)
	}
	return nil
}
