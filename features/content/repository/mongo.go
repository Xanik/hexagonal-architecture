package repository

import (
	"fmt"
	"study/config"
	"study/features/content"
	"study/models"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoContentRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "contents"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) content.Repository {
	return &mongoContentRespository{dbs}
}

func (m *mongoContentRespository) Create(a *models.Content) (*models.Content, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoContentRespository) Find(id string, collection string) (interface{}, error) {

	var account interface{}

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoContentRespository) FindBy(key string, value string) (*models.Content, error) {

	var account *models.Content

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoContentRespository) FindAll(skip int, limit int) (*models.ArrayResponse, error) {

	account := make([]*models.Content, 0)

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{}).Skip(skip * limit).Limit(limit).All(&account)

	if err != nil {
		return nil, err
	}

	var data *models.ArrayResponse

	data.Contents = account

	data.Count = len(account)

	return data, nil
}

func (m *mongoContentRespository) Update(a *models.Content, id string) (*models.Content, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	a.UpdatedAt = time.Now()

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoContentRespository) FindInInterest(elementMap []bson.ObjectId, collection string) (interface{}, error) {

	var account []models.Interest

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(bson.M{"_id": bson.M{"$in": elementMap}}).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoContentRespository) FindInContent(elementMap []string, collection string, skip int, limit int) (*models.ArrayResponse, error) {

	account := make([]*models.Content, 0)

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(bson.M{"tags": bson.M{"$in": elementMap}}).Skip(skip * limit).Limit(limit).All(&account)

	if err != nil {
		return nil, err
	}

	var data models.ArrayResponse

	res, e := m.CountContent(bson.M{"tags": bson.M{"$in": elementMap}})

	if e != nil {
		return nil, e
	}

	data.Contents = account

	data.Count = res.(int)

	return &data, nil
}

func (m *mongoContentRespository) FindAllWithQuery(match bson.M, skip int, limit int) (*models.ArrayResponse, error) {
	account := make([]*models.Content, 0)

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(match).Skip(skip * limit).Limit(limit).All(&account)

	if err != nil {
		return nil, err
	}

	var data models.ArrayResponse

	res, e := m.CountContent(match)

	if e != nil {
		return nil, e
	}

	data.Contents = account

	data.Count = res.(int)

	return &data, nil
}

func (m *mongoContentRespository) FindWith(query bson.M, collection string) (interface{}, error) {
	var account *models.Interactions

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(query).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoContentRespository) CountContent(query bson.M) (interface{}, error) {
	coll := m.dbs.DB(dbName).C(dbCollection)

	account, err := coll.Find(query).Count()

	fmt.Println(account, err, query)

	if err != nil {
		return nil, err
	}

	return account, nil
}
