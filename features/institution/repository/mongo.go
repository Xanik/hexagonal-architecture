package repository

import (
	"study/config"
	"study/features/institution"
	"study/models"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoInstitutionRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "institutions"
)

//NewMongoInstitutionRepository will create an object that represent the account.Repository interface
func NewMongoInstitutionRepository(dbs *mgo.Session) institution.Repository {
	return &mongoInstitutionRespository{dbs}
}

func (m *mongoInstitutionRespository) Create(a *models.Institution) (*models.Institution, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	a.ID = bson.NewObjectId()

	a.CreatedAT = time.Now()

	a.UpdatedAt = time.Now()

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoInstitutionRespository) Find(id string) (*models.Institution, error) {

	var account models.Institution

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.ObjectIdHex(id)).One(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (m *mongoInstitutionRespository) FindBy(key string, value string) (*models.Institution, error) {

	var account *models.Institution

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInstitutionRespository) FindAll(skip int, limit int) ([]*models.Institution, error) {

	var account []*models.Institution

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{}).Skip(skip).Limit(limit).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInstitutionRespository) Update(a interface{}, id string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}
