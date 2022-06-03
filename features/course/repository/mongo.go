package repository

import (
	"study/config"
	"study/features/course"
	"study/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoCourseRespository struct {
	dbs *mgo.Session
}

var (
	dbName       = config.Env.DBNAME
	dbCollection = "courses"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoAccountRepository(dbs *mgo.Session) course.Repository {
	return &mongoCourseRespository{dbs}
}

func (m *mongoCourseRespository) Create(a interface{}) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoCourseRespository) Find(id string) (interface{}, error) {
	var account interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)

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
		"resources":   1,
	}}

	o2 := bson.M{"$group": bson.M{"_id": "$_id", "title": bson.M{"$first": "$title"}, "resources": bson.M{"$first": "$resources"}, "summary": bson.M{"$first": "$summary"}, "description": bson.M{"$first": "$description"}, "photo_url": bson.M{"$first": "$photo_url"}, "video_url": bson.M{"$first": "$video_url"}, "duration": bson.M{"$first": "$duration"}, "rating": bson.M{"$first": "$rating"}, "provider": bson.M{"$first": "$provider"}, "created_at": bson.M{"$first": "$created_at"}, "updated_at": bson.M{"$first": "$updated_at"}, "modules": bson.M{"$addToSet": bson.M{"title": "$titles", "lessons": "$lessons"}}}}

	match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}

	pipe := coll.Pipe([]bson.M{match, unWind, lookUpLessons, project, o2})

	err := pipe.One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoCourseRespository) FindBy(key string, value string) (*models.Course, error) {
	var account *models.Course

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Find(bson.M{key: value}).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoCourseRespository) FindAll() ([]interface{}, error) {
	var account []interface{}

	coll := m.dbs.DB(dbName).C(dbCollection)
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
		"resources":   1,
	}}

	o2 := bson.M{"$group": bson.M{"_id": "$_id", "title": bson.M{"$first": "$title"}, "resources": bson.M{"$first": "$resources"}, "summary": bson.M{"$first": "$summary"}, "description": bson.M{"$first": "$description"}, "photo_url": bson.M{"$first": "$photo_url"}, "video_url": bson.M{"$first": "$video_url"}, "duration": bson.M{"$first": "$duration"}, "rating": bson.M{"$first": "$rating"}, "provider": bson.M{"$first": "$provider"}, "created_at": bson.M{"$first": "$created_at"}, "updated_at": bson.M{"$first": "$updated_at"}, "modules": bson.M{"$addToSet": bson.M{"title": "$titles", "lessons": "$lessons"}}}}

	pipe := coll.Pipe([]bson.M{unWind, lookUpLessons, project, o2})

	err := pipe.All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoCourseRespository) Update(a interface{}, id string) (interface{}, error) {

	coll := m.dbs.DB(dbName).C(dbCollection)

	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": a})

	if err != nil {
		return nil, err
	}

	return a, nil
}
