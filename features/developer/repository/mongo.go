package repository

import (
	"study/config"
	"study/features/developer"
	"study/models"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoRespository struct {
	dbs *mgo.Session
}

var (
	dbName                = config.Env.DBNAME
	dbContentCollection   = "contents"
	dbDeveloperCollection = "developers"
	dbCourseCollection    = "courses"
)

//NewMongoAccountRepository will create an object that represent the account.Repository interface
func NewMongoRepository(dbs *mgo.Session) developer.Repository {
	return &mongoRespository{dbs}
}

func (m *mongoRespository) RequestAccess(a *models.Developer) (*models.Developer, error) {

	coll := m.dbs.DB(dbName).C(dbDeveloperCollection)

	a.CreatedAT = time.Now()

	a.UpdatedAt = time.Now()

	err := coll.Insert(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *mongoRespository) GetDeveloperAccount(refreshToken string) (*models.Developer, error) {

	var account *models.Developer

	coll := m.dbs.DB(dbName).C(dbDeveloperCollection)

	err := coll.Find(bson.M{"refresh_token": refreshToken}).One(&account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoRespository) GetAllContents() ([]*models.Content, error) {

	var account []*models.Content

	coll := m.dbs.DB(dbName).C(dbContentCollection)

	err := coll.Find(bson.M{}).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoRespository) GetContentByInterest(elementMap []string) ([]*models.Content, error) {

	var account []*models.Content

	coll := m.dbs.DB(dbName).C(dbContentCollection)

	err := coll.Find(bson.M{"tags": bson.M{"$in": elementMap}}).All(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongoRespository) GetAllCourses() ([]*models.Course, error) {
	var account []*models.Course

	coll := m.dbs.DB(dbName).C(dbCourseCollection)
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
