package page_builder

import (
	"bizbundl/internal/modules/page_builder/service"
	"bizbundl/internal/server"
	"context"

	"github.com/rs/zerolog/log"
)

func Init(app *server.Server) *service.PageBuilderService {
	svc := service.NewPageBuilderService(app.GetDB())

	// Seed Defaults
	// Better to run in background or migration, but for MVP checking on startup is fine.
	// Use a detached context or app context if available?
	// Using Background for now.
	if err := svc.SeedDefaults(context.Background()); err != nil {
		log.Error().Err(err).Msg("Failed to seed default pages")
	} else {
		log.Info().Msg("PageBuilder: Default pages seeded/verified")
	}

	return svc
}
