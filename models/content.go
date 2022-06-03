package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Content model
type Content struct {
	ID          bson.ObjectId   `json:"_id" bson:"_id"`
	Title       string          `json:"title" bson:"title"`
	Summary     string          `json:"summary" bson:"summary"`
	Description string          `json:"description" bson:"description"`
	Image       string          `json:"image" bson:"image"`
	Platform    string          `json:"platform" bson:"platform"`
	Library     string          `json:"library" bson:"library"`
	Comments    []bson.ObjectId `json:"comments,omitempty" bson:"comments,omitempty"`
	Likes       int             `json:"likes" bson:"likes"`
	Keywords    []string        `json:"keywords,omitempty" bson:"keywords,omitempty"`
	Saves       int             `json:"saves" bson:"saves"`
	Shares      int             `json:"shares" bson:"shares"`
	IsLiked     bool            `json:"isliked" bson:"isliked"`
	IsSaved     bool            `json:"issaved" bson:"issaved"`
	Tags        []string        `json:"tags" bson:"tags"`
	Authors     []string        `json:"authors,omitempty" bson:"authors,omitempty"`
	// Media       []Media         `json:"media,omitempty" bson:"media,omitempty"`
	CreatedAT time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

//Media model
type Media struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Type      string        `json:"type" bson:"type"`
	Location  string        `json:"location" bson:"location"`
	CreatedAT time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

//Interactions model
type Interactions struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	Type         string        `json:"type" bson:"type"`
	Name         string        `json:"name,omitempty" bson:"name,omitempty"`
	Comment      string        `json:"comment,omitempty" bson:"comment,omitempty"`
	ContentID    bson.ObjectId `json:"content_id,omitempty" bson:"content_id,omitempty"`
	CollectionID bson.ObjectId `json:"collection_id,omitempty" bson:"collection_id,omitempty"`
	Count        interface{}   `json:"count,omitempty" bson:"count,omitempty"`
	Color        string        `json:"color,omitempty" bson:"color,omitempty"`
	UserID       bson.ObjectId `json:"user_id" bson:"user_id"`
	Recipient    bson.ObjectId `json:"recipient,omitempty" bson:"recipient,omitempty"`
	CreatedAT    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
}

//Collection model
type Collection struct {
	Name         string        `json:"name,omitempty" bson:"name,omitempty"`
	Content      []Content     `json:"content,omitempty" bson:"content,omitempty"`
	Color        string        `json:"color,omitempty" bson:"color,omitempty"`
	CollectionID bson.ObjectId `json:"collection_id,omitempty" bson:"collection_id,omitempty"`
}

//Data Model
type Data struct {
	Content []Content `json:"content"`
}

//Course model
type Course struct {
	ID          bson.ObjectId   `json:"_id" bson:"_id"`
	Title       string          `json:"title" bson:"title"`
	Summary     string          `json:"summary" bson:"summary"`
	Description string          `json:"description" bson:"description"`
	PhotoURL    string          `json:"photo_url" bson:"photo_url"`
	VideoURL    string          `json:"video_url" bson:"video_url"`
	Duration    int             `json:"duration" bson:"duration"`
	Rating      float32         `json:"rating" bson:"rating"`
	Comments    []bson.ObjectId `json:"comments,omitempty" bson:"comments,omitempty"`
	Provider    string          `json:"provider,omitempty" bson:"provider,omitempty"`
	Modules     []Module        `json:"modules,omitempty" bson:"modules,omitempty"`
	Resources    []Resource      `json:"resources,omitempty" bson:"resources,omitempty"`
	CreatedAT   time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" bson:"updated_at"`
}

//Module model
type Module struct {
	Title     string    `json:"title" bson:"title"`
	Lessons   []Lesson  `json:"lessons,omitempty" bson:"lessons,omitempty"`
	CreatedAT time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

//Lesson model
type Lesson struct {
	Title         string    `json:"title,omitempty" bson:"title"`
	MediaURL      string    `json:"media_url,omitempty" bson:"media_url"`
	PhotoURL      string    `json:"photo_url,omitempty" bson:"photo_url"`
	HighDef1      string    `json:"1080p,omitempty" bson:"1080p"`
	HighDef2      string    `json:"720p,omitempty" bson:"720p"`
	StandardDef1  string    `json:"540p,omitempty" bson:"540p"`
	StandardDef2  string    `json:"360p,omitempty" bson:"360p"`
	LiveStreaming string    `json:"live_streaming,omitempty" bson:"live_streaming"`
	Duration      int       `json:"duration,omitempty" bson:"duration"`
	MediaType     string    `json:"media_type,omitempty" bson:"media_type"`
	Author        string    `json:"author,omitempty" bson:"author,omitempty"`
	CreatedAT     time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type Resource struct {
	Title     string    `json:"title,omitempty" bson:"title"`
	MediaURL  string    `json:"media_url,omitempty" bson:"media_url"`
	MediaType string    `json:"media_type,omitempty" bson:"media_type"`
	Author    string    `json:"author,omitempty" bson:"author,omitempty"`
	CreatedAT time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type ArrayResponse struct {
	Contents []*Content `json:"contents" bson:"contents"`
	Count    int        `json:"count" bson:"count"`
}
