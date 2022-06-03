package chat

import (
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

//Usecase interface represents chat usecases
type Usecase interface {
	// Chat Usecases
	CreateDiscussion(interface{}) (interface{}, error)
	Find(id string) (map[string]interface{}, error)
	FindBy(key string, value string) (map[string]interface{}, error)
	FindAll(match bson.M, match2 bson.M) ([]map[string]interface{}, error)
	FindByCourse(string) ([]map[string]interface{}, error)
	Update(user interface{}, id string) (interface{}, error)
	AddComment(user interface{}, id bson.ObjectId, objectId bson.ObjectId) (interface{}, error)
	FindCommentsBy(key string, value string) ([]map[string]interface{}, error)
	ReadData(conn *websocket.Conn)
	BroadcastMessage(store map[string]interface{})
	FindDiscussionByUserID(match bson.M) ([]map[string]interface{}, error)
	FindAllAndGroup() ([]interface{}, error)
	TrendingDiscussions() ([]map[string]interface{}, error)
	UpdateDiscussionUsers(user string, id string) (interface{}, error)
	Subscribe(string, string) (interface{}, error)
	UnSubscribe(string, string) (interface{}, error)
	DeleteComment(string, string) (string, error)
	DeleteDiscussion(string, string) (string, error)
	UpdateComment(m interface{}, id string) (interface{}, error)
}
