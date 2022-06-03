package repository

import (
	"fmt"
	"study/config"
	"study/features/network"
	"study/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoNetworkRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "networks"
)

//NewMongoNetworkRepository will create an object that represent the account.Repository interface
func NewMongoNetworkRepository(dbs *mgo.Session) network.Repository {
	return &mongoNetworkRespository{dbs}
}

func (m *mongoNetworkRespository) Create(a interface{}) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoNetworkRespository) Find(id string) (*models.Network, error) {

	var account models.Network

	coll := m.dbs.DB(dbName).C(dbCollection)

	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id)}}

	pipe := coll.Pipe([]bson.M{match})

	err := pipe.One(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (m *mongoNetworkRespository) FindFollowers(id string, skip int, limit int) (interface{}, error) {

	var account []interface{}

	var count map[string]interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id)}}

	lookUpFollowers := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "followers",
		"foreignField": "_id",
		"as":           "followers",
	}}

	project := bson.M{"$project": bson.M{"followers": 1, "_id": 0}}

	size := bson.M{"$project": bson.M{"_id": 0, "count": bson.M{"$size": "$followers"}}}

	unwind := bson.M{"$unwind": "$followers"}

	projectOut := bson.M{"$project": bson.M{"_id": "$followers._id", "first_name": "$followers.first_name", "last_name": "$followers.last_name", "bio": "$followers.bio", "image": "$followers.image", "gender": "$followers.gender"}}

	skips := bson.M{"$skip": skip * limit}

	limits := bson.M{"$limit": limit}

	pipe := coll.Pipe([]bson.M{match, lookUpFollowers, project, unwind, projectOut, skips, limits})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	pipeline := coll.Pipe([]bson.M{match, size})

	errr := pipeline.One(&count)

	if errr != nil {
		return nil, errr
	}

	data := make(map[string]interface{})

	data["followers"] = account

	data["count"] = count["count"]

	return &data, nil
}

func (m *mongoNetworkRespository) FindFollowing(id string, skip int, limit int) (interface{}, error) {

	var account []interface{}

	var count map[string]interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id)}}

	lookUpFollowing := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "accounts",
		"localField":   "following",
		"foreignField": "_id",
		"as":           "following",
	}}

	project := bson.M{"$project": bson.M{"following": 1, "_id": 0}}

	size := bson.M{"$project": bson.M{"_id": 0, "count": bson.M{"$size": "$following"}}}

	unwind := bson.M{"$unwind": "$following"}

	projectOut := bson.M{"$project": bson.M{"_id": "$following._id", "first_name": "$following.first_name", "last_name": "$following.last_name", "bio": "$following.bio", "image": "$following.image", "gender": "$following.gender"}}

	skips := bson.M{"$skip": skip * limit}

	limits := bson.M{"$limit": limit}

	pipe := coll.Pipe([]bson.M{match, lookUpFollowing, project, unwind, projectOut, skips, limits})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	pipeline := coll.Pipe([]bson.M{match, size})

	errr := pipeline.One(&count)

	if errr != nil {
		return nil, errr
	}

	data := make(map[string]interface{})

	data["following"] = account

	data["count"] = count["count"]

	return &data, nil
}

func (m *mongoNetworkRespository) FindBy(key string, value string) (*models.Network, error) {

	var account *models.Network

	coll := m.dbs.DB(dbName).C(dbCollection)

	if key == "_id" {
		err := coll.Find(bson.M{"user_id": bson.ObjectIdHex(value)}).One(&account)
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

func (m *mongoNetworkRespository) FindAll() ([]*models.Network, error) {

	var account []*models.Network

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{}).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoNetworkRespository) Update(a interface{}, id string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Update(bson.M{"user_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoNetworkRespository) GetUsersSuggestedFollowers(id []bson.ObjectId, skip int, limit int) (interface{}, error) {

	var account []models.Follow

	coll := m.dbs.DB(dbName).C("accounts")

	match := bson.M{"$match": bson.M{"interest_id": bson.M{"$in": []bson.ObjectId(id)}}}

	projectOut := bson.M{"$project": bson.M{"_id": 1, "first_name": 1, "last_name": 1, "image": 1, "gender": 1, "bio": 1}}

	skips := bson.M{"$skip": skip * limit}

	limits := bson.M{"$limit": limit}

	pipe := coll.Pipe([]bson.M{match, projectOut, skips, limits})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	count, errr := coll.Find(bson.M{"interest_id": bson.M{"$in": []bson.ObjectId(id)}}).Count()

	if errr != nil {
		return nil, errr
	}

	var data models.Following

	data.Following = account

	data.Count = count

	return &data, nil
}

func (m *mongoNetworkRespository) FindIn(elementMap []string) (interface{}, error) {

	var account []interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{"_id": bson.M{"$in": elementMap}}).All(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (m *mongoNetworkRespository) FindUser(id string, collection string) (*models.Account, error) {

	var account *models.Account

	coll := m.dbs.DB(dbName).C(collection)

	match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}

	pipe := coll.Pipe([]bson.M{match})

	err := pipe.One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoNetworkRespository) FindWith(query bson.M) (interface{}, error) {
	var account *models.Network

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(query).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoNetworkRespository) Notify(a interface{}) (interface{}, error) {

	coll := m.dbs.DB(dbName).C("notifications")

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoNetworkRespository) Count(query bson.M) (interface{}, error) {
	coll := m.dbs.DB(dbName).C(dbCollection)

	account, err := coll.Find(query).Count()

	fmt.Println(account, err, query)

	if err != nil {
		return nil, err
	}

	return account, nil
}
