package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Network model
type Network struct {
	ID        bson.ObjectId   `json:"_id" bson:"_id"`
	UserID    bson.ObjectId   `json:"user_id" bson:"user_id"`
	Followers []bson.ObjectId `json:"followers,omitempty" bson:"followers,omitempty"`
	Following []bson.ObjectId `json:"following,omitempty" bson:"following,omitempty"`
	Follow    []Follow        `json:"follow,omitempty" bson:"follow,omitempty"`
	CreatedAT time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
}

type Following struct {
	Following []Follow    `json:"following" bson:"following"`
	Count     interface{} `json:"count" bson:"count"`
}

type Followers struct {
	Followers []Follow    `json:"followers" bson:"followers"`
	Count     interface{} `json:"count" bson:"count"`
}

type Follow struct {
	ID         bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	FirstName  string        `json:"first_name" bson:"first_name"`
	LastName   string        `json:"last_name" bson:"last_name"`
	Bio        string        `json:"bio,omitempty" bson:"bio,omitempty"`
	Image      string        `json:"image" bson:"image"`
	Gender     string        `json:"gender,omitempty" bson:"gender,omitempty"`
	IsFollowed bool          `json:"isfollowed" bson:"isfollowed"`
}

type UserNetwork struct {
	ID        bson.ObjectId   `json:"_id" bson:"_id"`
	UserID    bson.ObjectId   `json:"user_id" bson:"user_id"`
	Following []bson.ObjectId `json:"following" bson:"following"`
	Followers []bson.ObjectId `json:"followers" bson:"followers"`
	CreatedAT time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
}
