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

	var contacts []models.Contact
	contactStringMap, err := redis.StringMap(redisConn.Do("HGETALL", "contacts"))
	for _, v := range contactStringMap {
		var contact models.Contact
		err = json.Unmarshal([]byte(v), &contact)
		if err != nil {
			fmt.Println(err.Error())
			return []models.Contact{}
		}
		contacts = append(contacts, contact)
	}
	if err != nil {
		fmt.Println(err.Error())
		return []models.Contact{}
	}

	return contacts
}

func (r *ContactRepository) Update(contactId int, c *models.Contact) bool {
	return true
}

func (r *ContactRepository) Add(c *models.Contact) bool {
	redisConn, err := redis.Dial("tcp", ":6380")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer redisConn.Close()

	serializedContact, err := json.Marshal(*c)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = redisConn.Do("HSET", "contacts", c.Name, serializedContact)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func (r *ContactRepository) Remove(contactName string) bool {
	/* redisConn, err := redis.Dial("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	serializedContact, err := json.Marshal(*c)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = redisConn.Do("LREM", "contacts", serializedContact)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} */

	return true
}

func (r *ContactRepository) FindById(contactId int) models.Contact {
	return models.Contact{}
}

func (r *ContactRepository) FindByName(contactName string) []models.Contact {
	return []models.Contact{}
}
