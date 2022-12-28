DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable


postgres:
	docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:13-alpine

createdb:
	docker exec -it postgres13 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres13 dropdb simple_bank

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/hmoodallahma/123bank/db/sqlc Store

server:
	go run main.go

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown mock server sqlc