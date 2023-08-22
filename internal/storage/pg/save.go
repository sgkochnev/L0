package pg

import (
	"LO/internal/models"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Save(ctx context.Context, order *models.Order) error {
	const op = "storage.pg.Save"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// INSERT into orders (order_uid, data) VALUES ($1, $2)
	query, params, err := psql.
		Insert("orders").
		Columns("order_uid", "data").
		Values(order.OrderUID, order).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot build query %s: %w", op, err)
	}

	_, err = s.conn.Exec(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("cannot save order: %s: %w", op, err)
	}

	return nil
}
