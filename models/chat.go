package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Chat model
type Chat struct {
	ID          bson.ObjectId   `json:"_id" bson:"_id"`
	CreatedBy   bson.ObjectId   `json:"created_by" bson:"created_by"`
	Topic       string          `json:"topic" bson:"topic"`
	Body        string          `json:"body" bson:"body"`
	Comments    []bson.ObjectId `json:"comments" bson:"comments"`
	CourseID    bson.ObjectId   `json:"course_id,omitempty" bson:"course_id,omitempty"`
	Users       []bson.ObjectId `json:"users,omitempty" bson:"users,omitempty"`
	Subscribers []bson.ObjectId `json:"subscribers,omitempty" bson:"subscribers,omitempty"`
	Tags        []string        `json:"tags" bson:"tags"`
	Room        string          `json:"room" bson:"room"`
	Type        string          `json:"type" bson:"type"`
	CreatedAT   time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" bson:"updated_at"`
}

//Comment model
type Comment struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	CreatedBy    bson.ObjectId `json:"created_by" bson:"created_by"`
	DiscussionID bson.ObjectId `json:"discussion_id" bson:"discussion_id"`
	Message      string        `json:"message" bson:"message"`
	Attachment   string        `json:"attachment" bson:"attachment"`
	CreatedAT    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
}

const (
	Private string = "private"
	Public  string = "public"
)

func (c Chat) ValidChat() bool {
	// If it's a user input, you'd want to validate Type's underlying
	// string isn't out of the enum's range.
	if c.Type == Private || c.Type == Public {
		return true
	}
	return false
}
