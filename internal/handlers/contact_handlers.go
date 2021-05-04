package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/edmilsonrobson/go-phone-agenda/internal/models"
)

type ContactRepository interface {
	List() []models.Contact
	Add(*models.Contact) bool
	Update(string, *models.Contact) bool
	Remove(string) bool
	SearchByName(string) *models.Contact
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

func ListContacts(repo ContactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		contacts := repo.List()
		json.NewEncoder(w).Encode(contacts)
	}
}

func AddContact(repo ContactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var c models.Contact
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if !validateContact(&c) {
			http.Error(w, "Cannot add duplicate contacts", http.StatusBadRequest)
		} else {
			success := repo.Add(&c)
			if success {
				w.WriteHeader(200)
			} else {
				http.Error(w, "Cannot add duplicate contacts", http.StatusBadRequest)
			}
		}
	}
}

func UpdateContact(repo ContactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		} else {
			success := repo.Update(contactName, &updatedContact)
			if success {
				w.WriteHeader(200)
			} else {
				http.Error(w, "Could not update", http.StatusBadRequest)
			}
		}
	}
}

func DeleteContact(repo ContactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonRequest map[string]string
		err := json.NewDecoder(r.Body).Decode(&jsonRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		contactName := jsonRequest["name"]

		success := repo.Remove(contactName)
		if success {
			w.WriteHeader(200)
		} else {
			http.Error(w, "Could not delete requested contact", http.StatusBadRequest)
		}
	}
}

func SearchContactByName(repo ContactRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contactName := r.URL.Query().Get("name")

		c := repo.SearchByName(contactName)
		if *c != (models.Contact{}) {
			json.NewEncoder(w).Encode(c)
		} else {
			http.Error(w, "No contacts found", http.StatusNoContent)
		}
	}
}
