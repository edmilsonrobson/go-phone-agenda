package repositories

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/edmilsonrobson/go-phone-agenda/models"
	"github.com/gomodule/redigo/redis"
)

type ContactRepository struct{}

func (r *ContactRepository) List() []models.Contact {
	redisConn, err := redis.Dial("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	contactBytes, err := redis.ByteSlices(redisConn.Do("LRANGE", "contacts", 0, -1))
	if err != nil {
		fmt.Println(err.Error())
		return []models.Contact{}
	}

	var contacts []models.Contact
	for _, v := range contactBytes {
		var contact models.Contact
		err = json.Unmarshal(v, &contact)
		if err != nil {
			fmt.Println(err.Error())
			return []models.Contact{}
		}
		contacts = append(contacts, contact)
	}

	return contacts
}

func (r *ContactRepository) Update(contactId int, c *models.Contact) bool {
	return true
}

func (r *ContactRepository) Add(c *models.Contact) bool {
	redisConn, err := redis.Dial("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	serializedContact, err := json.Marshal(*c)
	if err != nil {
		log.Fatal(err)
	}

	_, err = redisConn.Do("RPUSH", "contacts", serializedContact)
	if err != nil {
		log.Fatal(err)
	}

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
