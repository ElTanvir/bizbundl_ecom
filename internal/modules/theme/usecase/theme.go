package usecase

import (
	"context"
	db "portfolioed/internal/db/sqlc"
	"portfolioed/internal/modules/theme/domain"
)

func (u *themeUseCase) GenerateCSS(ctx context.Context) (string, error) {
	return u.repo.GenerateCSS(ctx)
}
func (u *themeUseCase) InvalidateCache(ctx context.Context) error {
	return u.repo.InvalidateCache(ctx)
}

func (u *themeUseCase) GetAllVariables(ctx context.Context) ([]db.CssVariable, error) {
	return u.repo.GetAllVariables(ctx)
}

func (u *themeUseCase) GetVariablesByType(ctx context.Context, varType string) ([]db.CssVariable, error) {
	return u.repo.GetVariablesByType(ctx, varType)
}

func (u *themeUseCase) GetCSSVariableByName(ctx context.Context, name string) (db.CssVariable, error) {
	return u.repo.GetCSSVariableByName(ctx, name)
}
func (u *themeUseCase) CreateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error {
	return u.repo.CreateVariable(ctx, req)
}

func (u *themeUseCase) UpdateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error {
	return u.repo.UpdateVariable(ctx, req)
}

func (u *themeUseCase) DeleteVariable(ctx context.Context, id string) error {
	return u.repo.DeleteVariable(ctx, id)
}
