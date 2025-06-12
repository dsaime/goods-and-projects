.PHONY: test vet lint check run migrate

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

run:
	go run ./cmd/app/main.go

migrate-pgsl:
	migrate -database ${PGSQL_DSN} -path ./migrations/pgsql up

migrate-clickhouse:
	migrate -database ${CLICKHOUSE_DSN} -path ./migrations/clickhouse up