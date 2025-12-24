package main

import (
	"database/sql"
	"os"
	"portfolioed/internal/cloudflare"
	"portfolioed/internal/config"
	db "portfolioed/internal/db/sqlc"
	"portfolioed/internal/modules/menu"
	"portfolioed/internal/modules/root"
	"portfolioed/internal/modules/theme"
	"portfolioed/internal/server"
	"portfolioed/util"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := sql.Open("sqlite3", cfg.DBSource())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	defer conn.Close()
	migrationDir := "internal/db/migration"
	if cfg.InDocker == "true" {
		migrationDir = "/app/internal/db/migration"
	}
	err = util.RunMigrations(cfg.DBSource(), migrationDir)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	util.RegisterTagName()
	dbStore := db.NewStore(conn)
	app, err := server.NewServer(cfg, dbStore)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}
	theme.Init(app)
	root.Init(app)
	menu.Init(app)
	log.Fatal().Err(app.Start()).Msg("failed to start server")
}

func CacheRuleSetup(){
	cloudflare.Cache()
}