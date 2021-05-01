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

	if contacts == nil {
		return []models.Contact{}
	}
	return contacts
}

func (r *ContactRepository) Update(contactName string, updatedContact *models.Contact) bool {
	redisConn, err := redis.Dial("tcp", ":6380")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer redisConn.Close()

	serializedContact, err := json.Marshal(*updatedContact)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// Name is the same, so just update the value of that same key, because the name is used as the key
	if updatedContact.Name == contactName {
		_, err = redisConn.Do("HSET", "contacts", contactName, serializedContact)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	} else {
		// Name has changed. Delete the old record and create it under a new key, since names are used for keys
		_, err = redisConn.Do("HSET", "contacts", updatedContact.Name, serializedContact)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		_, err = redisConn.Do("HDEL", "contacts", contactName, serializedContact)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

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
	redisConn, err := redis.Dial("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	defer redisConn.Close()

	redisReturn, err := redis.Bool(redisConn.Do("HDEL", "contacts", contactName))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return redisReturn
}

func (r *ContactRepository) FindById(contactId int) models.Contact {

	return models.Contact{}
}

func (r *ContactRepository) FindByName(contactName string) []models.Contact {
	return []models.Contact{}
}
