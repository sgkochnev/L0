package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Nats       `yaml:"nats"        env:"NATS"        env-required:"true"`
	Database   `yaml:"database"    env:"DATABASE"    env-required:"true"`
	HttpServer `yaml:"http_server" env:"HTTP_SERVER" env-required:"true"`
}

type Nats struct {
	Address   string `yaml:"address"    env:"NATS_ADDRESS"    env-default:"localhost:4222"`
	ClusterID string `yaml:"cluster_id" env:"NATS_CLUSTER_ID" env-default:"test-cluster"`
	ClientID  string `yaml:"client_id"  env:"NATS_CLIENT_ID"  env-default:"publisher"`
	Subject   string `yaml:"subject"    env:"NATS_SUBJECT"    env-default:"order-topic"`
}

type Database struct {
	Username              string        `yaml:"username"                 env:"DB_USERNAME"                 env-required:"true"`
	Password              string        `yaml:"password"                 env:"DB_PASSWORD"                 env-required:"true"`
	Name                  string        `yaml:"name"                     env:"DB_NAME"                     env-required:"true"`
	Host                  string        `yaml:"host"                     env:"DB_HOST"                                         env-default:"localhost"`
	Port                  string        `yaml:"port"                     env:"DB_PORT"                                         env-default:"5432"`
	MaxConnLifetime       time.Duration `yaml:"max_conn_lifetime"        env:"DB_MAX_CONN_LIFETIME"                            env-default:"15s"`
	MaxConnLifetimeJitter time.Duration `yaml:"max_conn_lifetime_jitter" env:"DB_MAX_CONN_LIFETIME_JITTER"                     env-default:"3s"`
	MaxConnIdelTime       time.Duration `yaml:"max_conn_idel_time"       env:"DB_MAX_CONN_IDEL_TIME"                           env-default:"300s"`
	HealthCheckPeriod     time.Duration `yaml:"health_check_period"      env:"DB_HEALTH_CHECK_PERIOD"                          env-default:"60s"`
	MaxConn               int32         `yaml:"max_conn"                 env:"DB_MAX_CONN"                                     env-default:"15"`
	MinConn               int32         `yaml:"min_conn"                 env:"DB_MIN_CONN"                                     env-default:"5"`
}

type HttpServer struct {
	Address         string        `yaml:"address"          env:"HTTP_SERVER_ADDRESS"          env-default:"localhost:8080"`
	Timeout         time.Duration `yaml:"timeout"          env:"HTTP_SERVER_TIMEOUT"          env-default:"5s"`
	IdelTimeout     time.Duration `yaml:"idel_timeout"     env:"HTTP_SERVER_IDEL_TIMEOUT"     env-default:"30s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" env-default:"3s"`
}

func MastLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config: ", err)
	}

	return &cfg
}
