migrate_up:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...