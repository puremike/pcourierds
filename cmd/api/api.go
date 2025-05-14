package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type app struct {
	config *config
	logger *zap.SugaredLogger
}

type config struct {
	port string
	env  string
}

func (app *app) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.health)
		docURL := fmt.Sprintf("%s/swagger/doc.json", app.config.port)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docURL)))
	})

	return r
}

func (app *app) start(mux http.Handler) error {

	srv := http.Server{
		Addr:         app.config.port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	// implementing graceful shutdown
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Infow("Shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Starting server on port:", "port", app.config.port, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdown; err != nil {
		return err
	}

	return nil
}
