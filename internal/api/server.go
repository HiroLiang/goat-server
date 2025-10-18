package api

import (
	"net/http"
	"time"

	"github.com/HiroLiang/goat-server/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func NewServer(addr string) *http.Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Register Swagger
	routes.RegisterSwaggerRoutes(r)

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
