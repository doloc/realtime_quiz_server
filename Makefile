postgres:
	docker run --name postgreSQL16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:16

createdb:
	docker exec -it postgreSQL16 createdb --username=root --owner=root realtime_quiz

dropdb:
	docker exec -it postgreSQL16 dropdb realtime_quiz

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/realtime_quiz?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/realtime_quiz?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb sqlc test server