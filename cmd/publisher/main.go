package main

import (
	"LO/internal/config"
	"bytes"
	_ "embed"
	"log"

	"github.com/nats-io/stan.go"
)

//go:embed fake_data.txt
var data []byte

func main() {
	cfg := config.MastLoad()

	sc, err := stan.Connect(
		cfg.Nats.ClusterID,
		cfg.Nats.ClientID+"2",
		stan.NatsURL(cfg.Nats.Address),
	)
	if err != nil {
		log.Fatalf("cannot connect to nats-streaming: %v", err)
	}
	defer func() {
		if err = sc.Close(); err != nil {
			log.Printf("cannot close stan connection: %v\n", err)
		}
	}()

	sData := bytes.Split(data, []byte("\n\n"))

	for i, data := range sData {
		err = sc.Publish(cfg.Nats.Subject, data)
		if err != nil {
			log.Printf("cannot publish: %v\n", err)
		}
		log.Printf("published data %d\n", i+1)
	}
}
