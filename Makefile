include .env

postgres_migrate_up:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB}?sslmode=disable -path postgres up

postgres_migrate_down:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB}?sslmode=disable -path postgres down

scylla_migrate_up:
	@migrate -database cassandra://${SCYLLA_LISTEN_ADDRESS}:${SCYLLA_PORT}/${SCYLLA_KEYSPACE} -path scylla up

scylla_migrate_down:
	@migrate -database cassandra://${SCYLLA_LISTEN_ADDRESS}:${SCYLLA_PORT}/${SCYLLA_KEYSPACE} -path scylla down

build:
	@go build -o bin/main

run: build
	@./bin/main

test:
	@go test -v ./...


