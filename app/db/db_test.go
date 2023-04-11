package db

import (
	"fmt"
	"github.com/aheld/listservice/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

const bannerId = "f4bd6cdc-eb4b-4f74-8565-c243d3fdf20c"
const userId = 150

func setupSuite(t *testing.T) (func(t *testing.T), Database) {
	log.Println("setup DB connection")

	connStr := os.Getenv("POSTGRESQL_URL")
	if connStr == "" {
		fmt.Print("No connection string defined, skipping DB tests")
	}
	fmt.Printf("Running DB Tests against %s", connStr)
	db, _ := Initialize(connStr)
	// Return a function to teardown the test
	return func(t *testing.T) {
		log.Println("teardown suite")
		db.Conn.Close()
	}, db
}

func TestSqlSchema2(t *testing.T) {
	teardownSuite, db := setupSuite(t)
	defer teardownSuite(t)

	if db.SchemaVersion < 2 {
		t.Skip()
	}
	nowTime := time.Now().String() //just to have a unique string
	item := "I need Orgeno at %s"

	// Insert a row
	listId, err := db.InsertListItem(bannerId, userId, item)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("List Id of inserted row is %d", listId)
	require.Greater(t, listId, 0, "This should be a positive integer!")

	// Assert the row was inserted
	listItem, err := db.GetListItem(bannerId, userId, listId)
	if err != nil {
		fmt.Println(err)
	}

	assert.Nil(t, err, "DB Failed to get a single item")
	require.Equal(t, item, listItem.Item, "Item from db should match %s", item)

	//Update Item and assert changes
	newItem := fmt.Sprintf("Milkly Way Bar @ %s", nowTime)
	db.UpdateListItem(bannerId, userId, listId, newItem)
	// check what was just inserted
	listItem, err = db.GetListItem(bannerId, userId, listId)
	assert.Nil(t, err, "DB Failed to get a single item")
	assert.Equal(t, newItem, listItem.Item, "Item from db should match")
}

func TestSqlSchema3(t *testing.T) {
	teardownSuite, db := setupSuite(t)
	defer teardownSuite(t)

	if db.SchemaVersion < 3 {
		t.Skip()
	}
	nowTime := time.Now().String() //just to have a unique string
	item := fmt.Sprintf("I need Orgeno at %s", nowTime)

	// Insert a row
	_, err := db.InsertListItem(bannerId, userId, item)
	assert.Nil(t, err, "DB Failed to insert")

	items, err := db.GetListItems(bannerId, userId)
	assert.Nil(t, err, "DB Failed to get a all items")

	assert.Equal(t, item, findItem(items, item), "Failed to retrive our new item")
}

func findItem(items []domain.ListItem, item string) string {
	for _, it := range items {
		if item == it.Item {
			return it.Item
		}
	}
	return ""
}
