package pg

import (
	"LO/internal/storage"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn *pgxpool.Pool
}

func New(ctx context.Context, connectionString string, opts ...Option) (*Storage, error) {
	connConf, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Printf("cannot parse config: %v", err)
		return nil, storage.ErrConntctionString
	}

	for _, opt := range opts {
		opt(connConf)
	}

	connPool, err := pgxpool.NewWithConfig(ctx, connConf)
	if err != nil {
		log.Printf("cannot create pool: %v", err)
		return nil, storage.ErrCanNotConnect
	}

	err = connPool.Ping(context.Background())
	if err != nil {
		log.Printf("cannot ping: %v", err)
		return nil, storage.ErrNotResponding
	}

	return &Storage{conn: connPool}, nil
}

func (s *Storage) Close() {
	s.conn.Close()
}
