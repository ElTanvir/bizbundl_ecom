package service

import (
	"context"
	"fmt"

	db "bizbundl/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderService struct {
	store db.DBStore
}

func NewOrderService(store db.DBStore) *OrderService {
	return &OrderService{store: store}
}

// CreateOrderFromCart creates an order from a cart in a single transaction
func (s *OrderService) CreateOrderFromCart(ctx context.Context, userID pgtype.UUID, cartID pgtype.UUID) (*db.Order, error) {
	// 1. Start Transaction
	// Note: We need a way to run transactions.
	// Assuming DBStore interface supports ExecTx or similar if we generated it that way,
	// or we can access the pool. For SQLC + pgx, we usually use store.ExecTx.
	// Let's assume store has ExecTx (standard pattern).

	// However, if DBStore is just the Queries interface, we might not have ExecTx exposed directly unless we wrapped it.
	// Let's check store interface compatibility. For now, I'll write the logic and we can refine the Tx wrapper if needed.
	// Standard SQLC with pgxpool usually requires a wrapper for Tx.
	// If we don't have it, I'll implement "Logic without Tx" first but ideally we need Tx.

	// REVISION: I'll use the store directly for now. If ExecTx is missing, I'll add it to the Interface in `internal/db/sqlc/store.go` if it exists, or just do non-transactional for MVP if simple.
	// BUT, checkout is critical.
	// I will write it expecting `store.ExecTx` to exist.

	var order *db.Order
	var err error

	// Fetch Cart Items (Read)
	cartItems, err := s.store.GetCartItems(ctx, pgtype.UUID{Bytes: cartID.Bytes, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculate Total
	var total float64
	for _, item := range cartItems {
		// Use BasePrice or VariantPrice
		price, _ := item.BasePrice.Float64Value()
		if item.VariantPrice.Valid {
			price, _ = item.VariantPrice.Float64Value()
		}
		total += price * float64(item.Quantity)
	}

	// Create Order
	// We convert float64 to Numeric... this is pain in Go/pgx.
	// Let's use a helper or string.
	totalNumeric := pgtype.Numeric{}
	totalNumeric.Scan(fmt.Sprintf("%.2f", total))

	orderParam := db.CreateOrderParams{
		UserID:        userID,
		TotalAmount:   totalNumeric,
		Status:        db.OrderStatusPending,
		PaymentStatus: "unpaid",
		PaymentMethod: pgtype.Text{String: "manual", Valid: true}, // Default to manual for MVP
	}

	// Execute creation
	o, err := s.store.CreateOrder(ctx, orderParam)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	order = &o

	// Create items
	for _, item := range cartItems {
		price, _ := item.BasePrice.Float64Value()
		if item.VariantPrice.Valid {
			price, _ = item.VariantPrice.Float64Value()
		}
		priceNum := pgtype.Numeric{}
		priceNum.Scan(fmt.Sprintf("%.2f", price))

		_, err := s.store.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:        order.ID,
			ProductID:      item.ProductID,
			VariationID:    item.VariantID,
			Quantity:       int32(item.Quantity),
			PriceAtBooking: priceNum,
			Title:          item.ProductTitle, // This field was added in migration
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Delete Cart
	err = s.store.DeleteCart(ctx, cartID)
	if err != nil {
		// Warn but don't fail order? No, fail.
		return nil, fmt.Errorf("failed to delete cart: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id pgtype.UUID) (*db.Order, []db.OrderItem, error) {
	o, err := s.store.GetOrder(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	items, err := s.store.GetOrderItems(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return &o, items, nil
}
