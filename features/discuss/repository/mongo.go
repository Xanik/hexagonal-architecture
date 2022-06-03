package repository

import (
	"fmt"
	"study/config"
	chat "study/features/discuss"
	"study/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoChatRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "chats"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) chat.Repository {
	return &mongoChatRespository{dbs}
}

func (m *mongoChatRespository) Create(a interface{}, collection string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoChatRespository) Find(id string, collection string) (map[string]interface{}, error) {

	var account map[string]interface{}

	coll := m.dbs.DB(dbName).C(collection)

	if collection == "comments" {
		lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
			"from":         "accounts",
			"localField":   "created_by",
			"foreignField": "_id",
			"as":           "created_by",
		}}

		lookUpUsers := bson.M{"$lookup": bson.M{ // lookup the documents table here
			"from":         "accounts",
			"localField":   "users",
			"foreignField": "_id",
			"as":           "users",
		}}

		unWind := bson.M{"$unwind": "$created_by"}

		match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}

		pipe := coll.Pipe([]bson.M{match, lookUpAccount, lookUpUsers, unWind})

		err := pipe.One(&account)

		if err != nil {
			return nil, err
		}

	}

	if collection == "chats" {

		lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
			"from":         "accounts",
			"localField":   "created_by",
			"foreignField": "_id",
			"as":           "created_by",
		}}

		project := bson.M{"$project": bson.M{
			"comment_count": bson.M{"$size": "$comments"},
			"created_by":    1,
			"topic":         1,
			"tags":          1,
			"type":          1,
			"room":          1,
			"created_at":    1,
			"updated_at":    1,
			"body":          1,
			"users":         1,
			"subscribers":   1,
		}}

		projectOut := bson.M{"$project": bson.M{
			"created_by.password": 0,
			"created_by.code":     0,
		}}

		unWind := bson.M{"$unwind": "$created_by"}

		match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}

		pipe := coll.Pipe([]bson.M{match, lookUpAccount, unWind, project, projectOut})

		err := pipe.One(&account)

		if err != nil {
			return nil, err
		}
	}

	return account, nil

}

func (m *mongoChatRespository) FindComments(id string) (interface{}, error) {

	var account interface{}

	coll := m.dbs.DB(dbName).C("comments")
	lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "created_by",
		"foreignField": "_id",
		"as":           "created_by",
	}}

	lookUpUsers := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "users",
		"foreignField": "_id",
		"as":           "users",
	}}

	project := bson.M{"$project": bson.M{
		"created_by.password": 0,
		"created_by.code":     0,
	}}

	unWind := bson.M{"$unwind": "$created_by"}

	match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}

	pipe := coll.Pipe([]bson.M{match, lookUpAccount, lookUpUsers, unWind, project})

	err := pipe.One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil

}

func (m *mongoChatRespository) FindBy(key string, value string, collection string) (map[string]interface{}, error) {

	var account map[string]interface{}

	coll := m.dbs.DB(dbName).C(collection)

	if key == "_id" {
		err := coll.Find(bson.M{"_id": bson.ObjectIdHex(value)}).One(&account)
		if err != nil {
			return nil, err
		}
	} else {
		err := coll.Find(bson.M{key: value}).One(&account)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func (m *mongoChatRespository) FindSingle(key string, value string) (interface{}, error) {

	var account interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	if key == "_id" {
		err := coll.Find(bson.M{"_id": bson.ObjectIdHex(value)}).One(&account)
		if err != nil {
			return nil, err
		}
	} else {
		err := coll.Find(bson.M{key: value}).One(&account)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func (m *mongoChatRespository) FindAll(match bson.M) ([]map[string]interface{}, error) {

	var account []map[string]interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "created_by",
		"foreignField": "_id",
		"as":           "created_by",
	}}

	unWind := bson.M{"$unwind": "$created_by"}

	project := bson.M{"$project": bson.M{
		"comment_count": bson.M{"$size": "$comments"},
		"created_by":    1,
		"topic":         1,
		"tags":          1,
		"type":          1,
		"room":          1,
		"created_at":    1,
		"updated_at":    1,
		"body":          1,
	}}

	projectOut := bson.M{"$project": bson.M{
		"created_by.password": 0,
		"created_by.code":     0,
	}}

	pipe := coll.Pipe([]bson.M{match, lookUpAccount, unWind, project, projectOut})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoChatRespository) Trending() ([]map[string]interface{}, error) {

	var account []map[string]interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "created_by",
		"foreignField": "_id",
		"as":           "created_by",
	}}

	unWind := bson.M{"$unwind": "$created_by"}

	project := bson.M{"$project": bson.M{
		"comment_count": bson.M{"$size": "$comments"},
		"created_by":    1,
		"topic":         1,
		"tags":          1,
		"type":          1,
		"room":          1,
		"created_at":    1,
		"updated_at":    1,
		"body":          1,
	}}

	projectOut := bson.M{"$project": bson.M{
		"created_by.password": 0,
		"created_by.code":     0,
	}}

	match := bson.M{"$match": bson.M{"comments.3": bson.M{"$exists": true}, "type": "public"}}

	pipe := coll.Pipe([]bson.M{match, lookUpAccount, unWind, project, projectOut})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoChatRespository) FindAllAndGroup() ([]interface{}, error) {

	var account []interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	unWind := bson.M{"$unwind": "$tags"}

	group := bson.M{"$group": bson.M{"_id": "$tags", "count": bson.M{"$sum": 1}}}

	pipe := coll.Pipe([]bson.M{unWind, group})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoChatRespository) FindCommentsBy(key string, value string) ([]map[string]interface{}, error) {

	var account []map[string]interface{}

	coll := m.dbs.DB(dbName).C("comments")

	lookUpAccount := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "created_by",
		"foreignField": "_id",
		"as":           "created_by",
	}}

	projectOut := bson.M{"$project": bson.M{
		"created_by.password": 0,
		"created_by.code":     0,
	}}

	unWind := bson.M{"$unwind": "$created_by"}

	match := bson.M{"$match": bson.M{key: bson.ObjectIdHex(value)}}

	pipe := coll.Pipe([]bson.M{match, lookUpAccount, unWind, projectOut})

	err := pipe.All(&account)

	fmt.Println(account, coll)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoChatRespository) Update(a interface{}, id string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoChatRespository) FindWith(query bson.M) (interface{}, error) {
	var account *models.Chat

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(query).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoChatRespository) FindUser(id string) (*models.Account, error) {

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

func (m *mongoChatRespository) Delete(match bson.M, collection string) error {

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Remove(match)

	if err != nil {
		return err
	}

	return nil
}

func (m *mongoChatRespository) UpdateComment(a interface{}, id string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C("comments")

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}
