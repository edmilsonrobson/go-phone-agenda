package handlers

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("Update contact")
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	rawContactId := chi.URLParam(r, "contactId")
	contactId, err := strconv.Atoi(rawContactId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	contactRepository.Remove(contactId)
}

func SearchContactById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Search contact (By ID)")
}

func SearchContactByName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Search contact (by name)")
}
