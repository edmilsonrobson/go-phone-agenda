package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestListContacts(t *testing.T) {
	routes := Routes()
	routeToTest := "/contacts"

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + routeToTest)
	if err != nil {
		t.Log(err)
	}

	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type of 'application/json', but got %v", resp.Header.Get("Content-Type"))
	}
}

func TestAddContacts(t *testing.T) {
	routes := Routes()
	routeToTest := "/contacts"

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	values, err := json.Marshal(map[string]string{
		"name":  "Ryan",
		"phone": "+18855494912",
	})

	if err != nil {
		log.Fatal(err)
	}

	resp, err := ts.Client().Post(ts.URL+routeToTest, "application/json", bytes.NewBuffer(values))
	if err != nil {
		t.Log(err)
	}

	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type of 'application/json', but got %v", resp.Header.Get("Content-Type"))
	}

	teardown()
}

func TestDeleteContact(t *testing.T) {
	routes := Routes()
	routeToTest := "/contacts"

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	values, err := json.Marshal(map[string]string{
		"name": "Ryan",
	})

	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("DELETE", ts.URL+routeToTest, bytes.NewBuffer(values))
	if err != nil {
		t.Log(err)
	}

	resp, err := ts.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}

	expectedStatusCode := http.StatusBadRequest
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %v, but got %v", expectedStatusCode, resp.StatusCode)
	}

	teardown()
}

func TestUpdateContact(t *testing.T) {
	routes := Routes()
	routeToTest := "/contacts/update"

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	type UpdatedContact struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	type UpdateContactBody struct {
		Name           string         `json:"name"`
		UpdatedContact UpdatedContact `json:"updatedContact"`
	}

	testData := UpdateContactBody{
		Name: "Ryan",
		UpdatedContact: UpdatedContact{
			Name:  "John",
			Phone: "444",
		},
	}

	values, err := json.Marshal(testData)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := ts.Client().Post(ts.URL+routeToTest, "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}

	expectedStatusCode := http.StatusBadRequest
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %v, but got %v", expectedStatusCode, resp.StatusCode)
	}

	teardown()
}

func TestSearchContactByName(t *testing.T) {
	routes := Routes()
	routeToTest := "/contacts/search"

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	values := url.Values{}
	values.Add("name", "Joker")

	url := ts.URL + routeToTest + "?" + values.Encode()
	resp, err := ts.Client().Get(url)
	if err != nil {
		log.Fatal(err)
	}

	expectedStatusCode := http.StatusNoContent
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %v, but got %v", expectedStatusCode, resp.StatusCode)
	}

	teardown()
}

func teardown() {
	redisConn, err := redis.Dial("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	redisConn.Do("FLUSHALL")
}
