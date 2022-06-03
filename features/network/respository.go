package network

import (
	"study/models"

	"gopkg.in/mgo.v2/bson"
)

//Repository interface represents network services
type Repository interface {
	Create(interface{}) (interface{}, error)
	Find(string) (*models.Network, error)
	FindFollowers(string, int, int) (interface{}, error)
	FindFollowing(string, int, int) (interface{}, error)
	FindBy(key string, value string) (*models.Network, error)
	FindIn([]string) (interface{}, error)
	FindAll() ([]*models.Network, error)
	FindUser(string, string) (*models.Account, error)
	Update(user interface{}, id string) (interface{}, error)
	GetUsersSuggestedFollowers([]bson.ObjectId, int, int) (interface{}, error)
	FindWith(query bson.M) (interface{}, error)
	Notify(interface{}) (interface{}, error)
	Count(match bson.M) (interface{}, error)
}
