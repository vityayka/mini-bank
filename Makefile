migrate_up:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up 1

migrate_down:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go bank/db/sqlc Store