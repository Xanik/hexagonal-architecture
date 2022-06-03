package repository

import (
	"study/config"
	"study/features/organization"
	models "study/models"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoOrganizationRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "organizations"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) organization.Repository {
	return &mongoOrganizationRespository{dbs}
}

func (m *mongoOrganizationRespository) Create(a *models.Organization) (*models.Organization, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	a.CreatedAT = time.Now()

	a.UpdatedAt = time.Now()

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoOrganizationRespository) Find(id string) (*models.Organization, error) {
	var account *models.Organization

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoOrganizationRespository) FindBy(key string, value string) (*models.Organization, error) {
	var account *models.Organization

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoOrganizationRespository) FindAll() ([]*models.Organization, error) {
	var account []*models.Organization

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{}).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoOrganizationRespository) Update(a *models.Organization, id string) (*models.Organization, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	a.UpdatedAt = time.Now()

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}
