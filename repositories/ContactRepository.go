package repositories

import (
	"github.com/edmilsonrobson/go-phone-agenda/models"
)

type ContactRepository struct{}

var phoneBook = []models.Contact{
	{
		Name:  "Ed",
		Phone: "+5508511111111",
	},
	{
		Name:  "Santa Claus",
		Phone: "+5508522222222",
	},
}

func (r *ContactRepository) List() []models.Contact {
	return phoneBook
}

func (r *ContactRepository) Update(contactId int, c *models.Contact) bool {
	return true
}

func (r *ContactRepository) Add(c *models.Contact) bool {
	phoneBook = append(phoneBook, *c)
	return true
}

func (r *ContactRepository) Remove(contactId int) bool {
	return true
}

func (r *ContactRepository) FindById(contactId int) models.Contact {
	return models.Contact{}
}

func (r *ContactRepository) FindByName(contactName string) []models.Contact {
	return []models.Contact{}
}
