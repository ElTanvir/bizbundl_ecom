package handler

import (
	"context"
	db "portfolioed/internal/db/sqlc"
	"portfolioed/internal/modules/theme/domain"
	"portfolioed/internal/modules/theme/usecase"
)

type ThemeHandler interface {
	GetVariablesByType(ctx context.Context, varType string) ([]db.CssVariable, error)
	UpdateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error
}

type themeHandler struct {
	useCase usecase.ThemeUseCase
}

func NewThemeHandler(useCase usecase.ThemeUseCase) ThemeHandler {
	return &themeHandler{
		useCase: useCase,
	}
}
