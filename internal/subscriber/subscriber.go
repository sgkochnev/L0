package subscriber

import (
	"LO/internal/models"
	"context"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

type Cacher interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

type Saver interface {
	Save(ctx context.Context, order *models.Order) error
}

type Subscriber struct {
	conn         stan.Conn
	subscription stan.Subscription
	db           Saver
	cache        Cacher
}

func New(store Saver, cache Cacher, clusterID, clientID, natsURL string) (*Subscriber, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	sub := &Subscriber{
		conn:  sc,
		db:    store,
		cache: cache,
	}

	return sub, nil
}

func (s *Subscriber) Subscribe(subject string) error {
	sub, err := s.conn.Subscribe(subject, s.msgHandler)
	if err != nil {
		return err
	}

	s.subscription = sub

	return nil
}

func (s *Subscriber) Close() error {
	return s.conn.Close()
}

func (s *Subscriber) msgHandler(msg *stan.Msg) {
	order := models.Order{}
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Printf("cannot unmarshal msg: %v", err)
		return
	}

	err = validator.New().Struct(&order)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		log.Println(models.ValidationErrors(errs))
		return
	}

	err = s.db.Save(context.Background(), &order)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.cache.Set(msg.Subject, msg.Data)
	if err != nil {
		log.Printf("cannot set cache: %v", err)
		return
	}
}
