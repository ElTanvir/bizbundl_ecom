package util

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// FormatPrice converts a pgtype.Numeric to a formatted currency string.
// It assumes the numeric is valid. If not, consistent with 0.00.
func FormatPrice(n pgtype.Numeric) string {
	f, err := n.Float64Value()
	if err != nil {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", f.Float64)
}
