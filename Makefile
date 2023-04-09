
POSTGRESQL_URL='postgres://postgres:mypass@localhost:5432/postgres?sslmode=disable'
PGPASSWORD ?= mypass

.PHONY = start_pg
start_pg:
	docker run -d --name citus -p 5432:5432 -e POSTGRES_PASSWORD=$(PGPASSWORD) citusdata/citus:11.2
	psql -U postgres -h localhost -d postgres -c "SELECT * FROM citus_version();"

run_migrations:
	migrate -database $(POSTGRESQL_URL) -path db/migrations up

run_migrations_down:
	echo "Down one migration"
	migrate -database $(POSTGRESQL_URL) -path db/migrations down 1

run_go:
	cd app;\
	POSTGRESQL_URL=$(POSTGRESQL_URL) go run .

db_test:
	cd db/cmd;\
	POSTGRESQL_URL=$(POSTGRESQL_URL) go run .


.PHONY = boostrap
boostrap:
	brew install postrgresql
	brew install golang-migrate
	docker pull citusdata/citus:11.2
	cd list_service; \
		go mod tidy