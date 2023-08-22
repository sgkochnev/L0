package cache

import (
	"LO/internal/models"
	"context"
)

type Setter interface {
	Set(key string, value []byte) error
}

type Getter interface {
	GetAll(ctx context.Context) ([]models.OrderData, error)
}

func Recovery(c Setter, db Getter) error {
	orders, err := db.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := c.Set(order.OrderUID, order.Data); err != nil {
			return err
		}
	}
	return nil
}
