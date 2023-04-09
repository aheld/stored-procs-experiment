package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestSql(t *testing.T) {

	const banner_id = "f4bd6cdc-eb4b-4f74-8565-c243d3fdf20c"
	const user_id = 150

	item := faker.Sentence()

	db, _ := Initialize()

	// Insert a row
	list_id := -1
	db.Conn.QueryRow("select * from list_items_insert($1, $2, $3)", banner_id, user_id, item).Scan(&list_id)
	assert.Greater(t, list_id, 0, "This should be a positive integer")
	fmt.Printf("List ID is %d", list_id)
	fmt.Printf("select * from list_items_insert(%v, %v, %v)", banner_id, user_id, item)

	// Assert the row was inserted
	itemFromDB := getListIdAndItem(t, db, banner_id, user_id, list_id)
	assert.Equal(t, item, itemFromDB, "Item from db should match %s", item)

	newItem := "Milkly Way Bar"
	db.Conn.Exec("call list_items_update($1, $2, $3, $4);", banner_id, user_id, list_id, newItem)
	itemFromDB = getListIdAndItem(t, db, banner_id, user_id, list_id)
	assert.Equal(t, itemFromDB, newItem, "Item from db should match %s", newItem)

}

func getListIdAndItem(t *testing.T, db Database, banner_id string, user_id int, list_id int) string {
	itemFromDB := ""
	err := db.Conn.QueryRow("SELECT user_text from list_items_view where banner_id=$1 and user_id=$2 and id=$3",
		banner_id,
		user_id,
		list_id).Scan(&itemFromDB)
	if err != nil {
		assert.NotNil(t, err, "DB Failed to query")
	}
	return itemFromDB
}

type Database struct {
	Conn *sql.DB
}

func Initialize() (Database, error) {
	connStr := os.Getenv("POSTGRESQL_URL")
	if connStr == "" {
		log.Fatal("No connection string defined, exiting")
	}
	log.Println(connStr)

	db := Database{}
	log.Printf("Connected using:  %v", connStr)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	return db, nil
}
