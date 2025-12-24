package repository

import (
	"context"
	"fmt"
	db "portfolioed/internal/db/sqlc"
	"portfolioed/internal/modules/theme/domain"
	"portfolioed/internal/store"
	"strings"
	"time"
)

func (s *themeRepository) GenerateCSS(ctx context.Context) (string, error) {
	s.cacheMutex.RLock()
	if s.cssCache != "" {
		css := s.cssCache
		s.cacheMutex.RUnlock()
		return css, nil
	}
	s.cacheMutex.RUnlock()
	if err := s.regenerateCSS(ctx); err != nil {
		return "", fmt.Errorf("failed to generate CSS: %w", err)
	}
	s.cacheMutex.RLock()
	css := s.cssCache
	s.cacheMutex.RUnlock()

	if css == "" {
		return "", fmt.Errorf("failed to retrieve CSS after regeneration")
	}

	return css, nil
}

// InvalidateCache clears the current CSS cache and regenerates it from the database.
func (s *themeRepository) InvalidateCache(ctx context.Context) error {
	s.cacheMutex.Lock()
	s.cssCache = ""
	s.cacheMutex.Unlock()

	if err := s.regenerateCSS(ctx); err != nil {
		return fmt.Errorf("failed to invalidate and regenerate cache: %w", err)
	}

	return nil
}

// --- CRUD for CSS Variables ---

func (s *themeRepository) GetAllVariables(ctx context.Context) ([]db.CssVariable, error) {
	return s.db.GetAllCSSVariables(ctx)
}

func (s *themeRepository) GetVariablesByType(ctx context.Context, varType string) ([]db.CssVariable, error) {
	return s.db.GetCSSVariablesByType(ctx, varType)
}

func (s *themeRepository) GetCSSVariableByName(ctx context.Context, name string) (db.CssVariable, error) {
	return s.db.GetCSSVariableByName(ctx, name)
}

func (s *themeRepository) CreateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error {
	err := s.db.CreateCSSVariable(ctx, db.CreateCSSVariableParams{
		Name:         req.Name,
		Value:        req.Value,
		VariableType: req.VariableType,
	})
	if err != nil {
		return fmt.Errorf("failed to create variable: %w", err)
	}
	go s.regenerateCSS(context.Background())
	return nil
}

func (s *themeRepository) UpdateVariable(ctx context.Context, req domain.CreateOrUpdateVariableRequest) error {
	err := s.db.UpdateCSSVariable(ctx, db.UpdateCSSVariableParams{
		Name:         req.Name,
		Value:        req.Value,
		VariableType: req.VariableType,
	})
	if err != nil {
		return fmt.Errorf("failed to update variable: %w", err)
	}

	// Regenerate cache in the background
	go s.regenerateCSS(context.Background())
	return nil
}

func (s *themeRepository) DeleteVariable(ctx context.Context, id string) error {
	err := s.db.DeleteCSSVariable(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete variable: %w", err)
	}
	go s.regenerateCSS(context.Background())
	return nil
}

func (s *themeRepository) generateCSSFromVariables(vars []db.CssVariable) string {
	var builder strings.Builder
	builder.WriteString(":root {\n")
	for _, v := range vars {
		builder.WriteString(fmt.Sprintf("  %s: %s;\n", v.Name, v.Value))
	}
	builder.WriteString("}\n")
	data := builder.String()
	store.SetCssVariable(data)
	return data
}

// regenerateCSS fetches all variables from the DB, builds the CSS, and caches it.
func (s *themeRepository) regenerateCSS(ctx context.Context) error {
	vars, err := s.db.GetAllCSSVariables(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all variables for CSS generation: %w", err)
	}

	css := s.generateCSSFromVariables(vars)
	s.cacheMutex.Lock()
	s.cssCache = css
	s.cacheTime = time.Now()
	s.cacheMutex.Unlock()

	return nil
}
