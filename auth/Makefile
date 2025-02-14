test:
	go test -v -cover ./...

postgres:
	docker run --name authDB -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres:17.2

createdb:
	docker exec -it authDB createdb -U postgres auth

dropdb:
	docker exec -it authDB dropdb auth

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@host.docker.internal:5432/auth?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/auth?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

.PHONY: postgres createdb test dropdb migrateup migratedown sqlc server mock
