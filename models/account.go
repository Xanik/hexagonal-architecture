package models

import (
	"regexp"
	"time"
	"unicode/utf8"

	"gopkg.in/mgo.v2/bson"
)

//Account model
type Account struct {
	ID             bson.ObjectId   `json:"_id" bson:"_id,omitempty"`
	FirstName      string          `json:"first_name" bson:"first_name"`
	LastName       string          `json:"last_name" bson:"last_name"`
	Email          string          `json:"email,omitempty" bson:"email,omitempty"`
	Phone          string          `json:"phone,omitempty" bson:"phone,omitempty"`
	Password       string          `json:"password,omitempty" bson:"password,omitempty"`
	Bio            string          `json:"bio" bson:"bio"`
	Code           int             `json:"code" bson:"code"`
	Image          string          `json:"image" bson:"image"`
	Gender         string          `json:"gender,omitempty" bson:"gender,omitempty"`
	Followers      bson.ObjectId   `json:"followers,omitempty" bson:"followers,omitempty"`
	Following      bson.ObjectId   `json:"followering,omitempty" bson:"followering,omitempty"`
	Step           string          `json:"step,omitempty" bson:"step,omitempty"`
	Validated      bool            `json:"validated" bson:"validated"`
	AccountType    string          `json:"account_type,omitempty" bson:"account_type,omitempty"`
	Institution    bson.ObjectId   `json:"institution,omitempty" bson:"institution,omitempty"`
	InterestID     []bson.ObjectId `json:"interest_id,omitempty" bson:"interest_id,omitempty"`
	OrganizationID bson.ObjectId   `json:"organization_id" bson:"organization_id"`
	FirebaseTokens string          `json:"firebase_tokens" bson:"firebase_tokens"`
	CreatedWITH    string          `json:"created_with" bson:"created_with"`
	CreatedAT      time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" bson:"updated_at"`
	LastLogin      time.Time       `json:"last_login" bson:"last_login"`
}

//Student AccountType Options
const (
	Student    string = "student"
	Lecturer   string = "lecturer"
	Researcher string = "researcher"
	Industry   string = "organization"
)

//Login model
type Login struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Email    string        `json:"email" bson:"email"`
	Phone    string        `json:"phone,omitempty" bson:"phone,omitempty"`
	Password string        `json:"password,omitempty" bson:"password,omitempty"`
}

//Password model
type Password struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	NewPassword string        `json:"newpassword,omitempty" bson:"newpassword,omitempty"`
	OldPassword string        `json:"oldpassword,omitempty" bson:"oldpassword,omitempty"`
}

//FeedBack model
type FeedBack struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	Email       string        `json:"email" bson:"email"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
}

//Notification model
type Notification struct {
	ID        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Type      string        `json:"type" bson:"type"`
	Topic     string        `json:"topic" bson:"topic"`
	Message   string        `json:"message" bson:"message"`
	UserID    bson.ObjectId `json:"user_id" bson:"user_id"`
	SenderID  bson.ObjectId `json:"sender_id" bson:"sender_id"`
	ContentID bson.ObjectId `json:"content_id,omitempty" bson:"content_id,omitempty"`
	CreatedAT time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

//Valid Checks For Student AccountType
func (c Account) Valid() bool {
	// If it's a user input, you'd want to validate AccountType's underlying
	// string isn't out of the enum's range.
	if c.AccountType == Student || c.AccountType == Lecturer || c.AccountType == Researcher || c.AccountType == Industry {
		return true
	}
	return false
}

//Valid Checks For Empty Fields
func (c Account) ValidFieldWithEmail() bool {
	// If it's a user input, you'd want to validate the empty fields
	if c.FirstName == "" || c.LastName == "" || c.Email == "" || c.Institution == "" || c.OrganizationID == "" || c.Gender == "" {
		return false
	}

	return true
}

//Valid Checks For Empty Fields
func (c Account) ValidFieldWithPhone() bool {
	// If it's a user input, you'd want to validate the empty fields
	if c.FirstName == "" || c.LastName == "" || c.Phone == "" || c.Institution == "" || c.OrganizationID == "" || c.Gender == "" {
		return false
	}
	return true
}

func (c Account) ValidEmail() bool {

	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	return Re.MatchString(c.Email)

}

func (c Account) ValidPhone() bool {

	Re := regexp.MustCompile(`(?:^|[^0-9])(0[0-9][0-9]{9})(?:$|[^0-9])`)

	return Re.MatchString(c.Phone)

}

func (c Account) TransformPhone() string {
	s := c.Phone
	_, i := utf8.DecodeRuneInString(s)
	return "234" + s[i:]
}

func (c Account) ValidPassword() bool {

	pwdLen := len(c.Password)
	if pwdLen < 5 || c.Password == "" {
		return false
	}
	return true
}

//Validate Login
func (c Login) ValidPhone() bool {

	Re := regexp.MustCompile(`(?:^|[^0-9])(0[0-9][0-9]{9})(?:$|[^0-9])`)

	return Re.MatchString(c.Phone)

}

func (c Login) TransformPhone() string {
	s := c.Phone
	_, i := utf8.DecodeRuneInString(s)
	return "234" + s[i:]
}
