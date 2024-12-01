DB_URL = postgresql://root:secret@localhost:5432/membership_db?sslmode=disable
postgres:
	docker run --name membership-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it membership-db createdb --username=root --owner=root membership_db

dropdb:
	docker exec -it membership-db dropdb membership_db

new_migration:
	migrate create -ext sql -seq -dir internal/db/migrations $(name)

migrate_up:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down 1

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb new_migration migrate_up migrate_down