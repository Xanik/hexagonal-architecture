package repository

import (
	"study/config"
	"study/features/search"
	"study/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoSearchRespository struct {
	dbs *mgo.Session
}

var (
	dbName = config.Env.DBNAME
)

//NewMongoAccountRepository will create an object that represent the search.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) search.Repository {
	return &mongoSearchRespository{dbs}
}

func (m *mongoSearchRespository) FindAllAccounts(collection string) ([]map[string]interface{}, error) {
	var account []map[string]interface{}

	coll := m.dbs.DB(dbName).C(collection)

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

	match := bson.M{"$match": bson.M{"step": "SP"}}

	pipe := coll.Pipe([]bson.M{match, project})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoSearchRespository) FindAllContents(collection string) ([]*models.SearchContent, error) {
	var account []*models.SearchContent

	coll := m.dbs.DB(dbName).C(collection)

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

	pipe := coll.Pipe([]bson.M{project})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoSearchRespository) FindAllInterests(collection string) ([]interface{}, error) {
	var account []interface{}

	coll := m.dbs.DB(dbName).C(collection)

	project := bson.M{"$project": bson.M{
		"id":    "$_id",
		"_id":   0,
		"name":  1,
		"image": 1,
	}}

	pipe := coll.Pipe([]bson.M{project})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoSearchRespository) FindAllCourses(collection string) ([]map[string]interface{}, error) {
	var account []map[string]interface{}

	coll := m.dbs.DB(dbName).C(collection)
	lookUpLessons := bson.M{"$lookup": bson.M{ // lookup the documents table here
		"from":         "lessons",
		"localField":   "modules.lessons",
		"foreignField": "_id",
		"as":           "module_lessons",
	}}

	unWind := bson.M{"$unwind": "$modules"}

	project := bson.M{"$project": bson.M{
		"titles":      "$modules.title",
		"lessons":     "$module_lessons",
		"title":       1,
		"summary":     1,
		"description": 1,
		"photo_url":   1,
		"video_url":   1,
		"duration":    1,
		"rating":      1,
		"provider":    1,
		"created_at":  1,
		"updated_at":  1,
	}}

	o2 := bson.M{"$group": bson.M{"_id": "$_id", "title": bson.M{"$first": "$title"}, "summary": bson.M{"$first": "$summary"}, "description": bson.M{"$first": "$description"}, "photo_url": bson.M{"$first": "$photo_url"}, "video_url": bson.M{"$first": "$video_url"}, "duration": bson.M{"$first": "$duration"}, "rating": bson.M{"$first": "$rating"}, "provider": bson.M{"$first": "$provider"}, "created_at": bson.M{"$first": "$created_at"}, "updated_at": bson.M{"$first": "$updated_at"}, "modules": bson.M{"$addToSet": bson.M{"title": "$titles", "lessons": "$lessons"}}}}

	pipe := coll.Pipe([]bson.M{unWind, lookUpLessons, project, o2})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoSearchRespository) FindWith(query bson.M, collection string) (interface{}, error) {
	var account interface{}

	coll := m.dbs.DB(dbName).C(collection)

	err := coll.Find(query).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}
