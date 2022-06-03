package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Developer Struct
type Developer struct {
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	AppName      string        `json:"app_name,omitempty" bson:"app_name,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"  bson:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"  bson:"refresh_token,omitempty"`
	ExpiresIn    uint          `json:"expires_in,omitempty"  bson:"expires_in,omitempty"`
	CreatedAT    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
}

// DeveloperRequest Struct
type DeveloperRequest struct {
	Interests   []string `json:"interests,omitempty"`
}

// DeveloperResponse Struct
type DeveloperResponse struct {
	Interests   []string `json:"interests,omitempty"`
}