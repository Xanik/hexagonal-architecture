package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Institution model
type Institution struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Type      string        `json:"type" bson:"type"`
	CreatedAT time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
