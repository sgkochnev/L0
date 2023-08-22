package pg

import (
	"LO/internal/models"
	"LO/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	const op = "storage.GetByID"

	data, err := s.getByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var order models.Order
	err = json.Unmarshal(data.Data, &order)
	if err != nil {
		log.Printf("cannot unmarshal order %s: %w", op, err)
		return nil, fmt.Errorf("cannot unmarshal order %s: %w", op, err)
	}

	return &order, nil
}

func (s *Storage) GetByID(ctx context.Context, id string) (*models.OrderData, error) {
	return s.getByID(ctx, id)
}

func (s *Storage) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	const op = "storage.GetAll"

	rows, err := s.getAllRows(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var data models.OrderData

		err = rows.Scan(&data)
		if err != nil {
			return nil, fmt.Errorf("cannot scan order %s: %w", op, err)
		}

		var order models.Order
		err = json.Unmarshal(data.Data, &order)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal order %s: %w", op, err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]models.OrderData, error) {
	const op = "storage.GetAll"

	rows, err := s.getAllRows(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderData

	for rows.Next() {
		var order models.OrderData

		err = rows.Scan(&order.OrderUID, &order.Data)
		if err != nil {
			return nil, fmt.Errorf("cannot scan order %s: %w", op, err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) getAllRows(ctx context.Context) (pgx.Rows, error) {
	const op = "storage.getAllRows"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// SELECT order_uid, data FROM orders
	query, params, err := psql.
		Select("order_uid", "data").
		From("orders").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build query %s: %w", op, err)
	}

	rows, err := s.conn.Query(ctx, query, params...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrOrderNotFound
		}
		return nil, fmt.Errorf("cannot get orders %s: %w", op, err)
	}

	return rows, nil
}

func (s *Storage) getByID(ctx context.Context, id string) (*models.OrderData, error) {
	const op = "storage.pg.getByID"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// SELECT order_uid, data FROM orders WHERE order_uid = $1
	query, params, err := psql.
		Select("order_uid", "data").
		From("orders").
		Where(sq.Eq{
			"order_uid": id,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build query %s: %w", op, err)
	}

	var order models.OrderData

	err = s.conn.QueryRow(ctx, query, params...).Scan(&order.OrderUID, &order.Data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrOrderNotFound
		}
		return nil, fmt.Errorf("cannot get order %s: %w", op, err)
	}
	return &order, nil
}
