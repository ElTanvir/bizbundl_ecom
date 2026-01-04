package pages

import (
	db "bizbundl/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func calculateSubtotal(items []db.GetCartItemsRow) pgtype.Numeric {
	var total float64
	for _, item := range items {
		price, _ := item.BasePrice.Float64Value()
		total += price.Float64 * float64(item.Quantity)
	}

	// Convert back to Numeric (Simplified)
	var n pgtype.Numeric
	_ = n.Scan(total)
	return n
}
