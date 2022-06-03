package repository

import (
	"study/config"
	"study/features/interest"
	"study/models"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoInterestRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "interests"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) interest.Repository {
	return &mongoInterestRespository{dbs}
}

func (m *mongoInterestRespository) Create(a *models.Interest) (*models.Interest, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	a.CreatedAT = time.Now()

	a.UpdatedAt = time.Now()

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoInterestRespository) Find(id string, collection string) (interface{}, error) {

	var account interface{}

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInterestRespository) FindBy(key string, value string) (*models.Interest, error) {

	var account *models.Interest

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInterestRespository) FindAll() ([]*models.Interest, error) {

	var account []*models.Interest

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{}).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInterestRespository) Update(a *models.Interest, id string) (*models.Interest, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	a.UpdatedAt = time.Now()

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoInterestRespository) FindIn(elementMap []bson.ObjectId, collection string) (interface{}, error) {

	var account []interface{}

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(bson.M{"_id": bson.M{"$ne": elementMap}}).All(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
