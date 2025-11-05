package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/HiroLiang/goat-server/internal/api/routes"
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(addr string) *http.Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Init config
	initConfig(r)

	// Register Swagger
	if config.Env("APP_ENV", "main") == "dev" {
		routes.RegisterSwaggerRoutes(r)
	}

	// Register REST routes
	routes.RegisterRestRoutes(r)

	// Register WebSocket routes
	routes.RegisterWsRoutes(r)

	// Setting server
	return &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   0,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func initConfig(r *gin.Engine) {

	// setting cors
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasSuffix(origin, "hiroliang.com")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
