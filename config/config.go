package config

import (
	"os"
	"study/logger"

	"github.com/sirupsen/logrus"
)

var Env env

const defaultPort = "8082"
const defaultRedisUrl = "redis://localhost:6379"
const defaultMongoUrl = ""

func init() {
	e := getEnvironment()
	if err := e.Validate(); err != nil {
		logger.DefaultLogger.WithFields(logrus.Fields{"type": "env_error", "stack": err}).Error("Invalid environment variables")
	}

	Env = e
}

func getEnvironment() env {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = defaultRedisUrl
	}

	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		mongoUrl = defaultMongoUrl
	}

	return env{
		Env:              os.Getenv("ENV"),
		Port:             port,
		MongoURL:         mongoUrl,
		MongoSSL:         os.Getenv("MONGO_SSL"),
		RedisURL:         redisUrl,
		TokenSecret:      os.Getenv("TOKEN_SECRET"),
		MailGunDomain:    os.Getenv("MAILGUN_DOMAIN"),
		MailGunApiKey:    os.Getenv("MAILGUN_API_KEY"),
		WebClientBaseUrl: os.Getenv("WEB_CLIENT_BASE_URL"),
		ElasticsearchURL: "https://elastic:v9hKmmYyQg3dtht2lz759LGM@31b6486899c74ad095db11db761e5e7c.eu-west-1.aws.found.io:9243/",
		ElasticEnv:       os.Getenv("ELASTICSEARCH_ENV"),
		S3REGION:         os.Getenv("S3REGION"),
		S3BUCKET:         os.Getenv("S3BUCKET"),
		DBNAME:           os.Getenv("DBNAME"),
	}
}
