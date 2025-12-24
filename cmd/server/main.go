package main

import (
	"bizbundl/internal/config"
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/views"
	"bizbundl/internal/server"
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
	migrationDir := "internal/db/migration"
	if cfg.InDocker == "true" {
		migrationDir = "/app/internal/db/migration"
	}
	err = util.RunMigrations(cfg.DBSourceURL(), migrationDir)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	util.RegisterTagName()
	dbStore := db.NewStore(connPool)
	app, err := server.NewServer(cfg, dbStore)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}
	views.Init(app)
	log.Fatal().Err(app.Start()).Msg("failed to start server")
}