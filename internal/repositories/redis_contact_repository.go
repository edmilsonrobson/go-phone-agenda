package repositories

import (
	"encoding/json"

	"github.com/edmilsonrobson/go-phone-agenda/internal/logs"
	"github.com/edmilsonrobson/go-phone-agenda/internal/models"
	"github.com/gomodule/redigo/redis"
)

type RedisContactRepository struct {
	redisConn redis.Conn
}

func (r *RedisContactRepository) List() []models.Contact {
	var contacts []models.Contact
	contactStringMap, err := redis.StringMap(r.redisConn.Do("HGETALL", "contacts"))
	if err != nil {
		logs.WarningLogger.Printf(err.Error())
		return []models.Contact{}
	}

	for _, v := range contactStringMap {
		var contact models.Contact
		err = json.Unmarshal([]byte(v), &contact)
		if err != nil {
			logs.WarningLogger.Printf(err.Error())
			return []models.Contact{}
		}
		contacts = append(contacts, contact)
	}

	if contacts == nil {
		return []models.Contact{}
	}
	return contacts
}

func (r *RedisContactRepository) Update(contactName string, updatedContact *models.Contact) bool {
	serializedContact, err := json.Marshal(*updatedContact)
	if err != nil {
		logs.ErrorLogger.Printf(err.Error())
		return false
	}

	exists, err := redis.Bool(r.redisConn.Do("HEXISTS", "contacts", contactName))
	if err != nil {
		logs.WarningLogger.Printf(err.Error())
	}
	if !exists {
		return false
	}

	// If the name is the same, just update the value of that same key (since the name is used as the key)
	if updatedContact.Name == contactName {
		_, err = r.redisConn.Do("HSET", "contacts", contactName, serializedContact)
		if err != nil {
			logs.WarningLogger.Printf(err.Error())
			return false
		}
	} else {
		// Name has changed. Delete the old record and create it under a new key, since names are used for keys
		_, err = r.redisConn.Do("HSET", "contacts", updatedContact.Name, serializedContact)
		if err != nil {
			logs.WarningLogger.Printf(err.Error())
			return false
		}
		_, err = r.redisConn.Do("HDEL", "contacts", contactName, serializedContact)
		if err != nil {
			logs.WarningLogger.Printf(err.Error())
			return false
		}
	}

	return true
}

func (r *RedisContactRepository) Add(c *models.Contact) bool {
	serializedContact, err := json.Marshal(*c)
	if err != nil {
		logs.WarningLogger.Printf(err.Error())
		return false
	}

	exists, err := redis.Bool(r.redisConn.Do("HEXISTS", "contacts", c.Name))
	if err != nil {
		logs.WarningLogger.Printf(err.Error())
	}
	if exists {
		return false
	}

	_, err = r.redisConn.Do("HSET", "contacts", c.Name, serializedContact)
	if err != nil {
		logs.WarningLogger.Printf(err.Error())
		return false
	}

	return true
}

func (r *RedisContactRepository) Remove(contactName string) bool {
	redisReturn, err := redis.Bool(r.redisConn.Do("HDEL", "contacts", contactName))
	if err != nil {
		logs.WarningLogger.Printf(err.Error())
		return false
	}

	return redisReturn
}

func (r *RedisContactRepository) SearchByName(contactName string) *models.Contact {
	var contact models.Contact
	contactString, err := redis.String(r.redisConn.Do("HGET", "contacts", contactName))
	if err != nil {
		logs.WarningLogger.Printf(err.Error())
		return &models.Contact{}
	}

	json.Unmarshal([]byte(contactString), &contact)

	return &contact
}

func NewRedisContactRepository(redisConn *redis.Conn) *RedisContactRepository {
	return &RedisContactRepository{
		redisConn: *redisConn,
	}
}
