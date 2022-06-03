package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//Content model
type SearchContent struct {
	ID          bson.ObjectId `json:"id" bson:"id"`
	Title       string        `json:"title" bson:"title"`
	Summary     string        `json:"summary" bson:"summary"`
	Description string        `json:"description" bson:"description"`
	Image       string        `json:"image" bson:"image"`
	Library     string        `json:"library" bson:"library"`
	Likes       int           `json:"likes" bson:"likes"`
	Saves       int           `json:"saves" bson:"saves"`
}

//Account model
type SearchAccount struct {
	ID        bson.ObjectId `json:"id" bson:"id"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	Email     string        `json:"email,omitempty" bson:"email,omitempty"`
	Phone     string        `json:"phone,omitempty" bson:"phone,omitempty"`
	Bio       string        `json:"bio" bson:"bio"`
	Image     string        `json:"image" bson:"image"`
}

//Course model
type SearchCourse struct {
	ID          bson.ObjectId `json:"id" bson:"id"`
	Title       string        `json:"title" bson:"title"`
	Summary     string        `json:"summary" bson:"summary"`
	Description string        `json:"description" bson:"description"`
	PhotoURL    string        `json:"photo_url" bson:"photo_url"`
	VideoURL    string        `json:"video_url" bson:"video_url"`
	Duration    int           `json:"duration" bson:"duration"`
	Rating      float32       `json:"rating" bson:"rating"`
	Provider    string        `json:"provider,omitempty" bson:"provider,omitempty"`
	CreatedAT   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}
