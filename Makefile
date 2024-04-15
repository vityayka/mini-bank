migrate_up:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down
	
migrate_up_one:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down_one:
	migrate -path db/migration -database "postgres://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -cover -short ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go bank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank \
    proto/*.proto

.PHONY: proto