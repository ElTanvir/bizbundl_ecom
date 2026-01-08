package server

import (
	"bizbundl/internal/config"
	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/infra/elastic"
	"bizbundl/internal/infra/redis"
	"bizbundl/internal/middleware"
	cacheStore "bizbundl/internal/store"
	"bizbundl/token"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	redisClient "github.com/redis/go-redis/v9"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     *config.Config
	store      db.DBStore
	tokenMaker token.Maker
	router     *fiber.App
	redis      *redisClient.Client
	elastic    *elasticsearch.Client
}

func NewServer(config *config.Config, store db.DBStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// Initialize Infra
	rc := redis.NewRedisClient(config)
	// Initialize Global Cache Store (Settings etc)
	cacheStore.Init(cacheStore.NewRedisStore(rc))

	es, err := elastic.NewElasticClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to init elastic: %w", err)
	}

	app := fiber.New(fiber.Config{})
	app.Use(etag.New())
	app.Use(cache.New(cache.Config{
		Expiration:   1 * time.Minute,
		CacheControl: true,
		Storage:      redis.NewFiberStorage(rc),
	}))
	app.Use(recover.New())
	app.Use(middleware.TenancyMiddleware(store))
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
		ContentSecurityPolicy:     "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com https://cdn.jsdelivr.net https://cdn.tailwindcss.com https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src 'self' data: https:; font-src 'self' data: https://fonts.gstatic.com; connect-src 'self' https://cdn.jsdelivr.net https://unpkg.com;",
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
	// 	},
	// }))

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     app,
		redis:      rc,
		elastic:    es,
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
func (server *Server) GetDB() db.DBStore {
	return server.store
}
func (server *Server) GetTokenMaker() token.Maker {
	return server.tokenMaker
}
func (server *Server) GetConfig() *config.Config {
	return server.config
}

func (server *Server) GetRedis() *redisClient.Client {
	return server.redis
}

func (server *Server) GetElastic() *elasticsearch.Client {
	return server.elastic
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

	// Serve uploaded files publicly
	server.router.Static("/uploads", "./uploads", fiber.Static{
		MaxAge:        oneYearInSeconds,
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 24 * time.Hour,
	})
}
