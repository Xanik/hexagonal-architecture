package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Organization model
type Organization struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Address   string        `json:"address" bson:"address"`
	Email     string        `json:"email" bson:"email"`
	Image     string        `json:"image" bson:"image"`
	Theme1    string        `json:"theme_color_1" bson:"theme_color_1"`
	Theme2    string        `json:"theme_color_2" bson:"theme_color_2"`
	Folder1   string        `json:"folder_colour_1" bson:"folder_colour_1"`
	Folder2   string        `json:"folder_colour_2" bson:"folder_colour_2"`
	CreatedAT time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
