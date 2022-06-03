package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Upload model
type Upload struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Url       bson.ObjectId `json:"url" bson:"url"`
	CreatedAT time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
