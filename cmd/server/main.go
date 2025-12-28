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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HiroLiang/goat-server/internal/bootstrap"
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/database"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	_ "github.com/HiroLiang/goat-server/swag-docs"
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
	if err := redis.InitRedis(); err != nil {
		logger.Log.Fatal("Error initializing Redis", zap.Error(err))
	}

	// Init Database
	if err := database.InitDB(); err != nil {
		logger.Log.Fatal("Error initializing database", zap.Error(err))
	}
	defer database.CloseAllDBs()

	// Init all dependencies
	dependencies, err := bootstrap.BuildDeps()
	if err != nil {
		logger.Log.Fatal("Initialize dependencies failed", zap.Error(err))
	}

	// build use cases
	useCases := bootstrap.BuildUseCases(dependencies)

	// Start api server
	srv := bootstrap.NewServer(":"+config.Env("SERVER_PORT", "8080"), useCases, dependencies)
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

	// Block until receiving signal
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
