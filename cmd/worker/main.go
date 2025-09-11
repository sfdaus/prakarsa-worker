package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"prakarsa-app/config"
	"prakarsa-app/utils"
	"syscall"
	"time"

	httpdelivery "prakarsa-app/delivery/http"
	"prakarsa-app/infrastructure/datastore"
	"prakarsa-app/repository"
	"prakarsa-app/transport/clients"
	"prakarsa-app/worker"
	"prakarsa-app/worker/handlers"
)

// @title Go Boilerplate
// @version 1.0.4
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey JwtToken
// @in header
// @name Authorization
func main() {
	configApp := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	db, err := datastore.MustOpen(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	mux := http.NewServeMux()
	mux.HandleFunc("/", httpdelivery.HealthHandler(db))
	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 2 * time.Second}
	go func() { _ = srv.ListenAndServe() }()
	defer srv.Shutdown(context.Background())

	notif := clients.NewNotifClient(configApp.NotificationServiceURL)
	repo := repository.NewOutboxRepo(db)

	runner := &worker.Runner{
		Repo:  repo,
		H:     handlers.NewNotifApproved(notif),
		Batch: 50,
		Idle:  2 * time.Second,
	}
	if err := runner.Start(ctx); err != nil {
		log.Println("runner stopped:", err)
	}

}
