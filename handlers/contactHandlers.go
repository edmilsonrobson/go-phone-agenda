package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/edmilsonrobson/go-phone-agenda/models"
	"github.com/edmilsonrobson/go-phone-agenda/repositories"
	"github.com/go-chi/chi/v5"
)

var contactRepository = repositories.ContactRepository{}

func ListContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	contacts := contactRepository.List()
	json.NewEncoder(w).Encode(contacts)
}

func AddContact(w http.ResponseWriter, r *http.Request) {
	var c models.Contact
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contactRepository.Add(&c)
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	var jsonRequest map[string]json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&jsonRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var updatedContact models.Contact
	var contactName string
	json.Unmarshal(jsonRequest["name"], &contactName)
	json.Unmarshal([]byte(jsonRequest["updatedContact"]), &updatedContact)
	success := contactRepository.Update(contactName, &updatedContact)
	if !success {
		http.Error(w, "Could not update", http.StatusBadRequest)
	}
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	var jsonRequest map[string]string
	err := json.NewDecoder(r.Body).Decode(&jsonRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	contactName := jsonRequest["name"]

	success := contactRepository.Remove(contactName)
	if !success {
		http.Error(w, "Could not delete", http.StatusBadRequest)
	}
}

func SearchContactById(w http.ResponseWriter, r *http.Request) {
	rawContactId := chi.URLParam(r, "contactId")
	contactId, err := strconv.Atoi(rawContactId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c := contactRepository.FindById(contactId)
	json.NewEncoder(w).Encode(c)
}

func SearchContactByName(w http.ResponseWriter, r *http.Request) {
	contactName := chi.URLParam(r, "contactName")

	c := contactRepository.FindByName(contactName)
	json.NewEncoder(w).Encode(c)
}
