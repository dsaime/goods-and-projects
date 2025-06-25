.PHONY: test vet lint check migrate-clickhouse compose-up run-server run-event-listener

# Get number of CPU cores minus 1 for parallel execution
CORES := $(shell echo $$(( $$(nproc) - 1 )))

test:
	go test ./...

vet:
	go vet ./...

lint:
	echo "Running golangci-lint with $(CORES) workers..."
# The -j parameter for golangci-lint will use all available CPU cores minus one (to avoid overloading your system)
	golangci-lint run -v -j $(CORES)

check: vet lint

migrate-clickhouse:
	migrate -database ${CLICKHOUSE_DSN} -path ./migrations/clickhouse up

compose-up:
	docker compose -f docker/docker-compose.yaml up

run-server:
	go run ./cmd/goods-and-projects --http-addr :8080 --redis-url redis://localhost:6379 --pgsql-dsn postgresql://postgres:postgres@localhost:5432/gapdb?sslmode=disable --nats-url nats://localhost:4222

run-event-listener:
	go run ./cmd/goods-event-listener \
		--nats-url nats://localhost:4222 \
		--clickhouse-dsn http://localhost:8123?username=clickhouse&password=clickhouse