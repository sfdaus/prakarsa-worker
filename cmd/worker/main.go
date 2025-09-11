package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"prakarsa-app/config"
	"prakarsa-app/infrastructure/datastore"
	"prakarsa-app/usecase"
	"prakarsa-app/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	httpDelivery "prakarsa-app/delivery/http"
	pgsqlRepository "prakarsa-app/repository/pgsql"
	redisRepository "prakarsa-app/repository/redis"
)

// @title Go Boilerplate
// @version 1.0.4
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey JwtToken
// @in header
// @name Authorization
func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup infra
	dbInstance, err := datastore.NewDatabase(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	cacheInstance, err := datastore.NewCache(configApp.CacheURL)
	utils.PanicIfNeeded(err)

	// Setup repository
	redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	workerRepo := pgsqlRepository.NewPgsqlReferenceRepository(dbInstance)

	// Setup usecase
	ctxTimeout := time.Duration(configApp.ContextTimeout) * time.Second
	workerUC := usecase.NewReferenceUsecase(workerRepo, redisRepo, ctxTimeout)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Logger.Info("🚀 Server is alive and running")

	// Setup handler
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	httpDelivery.NewReferenceHandler(e, workerUC)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configApp.ContextTimeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
