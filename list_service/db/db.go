package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aheld/listservice/domain"
	_ "github.com/lib/pq"
)

// Make sure the DB has the tables and functions we expect
const SCHEMA_VERSION_REQUIRED = 2

type Database struct {
	Conn *sql.DB
}

func Initialize(connStr string) (Database, error) {
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

func (db Database) InsertListItem(userId int, item string) (int, error) {
	rows, err := db.Conn.Query("select * from list_items_insert('f4bd6cdc-eb4b-4f74-8565-c243d3fdf20a', $1, $2)", userId, item)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		var listId int
		err := rows.Scan(&listId)
		if err != nil {
			return listId, err
		}
	}
	return -1, fmt.Errorf("no item inserted")
}

func (db Database) GetListItems(userId int) ([]domain.ListItem, error) {
	rows, err := db.Conn.Query(
		"SELECT id, user_text from list_items_view where banner_id='f4bd6cdc-eb4b-4f74-8565-c243d3fdf20a' and user_id=$1 ",
		userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.ListItem{}
	for rows.Next() {
		var i domain.ListItem
		if err := rows.Scan(&i.Id, &i.Item); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (db Database) UpdateListItem(userId int, itemId int, item string) error {
	_, err := db.Conn.Exec("call list_items_update('f4bd6cdc-eb4b-4f74-8565-c243d3fdf20a', $1, $2, $3);", userId, itemId, item)
	return err
}

func (db Database) CheckVersion() (string, error) {
	var version int
	if err := db.Conn.QueryRow("select max(version) from schema_migrations where dirty=false").Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no schema version found")
		}
		return "", err
	}
	if version >= SCHEMA_VERSION_REQUIRED {
		return fmt.Sprintf("Schema Version is %d, which is good, because we need version %d", version, SCHEMA_VERSION_REQUIRED), nil
	}
	return "", fmt.Errorf("required schema version not found")
}
