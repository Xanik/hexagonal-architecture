package repository

import (
	"study/config"
	"study/features/account"
	"study/models"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoAccountRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "accounts"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) account.Repository {
	return &mongoAccountRespository{dbs}
}

func (m *mongoAccountRespository) Create(a *models.Account) (*models.Account, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	coll.EnsureIndex(mgo.Index{Key: []string{"email", "phone", "account_type"}, Unique: true})

	a.ID = bson.NewObjectId()

	a.CreatedAT = time.Now()

	a.UpdatedAt = time.Now()

	a.Step = "CA"

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoAccountRespository) Find(id string) (*models.Account, error) {

	var account models.Account

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.FindId(bson.ObjectIdHex(id)).One(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (m *mongoAccountRespository) FindBy(key string, value string) (*models.Account, error) {

	var account *models.Account

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoAccountRespository) FindAll() ([]*models.Account, error) {

	var account []*models.Account

	coll := m.dbs.DB(dbName).C(dbCollection)

	project := bson.M{"$project": bson.M{
		"password": 0,
	}}

	match := bson.M{"$match": bson.M{"step": "SP"}}

	pipe := coll.Pipe([]bson.M{match, project})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoAccountRespository) Update(a interface{}, id string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoAccountRespository) FindUserNotifications(id string) ([]map[string]interface{}, error) {

	var notifyShares []map[string]interface{}

	var notify []map[string]interface{}

	coll := m.dbs.DB(dbName).C("notifications")

	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id)}}

	lookUpUser := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "user_id",
		"foreignField": "_id",
		"as":           "user_id",
	}}

	lookUpSender := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "sender_id",
		"foreignField": "_id",
		"as":           "sender_id",
	}}

	lookUpContent := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "contents",
		"localField":   "content_id",
		"foreignField": "_id",
		"as":           "content_id",
	}}

	unwindUser := bson.M{"$unwind": "$user_id"}

	unwindSender := bson.M{"$unwind": "$sender_id"}

	unwindContent := bson.M{"$unwind": "$content_id"}

	project := bson.M{"$project": bson.M{"user_id.password": 0, "user_id.code": 0, "sender_id.password": 0, "sender_id.code": 0}}

	pipe := coll.Pipe([]bson.M{match, lookUpUser, lookUpSender, lookUpContent, unwindContent, unwindSender, unwindUser, project})

	err := pipe.All(&notifyShares)

	if err != nil {
		return nil, err
	}

	query := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id), "content_id": nil}}

	pipeline := coll.Pipe([]bson.M{query, lookUpUser, lookUpSender, unwindSender, unwindUser, project})
	e := pipeline.All(&notify)
	if e != nil {
		return nil, e
	}

	data := append(notify, notifyShares...)
	return data, nil
}

func (m *mongoAccountRespository) GetIndexedAccount(id string) (*models.SearchAccount, error) {
	var data *models.SearchAccount

	coll := m.dbs.DB(dbName).C(dbCollection)

	project := bson.M{"$project": bson.M{
		"id":         "$_id",
		"_id":        0,
		"first_name": 1,
		"last_name":  1,
		"email":      1,
		"phone":      1,
		"image":      1,
		"bio":        1,
	}}

	match := bson.M{"$match": bson.M{"step": "SP", "_id": bson.ObjectIdHex(id)}}

	pipe := coll.Pipe([]bson.M{match, project})

	err := pipe.One(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
