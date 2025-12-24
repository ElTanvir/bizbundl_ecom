package repository

import (
	"context"
	"fmt"
	db "portfolioed/internal/db/sqlc"
	"portfolioed/internal/modules/theme/domain"
	"sync"
	"time"
)

type ThemeRepository interface {
	GenerateCSS(ctx context.Context) (string, error)
	InvalidateCache(ctx context.Context) error
	GetAllVariables(ctx context.Context) ([]db.CssVariable, error)
	GetVariablesByType(ctx context.Context, varType string) ([]db.CssVariable, error)
	GetCSSVariableByName(ctx context.Context, name string) (db.CssVariable, error)
	CreateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error
	UpdateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error
	DeleteVariable(ctx context.Context, id string) error
}

type themeRepository struct {
	db         db.Store
	cssCache   string
	cacheTime  time.Time
	cacheMutex sync.RWMutex
}

func NewThemeRepository(database db.Store) ThemeRepository {
	repo := &themeRepository{
		db: database,
	}
	go func() {
		ctx := context.Background()
		if err := repo.regenerateCSS(ctx); err != nil {
			fmt.Printf("Warning: failed to pre-generate CSS cache on startup: %v\n", err)
		}
	}()

	return repo
}
