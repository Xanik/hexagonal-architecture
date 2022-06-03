package repository

import (
	"study/config"
	"study/features/upload"
	"study/models"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoUploadRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "uploads"
)

//NewMongoUploadRepository will create an object that represent the upload.Repository interface
func NewMongoUploadRepository(dbs *mgo.Session) upload.Repository {
	return &mongoUploadRespository{dbs}
}

var sessionTTL = mgo.Index{
	Key:         []string{"_id"},
	Unique:      false,
	DropDups:    false,
	Background:  true,
	ExpireAfter: 139 * 139 * 139} // one hour

func (m *mongoUploadRespository) Create(a *models.Upload) (*models.Upload, error) {

	if err := m.dbs.DB(dbName).C(dbCollection).EnsureIndex(sessionTTL); err != nil {
		panic(err)
	}

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

func (m *mongoUploadRespository) Find(id string) (*models.Upload, error) {

	var account models.Upload

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (m *mongoUploadRespository) FindBy(key string, value string) (*models.Upload, error) {

	var account *models.Upload

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoUploadRespository) FindAll() ([]*models.Upload, error) {

	var account []*models.Upload

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{}).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoUploadRespository) Update(a interface{}, id string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}
