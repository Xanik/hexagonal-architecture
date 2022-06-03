package dbs

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"study/config"
	"time"

	"gopkg.in/mgo.v2"
)

//ConnectMongodb returns a connection to a mongodb instance through a connection string in the configuration file
func ConnectMongodb() *mgo.Session {
	url := config.Env.MongoURL

	url = strings.Split(url, "?")[0]
	//connect URL:
	dialInfo, err := mgo.ParseURL(url)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), &tls.Config{})
		return conn, err
	}
	//Here is the session you are looking for. Up to you from here ;)
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Printf("db connection error: %s", err.Error())
	}
	log.Printf("db connection: %s", session)
	return session
}

//ConnectMongodbURL returns a connection to a mongodb instance through a connection string in the configuration file
func ConnectMongodbURL() *mgo.Session {
	url := config.Env.MongoURL

	url = strings.Split(url, "?")[0]

	dialInfo, err := mgo.ParseURL(url)
	params := &mgo.DialInfo{
		Username: dialInfo.Username,
		Password: dialInfo.Password,
		Addrs:    dialInfo.Addrs,
		Database: "admin",
	}
	params.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), &tls.Config{})
		return conn, err
	}
	session, errr := mgo.DialWithInfo(params)
	if errr != nil {
		fmt.Println(errr)
	}
	err = session.Ping()
	if err != nil {
		log.Printf("db connection error: %s", err.Error())
	}
	log.Printf("db connection: %s", session)
	return session
}

// TODO: Should not Fatalf if error occurs
func NewClient() *mgo.Session {

	dialInfo := getMongoDialInfo()
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	log.Printf("connected to MongoDB successfully")
	err = session.Ping()

	if err != nil {
		log.Printf("db connection error: %s", err.Error())
	}
	log.Printf("db connection: %s", session)

	return session
}

func getMongoDialInfo() *mgo.DialInfo {

	url := config.Env.MongoURL
	url = strings.Split(url, "?")[0]
	fmt.Println(url)

	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		fmt.Println(err)
	}
	dialInfo.Timeout = 5 * time.Second
	params := &mgo.DialInfo{
		Username: dialInfo.Username,
		Password: dialInfo.Password,
		Addrs:    dialInfo.Addrs,
		Database: "admin",
	}
	params.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), &tls.Config{})
		return conn, err
	}

	return params
}

var indexes = []struct {
	coll  string
	index mgo.Index
}{
	{
		coll:  "accounts",
		index: mgo.Index{Key: []string{"email"}, Unique: true},
	},
}

func ensureIndexes(db *mgo.Database) error {
	for _, index := range indexes {
		err := db.C(index.coll).EnsureIndex(index.index)
		if err != nil {
			return err
		}
	}
	return nil
}
