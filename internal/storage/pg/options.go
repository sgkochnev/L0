package pg

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Option func(*pgxpool.Config)

func WithMaxConnLifetime(timeout time.Duration) Option {
	return func(cfg *pgxpool.Config) {
		cfg.MaxConnLifetime = timeout
	}
}

func WithMaxConnLifetimeJitter(timeout time.Duration) Option {
	return func(cfg *pgxpool.Config) {
		cfg.MaxConnLifetimeJitter = timeout
	}
}

func WithMaxConnIdleTime(timeout time.Duration) Option {
	return func(cfg *pgxpool.Config) {
		cfg.MaxConnIdleTime = timeout
	}
}

func WithHealthCheckPeriod(timeout time.Duration) Option {
	return func(cfg *pgxpool.Config) {
		cfg.HealthCheckPeriod = timeout
	}
}

func WithMaxConns(count int32) Option {
	return func(cfg *pgxpool.Config) {
		cfg.MaxConns = count
	}
}

func WithMinConns(count int32) Option {
	return func(cfg *pgxpool.Config) {
		cfg.MinConns = count
	}
}
