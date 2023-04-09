package db

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSql(t *testing.T) {

	connStr := os.Getenv("POSTGRESQL_URL")
	if connStr == "" {
		fmt.Print("No connection string defined, skipping DB tests")
		return
	}

	log.Printf("Running DB Tests against %s", connStr)

	const bannerId = "f4bd6cdc-eb4b-4f74-8565-c243d3fdf20c"
	const userId = 150

	nowTime := time.Now().String() //just to have a unique string
	item := fmt.Sprintf("I need Orgeno at %s", nowTime)

	db, _ := Initialize(connStr)

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
