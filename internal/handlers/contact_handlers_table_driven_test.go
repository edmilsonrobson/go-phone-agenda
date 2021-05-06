package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/edmilsonrobson/go-phone-agenda/internal/models"
	"github.com/edmilsonrobson/go-phone-agenda/internal/repositories"
	"github.com/gomodule/redigo/redis"
	"github.com/google/go-cmp/cmp"
)

var testCases = []struct {
	name               string
	route              string
	method             string
	body               map[string]interface{}
	expectedContacts   []models.Contact
	expectedStatusCode int
}{
	{
		name:   "list base contacts",
		route:  "/contacts",
		method: "GET",
		body:   map[string]interface{}{},
		expectedContacts: []models.Contact{
			{Name: "Ryan", Phone: "456789"},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "add contact",
		route:  "/contacts",
		method: "POST",
		body: map[string]interface{}{
			"name":  "Santa Claus",
			"phone": "444879",
		},
		expectedContacts: []models.Contact{
			{Name: "Ryan", Phone: "456789"},
			{Name: "Santa Claus", Phone: "444879"},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "add contact with the same name as existing contact",
		route:  "/contacts",
		method: "POST",
		body: map[string]interface{}{
			"name":  "Santa Claus",
			"phone": "444879",
		},
		expectedContacts: []models.Contact{
			{Name: "Ryan", Phone: "456789"},
			{Name: "Santa Claus", Phone: "444879"},
		},
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:   "delete contact that exists",
		route:  "/contacts",
		method: "DELETE",
		body: map[string]interface{}{
			"name": "Santa Claus",
		},
		expectedContacts: []models.Contact{
			{Name: "Ryan", Phone: "456789"},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "delete contact that doesn't exist",
		route:  "/contacts",
		method: "DELETE",
		body: map[string]interface{}{
			"name": "Santa Claus",
		},
		expectedContacts: []models.Contact{
			{Name: "Ryan", Phone: "456789"},
		},
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:   "update contact that exists",
		route:  "/contacts/update",
		method: "POST",
		body: map[string]interface{}{
			"name": "Ryan",
			"updatedContact": map[string]string{
				"name":  "Kratos",
				"phone": "0",
			},
		},
		expectedContacts: []models.Contact{
			{Name: "Kratos", Phone: "0"},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "update contact that doesn't exist",
		route:  "/contacts/update",
		method: "POST",
		body: map[string]interface{}{
			"name": "Cloud Strife",
			"updatedContact": map[string]string{
				"name":  "The Ashen One",
				"phone": "123",
			},
		},
		expectedContacts: []models.Contact{
			{Name: "Kratos", Phone: "0"},
		},
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:   "delete contact to clear list",
		route:  "/contacts",
		method: "DELETE",
		body: map[string]interface{}{
			"name": "Kratos",
		},
		expectedContacts:   []models.Contact{},
		expectedStatusCode: http.StatusOK,
	},
}

var repo ContactRepository

func TestMain(m *testing.M) {
	redisConn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err)
	}

	redisConn.Do("FLUSHALL")

	// Test with redis
	repo = repositories.NewRedisContactRepository(&redisConn)
	setupInitialStorage()
	m.Run()

	// Test in memory
	repo = repositories.NewInMemoryContactRepository()
	setupInitialStorage()
	code := m.Run()

	os.Exit(code)
}

func setupInitialStorage() {
	repo.Add(&models.Contact{
		Name:  "Ryan",
		Phone: "456789",
	})
}

func TestHandlersOnExistingData(t *testing.T) {
	routes := getTestRoutes(repo)
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, testCase := range testCases {
		switch testCase.method {
		case "GET":
			resp, err := ts.Client().Get(ts.URL + testCase.route)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != testCase.expectedStatusCode {
				t.Errorf("For test '%s', expected status code %d but got %d", testCase.name, testCase.expectedStatusCode, resp.StatusCode)
			}
		case "POST":
			payload, err := json.Marshal(testCase.body)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := ts.Client().Post(ts.URL+testCase.route, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != testCase.expectedStatusCode {
				t.Errorf("For test '%s', expected status code %d but got %d", testCase.name, testCase.expectedStatusCode, resp.StatusCode)
			}
		case "DELETE":
			payload, err := json.Marshal(testCase.body)
			if err != nil {
				t.Fatal(err)
			}
			request, err := http.NewRequest("DELETE", ts.URL+"/contacts", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatal(err)
			}
			resp, err := ts.Client().Do(request)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != testCase.expectedStatusCode {
				t.Errorf("For test '%s', expected status code %d but got %d", testCase.name, testCase.expectedStatusCode, resp.StatusCode)
			}
		default:
			t.Errorf("Invalid method type: %s", testCase.method)
		}

		resp, err := ts.Client().Get(ts.URL + "/contacts")
		if err != nil {
			t.Fatal(err)
		}

		contactsReceived := []models.Contact{}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		json.Unmarshal(respBody, &contactsReceived)
		contactsMatch := cmp.Equal(contactsReceived, testCase.expectedContacts)
		if !contactsMatch {
			t.Errorf("For test '%s', contacts don't match. Expected %v, but got %v", testCase.name, testCase.expectedContacts, contactsReceived)
		}
		resp.Body.Close()
	}
}
