package interaction

import (
	"study/models"

	"gopkg.in/mgo.v2/bson"
)

//Repository interface represents account services
type Repository interface {
	Create(interface{}) (interface{}, error)
	Find(id string, collection string) (interface{}, error)
	FindBy(key string, value string) (interface{}, error)
	FindByCollection(collection string, skip int, limit int, id string) (interface{}, error)
	FindWith(query bson.M) (interface{}, error)
	FindByAndGroup(bson.M, int, int) ([]interface{}, error)
	FindByAndGroupCount(bson.M, int, int) ([]interface{}, error)
	FindAll(bson.M, int, int) (interface{}, error)
	Update(interface{}, string, string) (interface{}, error)
	Delete(string, string, string) error
	DeleteCollection(bson.M) error
	DeleteMany(bson.M) error
	UpdateCollectionName(interface{}, bson.M) (interface{}, error)
	CountContent(match bson.M) (int, error)
	FindUser(string) (*models.Account, error)
	Notify(interface{}) (interface{}, error)
	GetIndexedContent(id string) (*models.SearchContent, error)
	IndexDocument(interface{}, string, string)
	CreateIndex(string)
	DeleteIndex(string, string)
}
