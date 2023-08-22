docker_up:
	export POSTGRES_DB="L0" && export POSTGRES_USERNAME="l0" && export POSTGRES_PASSWORD="secret_password" && docker compose up -d

docker_down:
	docker compose down

migrations_up:
	migrate -path ./migrations -database 'postgres://l0:secret_password@localhost:5432/L0?sslmode=disable' up

migrations_down:
	migrate -path ./migrations -database 'postgres://l0:secret_password@localhost:5432/L0?sslmode=disable' down

run_publisher:
	export CONFIG_PATH="./config/config.yaml" && go run ./cmd/publisher/main.go

run_subscriber:
	export CONFIG_PATH="./config/config.yaml" && go run ./cmd/subscriber/main.go

