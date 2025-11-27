// @title Goat-Server
// @version 1.0.0
// @description Server for my all Goat application
// @host localhost:8080
// @basePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HiroLiang/goat-server/internal/api"
	"github.com/HiroLiang/goat-server/internal/api/handler"
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/HiroLiang/goat-server/internal/database"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	_ "github.com/HiroLiang/goat-server/docs"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] No .env file found")
	}
}

func main() {

	// Init logger
	logger.Init()
	defer logger.Stop()

	// Load Config Data
	if err := config.LoadConfig("./config"); err != nil {
		logger.Log.Fatal("Error loading config", zap.Error(err))
	}

	// Init Redis
	if err := database.InitRedis(); err != nil {
		logger.Log.Fatal("Error initializing Redis", zap.Error(err))
	}

	// Init Database
	if err := database.InitDB(); err != nil {
		logger.Log.Fatal("Error initializing database", zap.Error(err))
	}
	defer database.CloseAllDBs()

	// Test Database
	if db := database.GetDB(database.Postgres); db == nil {
		logger.Log.Fatal("Basic database not initialized")
	}

	// Init User Handler (Authentication)
	jwtSecret := config.Env("JWT_SECRET", "default-secret-change-in-production")
	jwtExpirationStr := config.Env("JWT_EXPIRATION_HOURS", "24")
	jwtExpiration := 24
	if _, err := fmt.Sscanf(jwtExpirationStr, "%d", &jwtExpiration); err != nil {
		logger.Log.Warn("Invalid JWT_EXPIRATION_HOURS, using default 24 hours")
	}
	handler.InitUserHandler(jwtSecret, jwtExpiration)

	// Start server
	srv := api.NewServer(":" + config.Env("SERVER_PORT", "8080"))
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Server start error", zap.Error(err))
		}
	}()
	defer func(srv *http.Server) {
		err := srv.Close()
		if err != nil {
			logger.Log.Error("Server close error", zap.Error(err))
		} else {
			logger.Log.Info("Server closed")
		}
	}(srv)

	// graceful shoutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutdown Server ...")

	// set timeout for closing server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// try to shoutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server shutdown error", zap.Error(err))
	} else {
		logger.Log.Info("Server exiting")
	}

}
