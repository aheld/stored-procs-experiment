
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

run_go_reload:
	cd app;\
	POSTGRESQL_URL=$(POSTGRESQL_URL) air

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
	# Setup live reload for go
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'