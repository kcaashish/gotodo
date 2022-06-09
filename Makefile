include .env

.PHONY: postgres adminer migrate_up migrate_down server

postgres:
	sudo docker run --rm -ti --network host -e POSTGRES_PASSWORD=$(PG_PASSWORD) postgres

adminer:
	sudo docker run --rm -ti --network host adminer

migrate_up:
	migrate -source file://migrations -database $(DB_URI) up

migrate_down:
	migrate -source file://migrations -database $(DB_URI) down

server:
	go run cmd/gotodo/main.go