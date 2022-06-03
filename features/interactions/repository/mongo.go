package repository

import (
	"fmt"
	"study/config"
	interaction "study/features/interactions"
	"study/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoInteractionRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "interactions"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) interaction.Repository {
	return &mongoInteractionRespository{dbs}
}

func (m *mongoInteractionRespository) Create(a interface{}) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoInteractionRespository) Find(id string, collection string) (interface{}, error) {
	var account interface{}

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInteractionRespository) FindBy(key string, value string) (interface{}, error) {
	var account *models.Interactions

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInteractionRespository) FindByCollection(collection string, skip int, limit int, id string) (interface{}, error) {

	match := bson.M{"$match": bson.M{"collection_id": bson.ObjectIdHex(collection), "user_id": bson.ObjectIdHex(id), "type": "collection"}}

	var account interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "user_id",
		"foreignField": "_id",
		"as":           "user_id",
	}}

	lookUpContent := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "contents",
		"localField":   "content_id",
		"foreignField": "_id",
		"as":           "content_id",
	}}

	unWind := bson.M{"$unwind": "$user_id"}

	unWindContent := bson.M{"$unwind": "$content_id"}

	skips := bson.M{"$skip": skip * limit}

	limits := bson.M{"$limit": limit}

	projectOut := bson.M{"$project": bson.M{"name": 1, "content_id": 1, "color": 1}}

	o2 := bson.M{"$group": bson.M{"_id": "$name", "content": bson.M{"$push": "$content_id"}, "color": bson.M{"$push": "$color"}}}

	project := bson.M{"$project": bson.M{"name": "$_id", "content": 1, "_id": 0, "color": bson.M{"$arrayElemAt": []interface{}{"$color", 0}}}}

	pipe := coll.Pipe([]bson.M{match, lookUpAccount, lookUpContent, unWind, unWindContent, skips, limits, projectOut, o2, project})

	err := pipe.One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInteractionRespository) FindWith(query bson.M) (interface{}, error) {
	var account *models.Interactions

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(query).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInteractionRespository) FindAll(match bson.M, skip int, limit int) (interface{}, error) {
	var account *models.Data

	coll := m.dbs.DB(dbName).C(dbCollection)

	lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "user_id",
		"foreignField": "_id",
		"as":           "user_id",
	}}

	lookUpContent := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "contents",
		"localField":   "content_id",
		"foreignField": "_id",
		"as":           "content",
	}}

	unWind := bson.M{"$unwind": "$user_id"}
	unWindContent := bson.M{"$unwind": "$content"}

	project := bson.M{"$project": bson.M{"content_id": 0}}

	o2 := bson.M{"$group": bson.M{"_id": "$name", "content": bson.M{"$push": "$content"}, "color": bson.M{"$push": "$color"}}}

	projectOut := bson.M{"$project": bson.M{"_id": 0, "content": 1, "color": bson.M{"$arrayElemAt": []interface{}{"$color", 0}}}}

	skips := bson.M{"$skip": skip * limit}

	limits := bson.M{"$limit": limit}

	pipe := coll.Pipe([]bson.M{match, lookUpAccount, lookUpContent, unWind, unWindContent, skips, limits, project, o2, projectOut})

	err := pipe.One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInteractionRespository) FindByAndGroup(match bson.M, skip int, limit int) ([]interface{}, error) {
	var account []interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "user_id",
		"foreignField": "_id",
		"as":           "user_id",
	}}

	lookUpContent := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "contents",
		"localField":   "content_id",
		"foreignField": "_id",
		"as":           "content_id",
	}}

	unWind := bson.M{"$unwind": "$user_id"}

	unWindContent := bson.M{"$unwind": "$content_id"}

	skips := bson.M{"$skip": skip * limit}

	limits := bson.M{"$limit": limit}

	projectOut := bson.M{"$project": bson.M{"name": 1, "content_id": 1, "color": 1}}

	o2 := bson.M{"$group": bson.M{"_id": "$name", "contents": bson.M{"$push": "$content_id"}, "color": bson.M{"$push": "$color"}}}

	project := bson.M{"$project": bson.M{"name": "$_id", "contents": 1, "_id": 0, "color": bson.M{"$arrayElemAt": []interface{}{"$color", 0}}}}

	pipe := coll.Pipe([]bson.M{match, lookUpAccount, lookUpContent, unWind, unWindContent, skips, limits, projectOut, o2, project})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (m *mongoInteractionRespository) FindByAndGroupCount(match bson.M, skip int, limit int) ([]interface{}, error) {
	var account []interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	//skips := bson.M{"$skip": skip * limit}
	//
	//limits := bson.M{"$limit": limit}

	o2 := bson.M{"$group": bson.M{"_id": bson.M{"name": "$name", "collection_id": "$collection_id"}, "count": bson.M{"$push": "$content_id"}, "color": bson.M{"$push": "$color"}}}

	project := bson.M{"$project": bson.M{"name": "$_id.name", "collection_id": "$_id.collection_id", "_id": 0, "count": bson.M{"$size": "$count"}, "color": bson.M{"$arrayElemAt": []interface{}{"$color", 0}}}}

	pipe := coll.Pipe([]bson.M{match, o2, project})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (m *mongoInteractionRespository) Update(a interface{}, id string, collection string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoInteractionRespository) Delete(account string, content string, types string) error {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Remove(bson.M{"user_id": bson.ObjectIdHex(account), "content_id": bson.ObjectIdHex(content), "type": types})

	if err != nil {
		return err
	}

	return nil
}

func (m *mongoInteractionRespository) DeleteCollection(match bson.M) error {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Remove(match)

	if err != nil {
		return err
	}

	return nil
}

func (m *mongoInteractionRespository) DeleteMany(match bson.M) error {

	coll := m.dbs.DB(dbName).C(dbCollection)

	data, err := coll.RemoveAll(match)

	if err != nil {
		return err
	}

	fmt.Println("data", data)

	return nil
}

func (m *mongoInteractionRespository) UpdateCollectionName(a interface{}, query bson.M) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	_, err := coll.UpdateAll(query, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoInteractionRespository) CountContent(query bson.M) (int, error) {
	coll := m.dbs.DB(dbName).C(dbCollection)

	account, err := coll.Find(query).Count()

	fmt.Println(account, err, query)

	if err != nil {
		return 0, err
	}

	return account, nil
}

func (m *mongoInteractionRespository) FindUser(id string) (*models.Account, error) {

	var account *models.Account

	coll := m.dbs.DB(dbName).C("accounts")

	match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}

	pipe := coll.Pipe([]bson.M{match})

	err := pipe.One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoInteractionRespository) Notify(a interface{}) (interface{}, error) {

	coll := m.dbs.DB(dbName).C("notifications")

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoInteractionRespository) GetIndexedContent(id string) (*models.SearchContent, error) {
	var data *models.SearchContent

	coll := m.dbs.DB(dbName).C("contents")

	project := bson.M{"$project": bson.M{
		"id":          "$_id",
		"_id":         0,
		"title":       1,
		"image":       1,
		"description": 1,
		"summary":     1,
		"library":     1,
		"likes":       1,
		"saves":       1,
	}}

	match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}

	pipe := coll.Pipe([]bson.M{match, project})

	err := pipe.One(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
