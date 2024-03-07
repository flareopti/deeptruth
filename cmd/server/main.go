package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/flareopti/deeptruth/internal/config"
	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/handlers/root"
	"github.com/flareopti/deeptruth/internal/handlers/static"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	router := chi.NewRouter()

	log.Info("Starting deeptruth backend", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/", func(r chi.Router) {
		r.Get("/", root.New(log))
	})

	static := static.New("frontend/static", "frontend/templates/index.html")
	router.Handle("/", static)

	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	testConn, err := pgxpool.New(context.Background(), cfg.Storage.Address)
	if err != nil {
		log.Error("Can't connect to db : ", err)
		os.Exit(1)
	}
	log.Info("Connected to db")

	testQuery := db.New(testConn)
	_ = testQuery

	log.Info("Starting server")
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server : ", err)
	}

	log.Error("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
