include .env

postgres_migrate_up:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB} -path postgres up

postgres_migrate_down:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}${POSTGRES_DB} -path postgres down

scylla_migrate_up:
	@migrate -database "cassandra://node-0.gce-us-east-1.235abdc69354fd53a8ad.clusters.scylla.cloud:9042/ticketing_api?username=scylla&password=n2Lhcj3fA7JbTtN" -path scylla up

scylla_migrate_down:
	@migrate -database cassandra://${SCYLLA_LISTEN_ADDRESS}:${SCYLLA_PORT}/${SCYLLA_KEYSPACE} -path scylla down

build:
	@go build -o bin/main

run: build
	@./bin/main

test:
	@go test -v ./...


