package server

import (
	"fmt"
	"portfolioed/internal/config"
	db "portfolioed/internal/db/sqlc"
	"portfolioed/internal/middleware"
	"portfolioed/token"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     *config.Config
	store      db.Store
	tokenMaker token.Maker
	router     *fiber.App
}

func NewServer(config *config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	app := fiber.New(fiber.Config{})
	app.Use(middleware.RouteCacheMiddleware())
	app.Use(recover.New())
	if config.Environment != "development" {
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))
	}

	app.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "SAMEORIGIN",
		HSTSMaxAge:                31536000,
		HSTSExcludeSubdomains:     false,
		HSTSPreloadEnabled:        true,
		ContentSecurityPolicy:     "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com https://cdn.jsdelivr.net https://cdn.tailwindcss.com https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src 'self' data: https:; font-src 'self' data: https://fonts.gstatic.com; connect-src 'self' https://cdn.jsdelivr.net;",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		CrossOriginEmbedderPolicy: "credentialless",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "cross-origin",
		OriginAgentCluster:        "?1",
		XDNSPrefetchControl:       "off",
		XDownloadOptions:          "noopen",
		XPermittedCrossDomain:     "none",
	}))

	// app.Use(cors.New(cors.Config{
	// 	AllowMethods: "GET,HEAD,PUT,PATCH,POST,DELETE,OPTIONS",
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// 	AllowOriginsFunc: func(origin string) bool {
	// 		u, err := url.Parse(origin)
	// 		if err != nil {
	// 			return false
	// 		}
	// 		h := u.Hostname()
	// 		return h == "localhost" || h == "127.0.0.1"
	// 	},
	// }))
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     app,
	}
	server.setupStatics()
	return server, nil
}

func (server *Server) Start() error {
	return server.router.Listen(":" + "8080")
}
func (server *Server) GetRouter() *fiber.App {
	return server.router
}
func (server *Server) GetDB() db.Store {
	return server.store
}
func (server *Server) GetTokenMaker() token.Maker {
	return server.tokenMaker
}
func (server *Server) GetConfig() *config.Config {
	return server.config
}
func (server *Server) setupStatics() {
	oneYearInSeconds := 31536000
	server.router.Static("/static", "./static", fiber.Static{
		MaxAge:        oneYearInSeconds,
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 365 * 24 * time.Hour,
	})
	server.router.Static("/demo", "./inspirations", fiber.Static{
		MaxAge:        oneYearInSeconds,
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 365 * 24 * time.Hour,
	})

	// Serve uploaded files publicly
	server.router.Static("/uploads", "./uploads", fiber.Static{
		MaxAge:        oneYearInSeconds,
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 24 * time.Hour,
	})
}
