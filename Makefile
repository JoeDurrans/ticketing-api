include .env

postgres_migrate_up:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB} -path postgres up

postgres_migrate_down:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB} -path postgres down

scylla_migrate_up:
	@migrate -database "cassandra://${SCYLLA_LISTEN_ADDRESS}:${SCYLLA_PORT}/${SCYLLA_KEYSPACE} -path scylla up

scylla_migrate_down:
	@migrate -database cassandra://${SCYLLA_LISTEN_ADDRESS}:${SCYLLA_PORT}/${SCYLLA_KEYSPACE} -path scylla down

test:
	@go test -v ./tests/...
