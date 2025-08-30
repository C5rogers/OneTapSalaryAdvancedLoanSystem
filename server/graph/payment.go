package graph

import (
	"context"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
)

type billing_payments_insert_input map[string]interface{}
type billing_payment_method_enum string
type billing_payment_status_enum string

func (gc *LMSGraphClient) CreatePayment(payment *models.InsertPayment) (p *models.Payment, err error) {
	var mutation struct {
		InsertPayment struct {
			ID string `graphql:"id"`
		} `graphql:"insert_billing_payments_one(object:$object)"`
	}

	variables := map[string]interface{}{
		"object": billing_payments_insert_input{
			"book_id":          uuid(payment.BookID),
			"base_price":       payment.BasePrice,
			"payment_method":   payment.PaymentMethod,
			"net_total":        payment.NetTotal,
			"phone_number":     payment.PhoneNumber,
			"status":           billing_payment_status_enum(payment.Status),
			"payment_provider": payment.PaymentProvider,
			"payment_via":      billing_payment_method_enum(payment.PaymentVia),
		},
	}

	if err = gc.Mutate(context.Background(), &mutation, variables); err != nil {
		return nil, err
	}

	p = &models.Payment{
		ID: mutation.InsertPayment.ID,
	}
	return p, nil

}
