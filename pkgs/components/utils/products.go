package utils

import db "bizbundl/internal/db/sqlc"

// GetProducList safely extracts a list of products from props.
func GetProducList(props map[string]interface{}, key string) []db.Product {
	if val, ok := props[key].([]db.Product); ok {
		return val
	}
	return nil
}
