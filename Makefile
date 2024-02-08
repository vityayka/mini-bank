migrate_up:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down\
	
migrate_up_one:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down_one:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go bank/db/sqlc Store