package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/flareopti/deeptruth/docs"
	"github.com/flareopti/deeptruth/internal/config"
	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/handlers/articles"
	"github.com/flareopti/deeptruth/internal/handlers/authors"
	"github.com/flareopti/deeptruth/internal/handlers/root"
	"github.com/flareopti/deeptruth/internal/handlers/static"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title DeepTruth
// @version 0.0.1
// @description Api for DeepTruth project

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

	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	conn, err := pgxpool.New(context.Background(), cfg.Storage.Address)
	if err != nil {
		log.Error("Can't connect to db : ", err)
		os.Exit(1)
	}
	log.Info("Connected to db")

	query := db.New(conn)

	router.Get("/", root.New(log))
	router.Get("/static/*", static.New(log, "frontend/static"))
	router.Route("/api/articles", func(r chi.Router) {
		r.Get("/", articles.List(log, query))
		r.Get("/search", articles.Search(log, query))
		r.Post("/", articles.Create(log, query))
		r.Route("/{articleID}", func(r chi.Router) {
			r.Get("/", articles.Get(log, query))
			r.Patch("/", articles.UpdateRating(log, query))
			r.Delete("/", articles.Delete(log, query))
		})
	})
	router.Route("/api/authors", func(r chi.Router) {
		r.Get("/", authors.List(log, query))
		r.Post("/", authors.Create(log, query))
		r.Route("/{authorID}", func(r chi.Router) {
			r.Get("/", authors.Get(log, query))
			r.Patch("/", authors.UpdateRating(log, query))
			r.Delete("/", authors.Delete(log, query))
		})
	})

	router.Get("/swagger/*", httpSwagger.Handler())
	log.Info("Starting server", slog.String("address", cfg.HTTPServer.Address))
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
