package handlers

import (
	"fmt"
	"net/http"
)

func ListContacts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List contacts")
}

func AddContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add contact")
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update contact")
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete contact")
}

func SearchContactById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Search contact (By ID)")
}

func SearchContactByName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Search contact (by name)")
}
