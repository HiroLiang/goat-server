package features

import (
	"fmt"
	"net/http/httptest"

	"github.com/HiroLiang/goat-server/internal/api"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/cucumber/godog"
)

var (
	testServer *httptest.Server
	baseURL    string
)

// InitializeSuite runs once before all features
func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		logger.InitTestEnv()

		app := api.NewServer(":8080")
		testServer = httptest.NewServer(app.Handler)
		baseURL = testServer.URL
		fmt.Println(testServer.URL)
	})
	ctx.AfterSuite(func() {
		testServer.Close()
		logger.Stop()
	})
}
