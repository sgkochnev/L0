package app

import (
	"LO/internal/cache"
	"LO/internal/config"
	"LO/internal/http-server/handlers/get"
	"LO/internal/http-server/middleware"
	libcache "LO/internal/lib/cache"
	"LO/internal/storage/pg"
	"LO/internal/subscriber"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpserver "LO/pkg/http_server"

	"github.com/gorilla/mux"
)

func Run() error {
	cfg := config.MastLoad()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	store, err := pg.New(context.Background(), connString,
		pg.WithMaxConnLifetime(cfg.Database.MaxConnLifetime),
		pg.WithMaxConnLifetimeJitter(cfg.Database.MaxConnLifetimeJitter),
		pg.WithMaxConnIdleTime(cfg.Database.MaxConnIdelTime),
		pg.WithHealthCheckPeriod(cfg.Database.HealthCheckPeriod),
		pg.WithMaxConns(cfg.Database.MaxConn),
		pg.WithMinConns(cfg.Database.MinConn),
	)
	if err != nil {
		return fmt.Errorf("app.Run.pg.New: %w", err)
	}
	defer func() {
		store.Close()
		log.Println("storage closed")
	}()

	appCache, err := cache.New(context.Background())
	if err != nil {
		return fmt.Errorf("app.Run.cache.New: %w", err)
	}
	defer func() {
		if err := appCache.Close(); err != nil {
			log.Println("app.Run.appCache.Close: ", err)
		}
		log.Println("cache deleted")
	}()

	err = libcache.Recovery(appCache, store)
	if err != nil {
		return fmt.Errorf("can not recover cache: %w", err)
	}

	sc, err := subscriber.New(
		store, appCache,
		cfg.Nats.ClusterID,
		cfg.Nats.ClientID,
		cfg.Nats.Address,
	)
	if err != nil {
		return fmt.Errorf("app.Run.subscriber.New: %w", err)
	}
	defer func() {
		if err := sc.Close(); err != nil {
			log.Println("app.Run.sc.Close: ", err)
		}
		log.Println("subscriber closed")
	}()

	err = sc.Subscribe(cfg.Nats.Subject)
	if err != nil {
		return fmt.Errorf("app.Run.sc.Subscribe: %w", err)
	}

	router := mux.NewRouter()

	router.Use(middleware.CORS)

	router.Handle("/orders/{order_uid}", get.New(appCache)).
		Methods(http.MethodGet)
	router.Handle("/orders", get.New(appCache)).
		Methods(http.MethodPost)

	server := httpserver.New(router,
		httpserver.WithAddress(cfg.HttpServer.Address),
		httpserver.WithIdleTimeout(cfg.HttpServer.IdelTimeout),
		httpserver.WithShutdownTimeout(cfg.HttpServer.ShutdownTimeout),
		httpserver.WithReadTimeout(cfg.HttpServer.Timeout),
		httpserver.WithWriteTimeout(cfg.HttpServer.Timeout),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app.Run.signal: %s", s)
	case err := <-server.Notify():
		log.Printf("app.Run.serverNotify: %s", err)
	}

	// shutdown server
	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("app.Run.serverShutdown: %w", err)
	}
	return nil
}
