package main

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/sixojke/crypto-service/internal/config"
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
}

func enableLogger(logLevel int) {
	logger.NewLogger(zerolog.Level(logLevel), os.Stdout)
}
