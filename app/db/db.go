package db

import (
	"database/sql"
	"fmt"
	"github.com/aheld/listservice/domain"
	_ "github.com/lib/pq"
	"log"
)

// Make sure the DB has the tables and functions we expect
const SCHEMA_VERSION_REQUIRED = 2

type Database struct {
	Conn          *sql.DB
	SchemaVersion int
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
	err = db.SetVersion()
	if err != nil {
		return db, err
	}
	db.CheckVersion()
	log.Printf("Database version %d", db.SchemaVersion)
	return db, nil
}

func (db Database) InsertListItem(bannerId string, userId int, item string) (int, error) {
	listId := -1
	err := db.Conn.QueryRow("select * from list_items_insert($1, $2, $3)",
		bannerId,
		userId,
		item).Scan(&listId)
	if err != nil {
		return -1, err
	}
	return listId, nil
}

func (db Database) GetListItems(bannerId string, userId int) ([]domain.ListItem, error) {
	switch db.SchemaVersion {
	case 3:
		return db.getListItems3(bannerId, userId)
	default:
		return db.getListItems2(bannerId, userId)
	}
}

func (db Database) getListItems2(bannerId string, userId int) ([]domain.ListItem, error) {
	rows, err := db.Conn.Query(
		"SELECT id, user_text from list_items_view where banner_id=$1 and user_id=$2;",
		bannerId,
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

func (db Database) getListItems3(bannerId string, userId int) ([]domain.ListItem, error) {
	rows, err := db.Conn.Query(
		"SELECT id, item from list_items_get($1,$2);",
		bannerId,
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

func (db Database) GetListItem(bannerId string, userId int, listId int) (domain.ListItem, error) {
	l := &domain.ListItem{}
	err := db.Conn.QueryRow("SELECT id, user_id, user_text from list_items_view where banner_id=$1 and user_id=$2 and id=$3",
		bannerId,
		userId,
		listId).Scan(&l.Id, &l.UserId, &l.Item)
	if err != nil {
		return *l, err
	}
	return *l, nil
}

func (db Database) UpdateListItem(bannerId string, userId int, itemId int, item string) error {
	_, err := db.Conn.Exec("call list_items_update($1, $2, $3, $4);", bannerId, userId, itemId, item)
	return err
}

func (db *Database) SetVersion() error {
	var version int
	if err := db.Conn.QueryRow(
		"select max(version) from schema_migrations where dirty=false",
	).Scan(&version); err != nil {
		return err
	}
	db.SchemaVersion = version
	return nil
}

func (db Database) CheckVersion() (string, error) {
	//version, err := db.GetVersion()
	//if err != nil {
	//	return "", err
	//}
	if db.SchemaVersion >= SCHEMA_VERSION_REQUIRED {
		return fmt.Sprintf("Schema Version is %d, which is good, because we need version %d", db.SchemaVersion, SCHEMA_VERSION_REQUIRED), nil
	}
	return "", fmt.Errorf("required schema version not found")
}
