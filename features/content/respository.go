package content

import (
	models "study/models"

	"gopkg.in/mgo.v2/bson"
)

//Repository interface represents Content services
type Repository interface {
	Create(*models.Content) (*models.Content, error)
	Find(id string, collection string) (interface{}, error)
	FindBy(key string, value string) (*models.Content, error)
	FindAll(int, int) (*models.ArrayResponse, error)
	FindWith(query bson.M, collection string) (interface{}, error)
	FindAllWithQuery(bson.M, int, int) (*models.ArrayResponse, error)
	FindInInterest([]bson.ObjectId, string) (interface{}, error)
	FindInContent([]string, string, int, int) (*models.ArrayResponse, error)
	Update(*models.Content, string) (*models.Content, error)
	CountContent(match bson.M) (interface{}, error)
}
