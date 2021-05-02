package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/edmilsonrobson/go-phone-agenda/internal/models"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type UpdatedContact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type UpdateContactBody struct {
	Name           string         `json:"name"`
	UpdatedContact UpdatedContact `json:"updatedContact"`
}

var routes http.Handler

func TestMain(m *testing.M) {
	loadTestEnv()
	routes = getTestRoutes()
	flushRedis()
	code := m.Run()
	flushRedis()
	os.Exit(code)
}

func loadTestEnv() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	path := filepath.Join(basepath, "../../.env.test")
	godotenv.Load(path)
}

func TestManipulatingContactAndCheckingTheList(t *testing.T) {
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	// Add contact
	values, err := json.Marshal(map[string]string{
		"name":  "Ryan",
		"phone": "+18855494912",
	})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := ts.Client().Post(ts.URL+"/contacts", "application/json", bytes.NewBuffer(values))
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

	// Request list of contacts to validate adding
	resp, err = ts.Client().Get(ts.URL + "/contacts")
	if err != nil {
		t.Log(err)
	}
	expectedStatusCode = http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type of 'application/json', but got %v", resp.Header.Get("Content-Type"))
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}
	expectedContacts := []models.Contact{
		{
			Name:  "Ryan",
			Phone: "+18855494912",
		},
	}
	var contacts []models.Contact
	json.Unmarshal(respBody, &contacts)
	if len(contacts) != len(expectedContacts) {
		t.Errorf("Expected %v contact, but found %v", len(expectedContacts), len(contacts))
	}
	if contacts[0] != expectedContacts[0] {
		t.Errorf("Expected first contact to be %v, but instead got %v", expectedContacts[0], contacts[0])
	}

	// Update contact that already exists

	testData := UpdateContactBody{
		Name: "Ryan",
		UpdatedContact: UpdatedContact{
			Name:  "John",
			Phone: "444",
		},
	}
	values, err = json.Marshal(testData)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = ts.Client().Post(ts.URL+"/contacts/update", "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}
	expectedStatusCode = http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %v, but got %v", expectedStatusCode, resp.StatusCode)
	}

	// Get list again to validate update

	resp, err = ts.Client().Get(ts.URL + "/contacts")
	if err != nil {
		t.Log(err)
	}
	expectedStatusCode = http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type of 'application/json', but got %v", resp.Header.Get("Content-Type"))
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}
	expectedContacts = []models.Contact{
		{
			Name:  "John",
			Phone: "444",
		},
	}

	json.Unmarshal(respBody, &contacts)
	if len(contacts) != len(expectedContacts) {
		t.Errorf("Expected %v contact, but found %v", len(expectedContacts), len(contacts))
	}
	if contacts[0] != expectedContacts[0] {
		t.Errorf("Expected first contact to be %v, but instead got %v", expectedContacts[0], contacts[0])
	}

	// Try updating a contact that doesn't exist

	testData = UpdateContactBody{
		Name: "PersonThatDoesntExist",
		UpdatedContact: UpdatedContact{
			Name:  "Jake",
			Phone: "123",
		},
	}
	values, err = json.Marshal(testData)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = ts.Client().Post(ts.URL+"/contacts/update", "application/json", bytes.NewBuffer(values))
	if err != nil {
		log.Fatal(err)
	}
	expectedStatusCode = http.StatusBadRequest
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %v, but got %v", expectedStatusCode, resp.StatusCode)
	}

	// Get list again to validate nothing changed after updating inexistent contact

	resp, err = ts.Client().Get(ts.URL + "/contacts")
	if err != nil {
		t.Log(err)
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}
	expectedContacts = []models.Contact{
		{
			Name:  "John",
			Phone: "444",
		},
	}

	json.Unmarshal(respBody, &contacts)
	if len(contacts) != len(expectedContacts) {
		t.Errorf("Expected %v contact, but found %v", len(expectedContacts), len(contacts))
	}
	if contacts[0] != expectedContacts[0] {
		t.Errorf("Expected first contact to be %v, but instead got %v", expectedContacts[0], contacts[0])
	}

	// Try deleting contact

	values, err = json.Marshal(map[string]string{
		"name": "John",
	})

	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("DELETE", ts.URL+"/contacts", bytes.NewBuffer(values))
	if err != nil {
		t.Log(err)
	}

	resp, err = ts.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}

	expectedStatusCode = http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %v, but got %v", expectedStatusCode, resp.StatusCode)
	}

	// Get List to see it has no more contacts

	resp, err = ts.Client().Get(ts.URL + "/contacts")
	if err != nil {
		t.Log(err)
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}
	expectedContacts = []models.Contact{}

	json.Unmarshal(respBody, &contacts)
	if len(contacts) != len(expectedContacts) {
		t.Errorf("Expected %v contact, but found %v", len(expectedContacts), len(contacts))
	}

	flushRedis()
}

func TestCreatingDuplicateEntries(t *testing.T) {

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Add a contact
	values, err := json.Marshal(map[string]string{
		"name":  "Ryan",
		"phone": "+18855494912",
	})

	if err != nil {
		log.Fatal(err)
	}

	resp, err := ts.Client().Post(ts.URL+"/contacts", "application/json", bytes.NewBuffer(values))
	if err != nil {
		t.Log(err)
	}

	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}

	// Try adding a contact with the same name
	values, err = json.Marshal(map[string]string{
		"name":  "Ryan",
		"phone": "+897987987",
	})

	if err != nil {
		log.Fatal(err)
	}

	resp, err = ts.Client().Post(ts.URL+"/contacts", "application/json", bytes.NewBuffer(values))
	if err != nil {
		t.Log(err)
	}

	expectedStatusCode = http.StatusBadRequest
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}

	// Check list to see if length is still 1

	resp, err = ts.Client().Get(ts.URL + "/contacts")
	if err != nil {
		t.Log(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
	}
	var contacts []models.Contact
	expectedContacts := []models.Contact{
		{
			Name:  "Ryan",
			Phone: "+18855494912",
		},
	}

	json.Unmarshal(respBody, &contacts)
	if len(contacts) != len(expectedContacts) {
		t.Errorf("Expected %v contact, but found %v", len(expectedContacts), len(contacts))
	}

	flushRedis()
}

func TestAddingIncompleteContact(t *testing.T) {

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Add a contact with no phone!
	values, err := json.Marshal(map[string]string{
		"name": "Ryan",
	})

	if err != nil {
		log.Fatal(err)
	}

	resp, err := ts.Client().Post(ts.URL+"/contacts", "application/json", bytes.NewBuffer(values))
	if err != nil {
		t.Log(err)
	}

	expectedStatusCode := http.StatusBadRequest
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}
}

func TestUpdatingIncompleteContact(t *testing.T) {

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Update a contact with no phone!
	values, err := json.Marshal(map[string]string{
		"name": "Ryan",
	})

	if err != nil {
		log.Fatal(err)
	}

	resp, err := ts.Client().Post(ts.URL+"/contacts/update", "application/json", bytes.NewBuffer(values))
	if err != nil {
		t.Log(err)
	}

	expectedStatusCode := http.StatusBadRequest
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected %d, but got %d", expectedStatusCode, resp.StatusCode)
	}
}

func TestSearchContactByName(t *testing.T) {
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

	flushRedis()
}

func flushRedis() {
	redisConn, err := redis.Dial("tcp", os.Getenv("REDIS_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	redisConn.Do("FLUSHALL")
}
