# A Fresh look at Stored Procedures!

Experiment to see how we can build use stored procedures and views to ensure all the data logic is encapsulated in PostgreSQL, while freeing the application code to only handle business logic.  

In this case the business logic will involve manipulating the list items and the data logic will ensure data partitioning across tenants, as well as common data concerns like record timestamps.

I am experimenting with doing data migrations and stored proc testing in a dedicated go application as opposed to building the tests directing in the main application.

## "Shopping List" 
a variation of the classic ToDo example app.  Start with a list of text items and evolve to more complex types such as recipes with ingredients.

Every list item has a simple model with an id, user_id and item text.  Its a multi-tenant application as well, so the database data model has an `tenant_id` which the application model does not. 

There 
## Rest Application
`app` directory contains a simple go app.  The interesting files are `db/db_test.go` and `db/db.go` which connect the database.

Its a multi-tenant application, but the domain and app logic doesn't know anything about tenants.  It just plucks a tenantId out of the middleware and passes it to the database.

The Application uses a /startup probe to ensure that it will not take traffic until the database it is connected to has the schema version that the application expects.  I will use this for rolling kubernetes deployments to let multiple versions of the be pushed out independently of the schema updates. 

```# Startup will return 200 when the database had the proper schema
### Startup Probe

http://localhost:3333/startz


### New ListItem Post

POST http://localhost:3333/lists HTTP/1.1
content-type: application/json

{
    "user_id": 100,
    "item": "Oregano"
}


### GET all list items for a user

http://localhost:3333/lists/100


### Update a specific list item for a user

PUT http://localhost:3333/lists/100/1
content-type: application/json

{
    "item": "Thyme & Oregano & Rosemary"
}
```

## Citus Distributed PostgresQL
Using [Citus PostgreSQL](https://www.citusdata.com/download/), a distributed database that has very strong row level partitioning.  It is a great data foundation for a multi-tenant application.

## Migration Application
[go-Migrate](https://github.com/golang-migrate/migrate) Is used to modify the DB schema.  The `Makefile` drives command like migrations and the go app under `db/cmd/` programmatically applies migrations and tests them.  There is also a docker based migration runner, which my be the right answer for kubernetes deployments. 


# Development

## Setup

You need to have go and should have the postgreSQL client installed, or some sql admin tool.  I run citus postgres as a docker image, but you can also choose to the the native version.

If you use your own postgreSQL, edit the `POSTGRESQL_URL` variable in the Makefile.

### Makefile helpers

### `make boostrap`
Will install postgreSQL (on a mac)
Install `golang-migrate` as a CLI tool (on a mac)
Download the citus distributed docker image
Download the golang dependencies for the rest app.

### `make start_pg`
Runs the Docker version of Citus postgres

### `run_migrations`
Runs migrations up to the latest version.  If it fails you will have to manually repair the database and then run `migrate force` to clean it up.

### `run_migrations_down`
Backs out the current migration.

### `make run_go`
Runs the go rest application.  If you want to run it via CLI or your editor, you will need the `POSTGRESQL_URL` in your environment.

### `make db_test`
Runs the go db migration application.  If you want to run it via CLI or your editor, you will need the `POSTGRESQL_URL` in your environment.
It will run the migrations and then some tests against the DB.
