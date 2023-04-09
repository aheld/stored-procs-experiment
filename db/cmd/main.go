package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	connStr := os.Getenv("POSTGRESQL_URL")
	if connStr == "" {
		log.Fatal("No connection string defined, exiting")
	}
	log.Println(connStr)

	runMigration(connStr, 2)
}

func runMigration(connStr string, version int) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("cound not Open postgres connection %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("cound not create migration driver %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations/",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to apply instantiate migration %v", err)
	}

	err = m.Migrate(uint(version))
	if err != nil {
		log.Fatalf("Failed to apply Migration %v", err)
	}
	ver, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("Failed to get Migration Version %v", err)
	}
	if dirty {
		log.Fatal("DIRTY Migration!! Go Clean it up")
	}
	fmt.Printf("Current DB on version %v", ver)
}
