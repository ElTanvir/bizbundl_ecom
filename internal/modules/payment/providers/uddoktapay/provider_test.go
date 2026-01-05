package uddoktapay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "bizbundl/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestInitPayment(t *testing.T) {
	// 1. Mock Server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Method
		assert.Equal(t, "POST", r.Method)

		// Verify Headers
		assert.Equal(t, "test-api-key", r.Header.Get("RT-UDDOKTAPAY-API-KEY"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify Body
		var reqBody initRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "100.50", reqBody.Amount)
		assert.Equal(t, "customer@example.com", reqBody.Email)

		// Send Mock Response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(initResponse{
			Status:     true,
			Message:    "Success",
			PaymentURL: "https://sandbox.uddoktapay.com/checkout/12345",
		})
	}))
	defer mockServer.Close()

	// 2. Init Provider with Mock URL
	provider := New("test-api-key")
	provider.BaseURL = mockServer.URL // Override BaseURL for testing

	// 3. Create Mock Order
	amount := pgtype.Numeric{}
	amount.Scan("100.50")
	id := pgtype.UUID{}
	id.Scan("550e8400-e29b-41d4-a716-446655440000") // Valid UUID

	order := &db.Order{
		ID:          id,
		TotalAmount: amount,
	}

	// 4. Call InitPayment
	url, err := provider.InitPayment(order, "customer@example.com")

	// 5. Assertions
	assert.NoError(t, err)
	assert.Equal(t, "https://sandbox.uddoktapay.com/checkout/12345", url)
}

func TestInitPayment_Error(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(initResponse{
			Status:  false,
			Message: "Invalid API Key",
		})
	}))
	defer mockServer.Close()

	provider := New("bad-key")
	provider.BaseURL = mockServer.URL

	amount := pgtype.Numeric{}
	amount.Scan("50.00")
	orderID := pgtype.UUID{}
	orderID.Scan("550e8400-e29b-41d4-a716-446655440000")

	order := &db.Order{ID: orderID, TotalAmount: amount}

	_, err := provider.InitPayment(order, "test@test.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid API Key")
}
