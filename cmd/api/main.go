package main

import (
	"os"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/logger"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/app"
)

func main() {
	logger := logger.New()

	logger.Info("starting app")
	app.LoadEnv(".env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:root@localhost:5432/scholarly?sslmode=disable"
	}
	logger.Info("trying to connect to db")
	db, err := postgres.NewConn(dbURL, 20, 10, time.Minute)
	if err != nil {
		logger.Error("main: error while connecting to db", "error", err)
		os.Exit(1)
	}
	logger.Info("db connection successful")
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("main: error while closing db conn", "error", err)
		}
	}()

	logger.Info("trying to connect to cache")
	cacheURL := os.Getenv("CACHE_URL")
	if cacheURL == "" {
		cacheURL = "localhost:6379"
	}
	rdb := valkey.NewConn(cacheURL, "", 0)
	if rdb == nil {
		logger.Error("main: error while connecting to cache")
		os.Exit(1)
	}
	logger.Info("cache connection successful")

	logger.Info("trying to start server")

	app := app.New(db, rdb, logger)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	if err := app.Handler.Run(app.Handler.Mount(), addr); err != nil {
		logger.Error("main: error while booting HTTP server", "error", err)
		os.Exit(1)
	}
	logger.Info("server has been shut down")
}
