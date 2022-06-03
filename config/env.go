package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type env struct {
	Env              string
	Port             string
	MongoURL         string
	MongoSSL         string
	RedisURL         string
	TokenSecret      string
	MailGunDomain    string
	MailGunApiKey    string
	WebClientBaseUrl string
	ElasticsearchURL string
	ElasticEnv       string
	S3REGION         string
	S3BUCKET         string
	DBNAME           string
}

func (e env) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Env, validation.Required),
		validation.Field(&e.Port, validation.Required),
		validation.Field(&e.MongoURL, validation.Required),
		validation.Field(&e.MongoSSL, validation.Required),
		validation.Field(&e.RedisURL, validation.Required),
		validation.Field(&e.TokenSecret, validation.Required),
		validation.Field(&e.MailGunDomain, validation.Required),
		validation.Field(&e.MailGunApiKey, validation.Required),
		validation.Field(&e.WebClientBaseUrl, validation.Required),
		validation.Field(&e.ElasticsearchURL, validation.Required),
		validation.Field(&e.ElasticEnv, validation.Required),
		validation.Field(&e.S3BUCKET, validation.Required),
		validation.Field(&e.S3REGION, validation.Required),
		validation.Field(&e.DBNAME, validation.Required),
	)
}
