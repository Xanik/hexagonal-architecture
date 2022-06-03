package search

import (
	"study/models"

	"gopkg.in/mgo.v2/bson"
)

//Repository interface represents search services
type Repository interface {
	FindAllAccounts(collection string) ([]map[string]interface{}, error)
	FindAllContents(collection string) ([]*models.SearchContent, error)
	FindAllInterests(collection string) ([]interface{}, error)
	FindAllCourses(collection string) ([]map[string]interface{}, error)
	Elastic(id string) (interface{}, error)
	IndexDocument(interface{}, string, string) (string, error)
	SearchElastic(string, string, string) ([]map[string]interface{}, error)
	CreateIndex(string) (string, error)
	FindWith(query bson.M, collection string) (interface{}, error)
}
