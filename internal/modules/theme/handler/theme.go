package handler

import (
	"context"
	db "portfolioed/internal/db/sqlc"
	"portfolioed/internal/modules/theme/domain"
)

func (h *themeHandler) GetVariablesByType(ctx context.Context, varType string) ([]db.CssVariable, error) {
	return h.useCase.GetVariablesByType(ctx, varType)
}
func (h *themeHandler) UpdateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error {
	return h.useCase.UpdateVariable(ctx, req)
}
