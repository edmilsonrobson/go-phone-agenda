package repositories

import (
	"fmt"

	"github.com/edmilsonrobson/go-phone-agenda/internal/models"
)

type InMemoryContactRepository struct {
	contacts map[string]models.Contact
}

func (r *InMemoryContactRepository) List() []models.Contact {
	list := make([]models.Contact, 0, len(r.contacts))
	fmt.Println(list)
	for _, v := range r.contacts {
		list = append(list, v)
	}

	return list
}

func (r *InMemoryContactRepository) Update(contactName string, updatedContact *models.Contact) bool {
	if _, ok := r.contacts[contactName]; ok {
		delete(r.contacts, contactName)
		r.contacts[updatedContact.Name] = *updatedContact
		return true
	}
	return false
}

func (r *InMemoryContactRepository) Add(c *models.Contact) bool {
	if _, ok := r.contacts[c.Name]; ok {
		return false
	}
	fmt.Println("Adding ", c.Name)
	r.contacts[c.Name] = *c
	return true
}

func (r *InMemoryContactRepository) Remove(contactName string) bool {
	if _, ok := r.contacts[contactName]; ok {
		delete(r.contacts, contactName)
		return true
	}
	return false
}

func (r *InMemoryContactRepository) SearchByName(contactName string) *models.Contact {
	if contact, ok := r.contacts[contactName]; ok {
		return &contact
	}
	return &models.Contact{}
}

func NewInMemoryContactRepository() *InMemoryContactRepository {
	return &InMemoryContactRepository{
		contacts: map[string]models.Contact{},
	}
}
