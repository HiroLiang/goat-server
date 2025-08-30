package api

import (
	"net/http"
	"time"

	"github.com/HiroLiang/goat-chat-server/internal/api/middleware"
	routes2 "github.com/HiroLiang/goat-chat-server/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func NewServer(addr string) *http.Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.ErrorHandler())

	// Register REST routes
	routes2.RegisterRestRoutes(r)

	// Register WebSocket routes
	routes2.RegisterWsRoutes(r)

	return &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   0,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
