package chat

import (
	"study/models"

	"gopkg.in/mgo.v2/bson"
)

//Repository interface represents chat services
type Repository interface {
	Create(interface{}, string) (interface{}, error)
	Find(id string, collection string) (map[string]interface{}, error)
	FindBy(key string, value string, collection string) (map[string]interface{}, error)
	FindAll(match bson.M) ([]map[string]interface{}, error)
	FindCommentsBy(key string, value string) ([]map[string]interface{}, error)
	Update(user interface{}, id string) (interface{}, error)
	FindAllAndGroup() ([]interface{}, error)
	Trending() ([]map[string]interface{}, error)
	FindWith(query bson.M) (interface{}, error)
	FindComments(id string) (interface{}, error)
	FindSingle(key string, value string) (interface{}, error)
	FindUser(string) (*models.Account, error)
	Delete(match bson.M, collection string) error
	UpdateComment(m interface{}, id string) (interface{}, error)
}
