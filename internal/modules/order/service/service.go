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

// OrderItemDTO helper for internal use
type OrderItemDTO struct {
	ProductID    pgtype.UUID
	VariantID    pgtype.UUID
	Quantity     int32
	UnitPrice    float64
	ProductTitle string
}

// createOrderCore handles the actual DB insertion
func (s *OrderService) createOrderCore(ctx context.Context, userID pgtype.UUID, items []OrderItemDTO, total float64) (*db.Order, error) {
	totalNumeric := pgtype.Numeric{}
	totalNumeric.Scan(fmt.Sprintf("%.2f", total))

	paymentMethod := "manual"
	paymentStatus := "unpaid"

	orderParam := db.CreateOrderParams{
		UserID:        userID,
		TotalAmount:   totalNumeric,
		Status:        db.NullOrderStatus{OrderStatus: db.OrderStatusPending, Valid: true},
		PaymentStatus: &paymentStatus,
		PaymentMethod: &paymentMethod,
	}

	// Execute creation
	o, err := s.store.CreateOrder(ctx, orderParam)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create items
	for _, item := range items {
		priceNum := pgtype.Numeric{}
		priceNum.Scan(fmt.Sprintf("%.2f", item.UnitPrice))

		_, err := s.store.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:        o.ID,
			ProductID:      item.ProductID,
			VariationID:    item.VariantID,
			Quantity:       item.Quantity,
			PriceAtBooking: priceNum,
			Title:          item.ProductTitle,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}
	return &o, nil
}

// CreateOrderFromCart creates an order from a cart
func (s *OrderService) CreateOrderFromCart(ctx context.Context, userID pgtype.UUID, cartID pgtype.UUID) (*db.Order, error) {
	// Fetch Cart Items
	cartItems, err := s.store.GetCartItems(ctx, pgtype.UUID{Bytes: cartID.Bytes, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Prepare DTOs and Total
	var total float64
	var orderItems []OrderItemDTO

	for _, item := range cartItems {
		var price float64
		baseF8, _ := item.BasePrice.Float64Value()
		price = baseF8.Float64

		if item.VariantPrice.Valid {
			varF8, _ := item.VariantPrice.Float64Value()
			price = varF8.Float64
		}
		total += price * float64(item.Quantity)

		orderItems = append(orderItems, OrderItemDTO{
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			Quantity:     int32(item.Quantity),
			UnitPrice:    price,
			ProductTitle: item.ProductTitle,
		})
	}

	// Core Creation
	order, err := s.createOrderCore(ctx, userID, orderItems, total)
	if err != nil {
		return nil, err
	}

	// Delete Cart
	if err := s.store.DeleteCart(ctx, cartID); err != nil {
		return nil, fmt.Errorf("failed to delete cart: %w", err)
	}

	return order, nil
}

// CreateOrderDirect creates an order for a single product/variant immediately
func (s *OrderService) CreateOrderDirect(ctx context.Context, userID pgtype.UUID, productID pgtype.UUID, variantID pgtype.UUID, quantity int32) (*db.Order, error) {
	// 1. Fetch Product
	p, err := s.store.GetProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	var price float64
	baseF8, _ := p.BasePrice.Float64Value()
	price = baseF8.Float64
	title := p.Title

	// 2. Fetch Variant if provided
	if variantID.Valid {
		v, err := s.store.GetProductVariant(ctx, variantID)
		if err != nil {
			return nil, fmt.Errorf("variant not found: %w", err)
		}
		varF8, _ := v.Price.Float64Value()
		price = varF8.Float64
		// Append variant title? Or just send main title.
		// Usually "Product Title - Variant Title" or just "Product Title" and let OrderItem store IDs.
		// OrderItem has Title snapshot field.
		title = fmt.Sprintf("%s - %s", p.Title, v.Title)
	}

	total := price * float64(quantity)

	item := OrderItemDTO{
		ProductID:    productID,
		VariantID:    variantID,
		Quantity:     quantity,
		UnitPrice:    price,
		ProductTitle: title,
	}

	return s.createOrderCore(ctx, userID, []OrderItemDTO{item}, total)
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
