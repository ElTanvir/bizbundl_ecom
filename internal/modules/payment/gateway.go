package payment

import (
	db "bizbundl/internal/db/sqlc"
)

type PaymentInfo struct {
	TransactionID string
	Status        string // "COMPLETED", "PENDING", "FAILED"
	Amount        float64
}

type Gateway interface {
	// InitPayment initializes a payment request and returns the payment URL
	InitPayment(order *db.Order, customerEmail string) (string, error)

	// VerifyPayment checks the status of a payment via invoice/transaction ID
	VerifyPayment(invoiceID string) (*PaymentInfo, error)
}
