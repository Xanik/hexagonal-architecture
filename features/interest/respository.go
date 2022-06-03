package interest

import (
	models "study/models"

	"gopkg.in/mgo.v2/bson"
)

//Repository interface represents account services
type Repository interface {
	Create(*models.Interest) (*models.Interest, error)
	Find(id string, collection string) (interface{}, error)
	FindBy(key string, value string) (*models.Interest, error)
	FindAll() ([]*models.Interest, error)
	FindIn([]bson.ObjectId, string) (interface{}, error)
	Update(*models.Interest, string) (*models.Interest, error)
}
