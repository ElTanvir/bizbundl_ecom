package main

import (
	"bizbundl/internal/config"
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/modules/auth"
	"bizbundl/internal/modules/cart"
	"bizbundl/internal/modules/catalog"
	"bizbundl/internal/modules/order"
	"bizbundl/internal/modules/platform"
	"bizbundl/internal/server"
	"bizbundl/internal/views/frontend"
	"bizbundl/util"
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	connPool, err := pgxpool.New(context.Background(), cfg.DBSource())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	// Run Platform Migrations (Public Schema) on Startup
	// Tenant migrations are handled by the worker or when a shop is created.
	platformMsgDir := "internal/db/migration/platform"
	if cfg.InDocker == "true" {
		platformMsgDir = "/app/internal/db/migration/platform"
	}
	// Note: We might want to force "search_path=public" here to be safe,
	// though the default connection usually defaults to public.
	err = util.RunMigrations(cfg.DBSourceURL(), platformMsgDir)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run platform migrations")
	}
	util.RegisterTagName()
	dbStore := db.NewStore(connPool)
	app, err := server.NewServer(cfg, dbStore)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}
	// Initialize Modules
	auth.Init(app)
	catalogSvc := catalog.Init(app)
	cartSvc := cart.Init(app)
	order.Init(app, cartSvc, catalogSvc)
	platform.Init(app)

	frontend.Init(app)
	log.Fatal().Err(app.Start()).Msg("failed to start server")
}
