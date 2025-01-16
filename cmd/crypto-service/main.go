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

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/sixojke/crypto-service/internal/config"
	"github.com/sixojke/crypto-service/internal/delivery"
	"github.com/sixojke/crypto-service/internal/repository"
	"github.com/sixojke/crypto-service/internal/server"
	"github.com/sixojke/crypto-service/internal/service"
	"github.com/sixojke/crypto-service/pkg/database"
	"github.com/sixojke/crypto-service/pkg/logger"
)

func main() {
	// Init config
	cfg, err := config.Init([]string{"configs"}, ".env")
	if err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	// Init logger
	enableLogger(cfg.Logger.LogLevel)

	// Init postgres
	postgres, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logger.Fatalf("error connect postgres db: %v", err)
	}
	defer postgres.Close()
	logger.Info("[POSTGRES] Connection successful")

	// Init repository
	repo := repository.NewCurrencyPostgres(postgres)

	// Init service
	service := service.NewCurrencyService(repo)

	// Init handler
	handler := delivery.NewHandler(service)

	// Start server
	srv := server.NewServer(cfg.HTTPServer, handler.Init())
	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %v\n", err)
		}
	}()
	logger.Infof("[SERVER] Started on port :%v", cfg.HTTPServer.Port)

	shutdown(srv, postgres)
}

// enableLogger initializes the logger
func enableLogger(logLevel int) {
	logger.NewLogger(zerolog.Level(logLevel), os.Stdout)
}

func shutdown(srv *server.Server, postgres *sqlx.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 3 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	postgres.Close()
}
