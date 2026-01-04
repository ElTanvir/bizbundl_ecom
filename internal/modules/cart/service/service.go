package service

import (
	"context"
	"errors"
	"fmt"

	db "bizbundl/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CartService struct {
	store db.DBStore
}

func NewCartService(store db.DBStore) *CartService {
	return &CartService{store: store}
}

// GetOrCreateCart retrieves an active cart or creates one
func (s *CartService) GetOrCreateCart(ctx context.Context, sessionID pgtype.UUID, userID pgtype.UUID) (db.Cart, error) {
	// 1. Try by Session ID first (Guest or User)
	if sessionID.Valid {
		cart, err := s.store.GetCartBySession(ctx, sessionID)
		if err == nil {
			return cart, nil
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return db.Cart{}, err
		}
	}

	// 2. If User ID provided, check by User
	if userID.Valid {
		cart, err := s.store.GetCartByUser(ctx, userID)
		if err == nil {
			return cart, nil
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return db.Cart{}, err
		}
	}

	// 3. Create New Cart
	// Note: If userID is valid, we link it. If sessionID is valid, we link it.
	// If both valid, link both.
	return s.store.CreateCart(ctx, db.CreateCartParams{
		SessionID: sessionID,
		UserID:    userID,
	})
}

func (s *CartService) AddToCart(ctx context.Context, sessionID pgtype.UUID, userID pgtype.UUID, productID pgtype.UUID, variantID pgtype.UUID, qty int32) (db.CartItem, error) {
	cart, err := s.GetOrCreateCart(ctx, sessionID, userID)
	if err != nil {
		return db.CartItem{}, fmt.Errorf("failed to get cart: %w", err)
	}

	// Add Item (Upsert)
	item, err := s.store.AddCartItem(ctx, db.AddCartItemParams{
		CartID:    cart.ID,
		ProductID: productID,
		VariantID: variantID, // Can be null
		Quantity:  qty,
	})
	return item, err
}

func (s *CartService) RemoveItem(ctx context.Context, itemID pgtype.UUID, cartID pgtype.UUID) error {
	return s.store.RemoveCartItem(ctx, db.RemoveCartItemParams{
		ID:     itemID,
		CartID: cartID,
	})
}

func (s *CartService) UpdateItemQuantity(ctx context.Context, itemID pgtype.UUID, cartID pgtype.UUID, quantity int32) (db.CartItem, error) {
	return s.store.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
		ID:       itemID,
		CartID:   cartID,
		Quantity: quantity,
	})
}

func (s *CartService) MergeCarts(ctx context.Context, sessionID pgtype.UUID, userID pgtype.UUID) error {
	// 1. Get Guest Cart (Session)
	guestCart, err := s.store.GetCartBySession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // check user cart existence?
		}
		return err
	}

	// 2. Get User Cart
	userCart, err := s.store.GetCartByUser(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// No user cart, just assign the guest cart to user
		return s.store.UpdateCartUser(ctx, db.UpdateCartUserParams{
			ID:     guestCart.ID,
			UserID: userID,
		})
	}

	// 3. Both exist (Merge Logic)
	// Complex: Move items from Guest Cart to User Cart.
	// Implementation: Get Items from Guest, Add to User, Delete Guest Cart.
	// Transaction recommended.
	// For MVP, simplistic loop.

	items, err := s.store.GetCartItems(ctx, guestCart.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		// Add to User Cart Upsert
		_, err := s.store.AddCartItem(ctx, db.AddCartItemParams{
			CartID:    userCart.ID,
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		})
		if err != nil {
			return err
		}
	}

	// Clear Guest Cart
	err = s.store.ClearCart(ctx, guestCart.ID)
	return err
}

func (s *CartService) GetCartItems(ctx context.Context, cartID pgtype.UUID) ([]db.GetCartItemsRow, error) {
	return s.store.GetCartItems(ctx, cartID)
}
