package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Interest model
type Interest struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Image     string        `json:"image" bson:"image"`
	CreatedAT time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
