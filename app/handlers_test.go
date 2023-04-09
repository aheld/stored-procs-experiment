package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aheld/listservice/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type fakeDB struct{ mock.Mock }

func (d *fakeDB) GetListItems(bannerId string, userId int) ([]domain.ListItem, error) {
	return []domain.ListItem{}, nil
}
func (d *fakeDB) UpdateListItem(bannerId string, userId int, itemId int, item string) error {
	return nil
}
func (d *fakeDB) CheckVersion() (string, error) { return "good to go", nil }
func (d *fakeDB) InsertListItem(bannerId string, userId int, item string) (int, error) {
	args := d.Called(userId, item)
	return args.Int(0), args.Error(1)
}

func TestAlive(t *testing.T) {
	testDB := &fakeDB{}
	s := CreateNewServer(testDB)
	s.MountInfrastructureHandlers()
	req, _ := http.NewRequest("GET", "/startz", nil)
	response := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, response.Code)
	require.Equal(t, "good to go", response.Body.String())
}

func TestInsert(t *testing.T) {
	item := "Test List Item"

	testDB := new(fakeDB)
	testDB.On("InsertListItem", 100, item).Return(150, nil)
	testDB.On("InsertListItem", 0, item).Return(150, nil)

	s := CreateNewServer(testDB)
	s.MountHandlers()

	response := executeListCreatePost(s, 100, item)
	checkResponseCode(t, http.StatusCreated, response.Code)
	checkResponseBody(t, response, item)

	response = executeListCreatePost(s, 0, item)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	response = executeListCreatePost(s, 5, item)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func checkResponseBody(t *testing.T, response *httptest.ResponseRecorder, item string) {
	var responseJson map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &responseJson))
	require.Equal(t, item, responseJson["item"])
}

func executeListCreatePost(s *Server, userId int, item string) *httptest.ResponseRecorder {
	body := ""
	switch userId {
	case 0:
		body = `{}`
	case 100:
		body = fmt.Sprintf(`{"user_id": 100, "item": "%s"}`, item)
	case 5:
		body = fmt.Sprintf(`{"item": %s}`, item)
	}
	jsonBody := []byte(body)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/lists", bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return executeRequest(req, s)
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
