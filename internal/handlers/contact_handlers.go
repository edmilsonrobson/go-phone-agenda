package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/edmilsonrobson/go-phone-agenda/internal/models"
	"github.com/edmilsonrobson/go-phone-agenda/internal/repositories"
)

var contactRepository = repositories.ContactRepository{}

func ListContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	contacts := contactRepository.List()
	json.NewEncoder(w).Encode(contacts)
}

func validateContact(contact *models.Contact) bool {
	fields := []string{contact.Name, contact.Phone}

	for _, v := range fields {
		if v == "" {
			return false
		}
	}

	return true
}

func AddContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c models.Contact
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateContact(&c) {
		http.Error(w, "Cannot add duplicate contacts", http.StatusBadRequest)
		return
	}
	success := contactRepository.Add(&c)
	if !success {
		http.Error(w, "Cannot add duplicate contacts", http.StatusBadRequest)
	}
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

	if !validateContact(&updatedContact) {
		http.Error(w, "Cannot add duplicate contacts", http.StatusBadRequest)
		return
	}
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
		http.Error(w, "Could not delete requested contact", http.StatusBadRequest)
	}
}

func SearchContactByName(w http.ResponseWriter, r *http.Request) {
	contactName := r.URL.Query().Get("name")

	c := contactRepository.FindByName(contactName)
	if c != (models.Contact{}) {
		json.NewEncoder(w).Encode(c)
	} else {
		http.Error(w, "No contacts found", http.StatusNoContent)
	}

}
