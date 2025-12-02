package features

import (
	"fmt"
	"net/http/httptest"

	"github.com/HiroLiang/goat-server/internal/bootstrap"
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/cucumber/godog"
	"go.uber.org/zap"
)

var (
	testServer *httptest.Server
	baseURL    string
)

// InitializeSuite runs once before all features
func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		logger.InitTestEnv()

		if err := config.LoadConfig("../config"); err != nil {
			logger.Log.Fatal("Error loading config", zap.Error(err))
		}

		dependencies, err := bootstrap.BuildDeps()
		if err != nil {
			logger.Log.Fatal("Initialize dependencies failed", zap.Error(err))
		}

		useCases := bootstrap.BuildUseCases(dependencies)

		// For feature tests we only hit /api/test, so pass empty services
		app := bootstrap.NewServer(":8080", useCases, dependencies)
		testServer = httptest.NewServer(app.Handler)
		baseURL = testServer.URL
		fmt.Println(testServer.URL)
	})
	ctx.AfterSuite(func() {
		testServer.Close()
		logger.Stop()
	})
}
