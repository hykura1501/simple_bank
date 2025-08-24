postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:alpine

createdb: 
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it postgres dropdb simple_bank

startdb:
	docker start postgres

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down 1


sqlc: 
	sqlc generate

test: 
	go test -v -cover ./...

run:
	go run main.go

mock:
	mockgen -destination=db/mock/store.go -package=mockdb github.com/hykura1501/simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test run mock startdb