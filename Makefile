.PHONY: postgres adminer migrate_up migrate_force migrate_down

postgres:
	sudo docker run --rm -ti --network host -e POSTGRES_PASSWORD=secret postgres

adminer:
	sudo docker run --rm -ti --network host adminer

migrate_up:
	migrate -source file://migrations \
	-database postgres://postgres:secret@localhost/postgres?sslmode=disable up

migrate_force:
	migrate -path migrations \
	-database postgres://postgres:secret@localhost/postgres?sslmode=disable force 1

migrate_down:
	migrate -source file://migrations \
	-database postgres://postgres:secret@localhost/postgres?sslmode=disable down