package usecase

import (
	"fmt"
	"log"
	chat "study/features/discuss"
	"study/libs/notification"
	"study/models"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type chatUsecase struct {
	chatRepo chat.Repository
	ch       chan messageModel
}

type messageModel struct {
	CreatedBy interface{} `json:"created_by"`
	UserID    string      `json:"user_id"`
	Group     bool        `json:"group"`
	Message   string      `json:"message"`
	Receiver  string      `json:"receiver"`
	Room      string      `json:"room"`
	Time      time.Time   `json:"time"`
}

// NewAccountUsecase will create new an accountUsecase object representation of Usecase interface
func NewAccountUsecase(a chat.Repository) chat.Usecase {
	ch := make(chan messageModel)
	return &chatUsecase{
		chatRepo: a,
		ch:       ch,
	}
}

func (a *chatUsecase) Find(id string) (map[string]interface{}, error) {

	resAccount, err := a.chatRepo.Find(id, "chats")

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *chatUsecase) FindBy(key string, value string) (map[string]interface{}, error) {
	resAccount, err := a.chatRepo.FindBy(key, value, "chats")

	if err != nil {
		return nil, err
	}

	return resAccount, nil
}

func (a *chatUsecase) FindByCourse(id string) ([]map[string]interface{}, error) {
	match := bson.M{"$match": bson.M{"course_id": bson.ObjectIdHex(id)}}

	resAccount, err := a.chatRepo.FindAll(match)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *chatUsecase) FindAll(match bson.M, query bson.M) ([]map[string]interface{}, error) {
	resAccount, err := a.chatRepo.FindAll(match)

	if err != nil {
		return nil, err
	}

	res, _ := a.chatRepo.FindAll(query)

	if res == nil {
		return resAccount, nil
	}

	var data []map[string]interface{}

	for _, value := range res {
		data = append(resAccount, value)
	}

	return data, nil
}

func (a *chatUsecase) CreateDiscussion(m interface{}) (interface{}, error) {

	x := m.(*models.Chat)

	if !x.ValidChat() {
		return &models.Chat{}, errors.New("invalid chat type")
	}

	x.ID = bson.NewObjectId()

	x.CreatedAT = time.Now()

	x.UpdatedAt = time.Now()

	if len(x.Users) != 0 {
		x.Subscribers = x.Users
		x.Subscribers = append(x.Subscribers, x.CreatedBy)
	}

	_, err := a.chatRepo.Create(x, "chats")

	if err != nil {
		return nil, err
	}

	newAccount, errr := a.Find(x.ID.Hex())

	if errr != nil {
		return nil, errr
	}
	// Get User That Sent Comment
	u, _ := a.chatRepo.FindUser(x.CreatedBy.Hex())

	//Loop And Send Push Notifications
	for _, value := range x.Users {

		// Get User Detail To Send Push Notification
		r, _ := a.chatRepo.FindUser(value.Hex())
		//Get Discussion

		//Send Push Notification
		notification.SendPushNotification(u.FirstName+" "+u.LastName, "just added you to "+"'"+x.Topic+"'", r.FirebaseTokens, x.ID.Hex(), newAccount, "discuss")
	}

	return newAccount, nil
}

func (a *chatUsecase) Update(m interface{}, id string) (interface{}, error) {
	existedAccount, _ := a.FindBy("_id", id)

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	_, err := a.chatRepo.Update(m, id)

	if err != nil {
		return nil, err
	}

	res, _ := a.FindBy("_id", id)

	return res, nil
}

func (a *chatUsecase) AddComment(m interface{}, id bson.ObjectId, discussion bson.ObjectId) (interface{}, error) {

	//First Add User To Discussion Channel
	chat := m.(*models.Comment)

	query := bson.M{"users": bson.ObjectIdHex(chat.CreatedBy.Hex()), "_id": bson.ObjectIdHex(discussion.Hex())}

	interact, _ := a.chatRepo.FindWith(query)

	if interact == nil {
		_, e := a.UpdateDiscussionUsers(chat.CreatedBy.Hex(), discussion.Hex())
		if e != nil {
			return nil, e
		}
	}

	//Go Ahead And Add comment to Chat
	_, err := a.chatRepo.Create(m, "comments")

	if err != nil {
		return nil, err
	}

	existedAccount, er := a.chatRepo.FindSingle("_id", discussion.Hex())

	if er != nil {
		return nil, er
	}

	x := m.(*models.Comment)

	x.UpdatedAt = time.Now()

	var y *models.Chat

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(existedAccount)
	bson.Unmarshal(bsonBytes, &y)

	y.Comments = append(y.Comments, x.ID)

	_, errr := a.chatRepo.Update(y, discussion.Hex())

	if errr != nil {
		return nil, err
	}

	// Get User That Sent Comment
	u, _ := a.chatRepo.FindUser(x.CreatedBy.Hex())
	c, _ := a.chatRepo.Find(discussion.Hex(), "chats")

	//Loop And Send Push Notifications
	for _, value := range y.Subscribers {

		// Get User Detail To Send Push Notification
		r, _ := a.chatRepo.FindUser(value.Hex())
		//Get Discussion

		//Send Push Notification
		notification.SendPushNotification(u.FirstName+" "+u.LastName, "just commented in "+"'"+y.Topic+"'", r.FirebaseTokens, y.ID.Hex(), c, "discuss")
	}

	newAccount, errr := a.chatRepo.FindComments(id.Hex())

	if errr != nil {
		return nil, errr
	}

	return newAccount, nil
}

func (a *chatUsecase) FindCommentsBy(key string, value string) ([]map[string]interface{}, error) {
	resAccount, err := a.chatRepo.FindCommentsBy(key, value)

	if err != nil {
		return nil, err
	}

	return resAccount, nil
}

func (a *chatUsecase) FindDiscussionByUserID(match bson.M) ([]map[string]interface{}, error) {
	resAccount, err := a.chatRepo.FindAll(match)

	if err != nil {
		return nil, err
	}

	return resAccount, nil
}

func (a *chatUsecase) ReadData(conn *websocket.Conn) {
	for {
		var message messageModel
		message.Time = time.Now()
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(message.UserID)
			return
		}

		a.ch <- message
	}
}

func (a *chatUsecase) BroadcastMessage(store map[string]interface{}) {
	for message := range a.ch {
		switch message.Group {
		case true:
			sendGroupMessage(a, store, message)
		case false:
			sendSingleMessage(store, message)
		}
	}
}

func sendSingleMessage(store map[string]interface{}, message messageModel) {
	receiver, ok := store[message.Receiver]
	sender, _ := store[message.UserID]
	if !ok {
		v, ok := sender.(*websocket.Conn)

		if !ok {
			log.Println("can not convert")
		}
		log.Println(v)

		//err := v.WriteMessage(websocket.TextMessage, []byte("user is not online"))
		//if err != nil {
		//	log.Print(err)
		//}
	} else {
		v, ok := receiver.(*websocket.Conn)
		if !ok {
			log.Println("can not convert")
		}
		err := v.WriteJSON(message)
		if err != nil {
			v, ok := sender.(*websocket.Conn)
			if !ok {
				log.Println("can not convert")
			}
			log.Println(v)
			//err := v.WriteMessage(websocket.TextMessage, []byte("user is not online"))
			//if err != nil {
			//	log.Print(err)
			//}
		}
	}
}

func sendGroupMessage(a *chatUsecase, store map[string]interface{}, message messageModel) {

	room := message.Room

	res, err := a.chatRepo.FindSingle("room", room)

	var listAccount *models.Chat

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(res)
	bson.Unmarshal(bsonBytes, &listAccount)

	user, er := a.chatRepo.FindBy("_id", message.UserID, "accounts")

	if er != nil {
		log.Println("An error occured %e", er)
	}

	message.CreatedBy = user

	fmt.Println(message)

	if err != nil {
		log.Println("An error occured %e", err)
	}

	if err != nil {
		log.Println("An error occured %e", err)
	}

	onlineUsers := make([]string, len(listAccount.Users))

	for i := range listAccount.Users {
		onlineUsers[i] = listAccount.Users[i].Hex()
	}

	fmt.Println(onlineUsers)
	for _, value := range onlineUsers {
		message.Receiver = value
		sendSingleMessage(store, message)
	}
}

func (a *chatUsecase) FindAllAndGroup() ([]interface{}, error) {
	resAccount, err := a.chatRepo.FindAllAndGroup()

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *chatUsecase) UpdateDiscussionUsers(m string, id string) (interface{}, error) {
	query := bson.M{"_id": bson.ObjectIdHex(id)}

	interact, _ := a.chatRepo.FindWith(query)

	if interact == nil {
		return nil, errors.New("Invalid Discussion")
	}
	var listAccount *models.Chat

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(interact)
	bson.Unmarshal(bsonBytes, &listAccount)

	var list []bson.ObjectId

	if len(listAccount.Users) == 0 {
		list = append(list, bson.ObjectIdHex(m))

	} else {
		for _, value := range listAccount.Users {
			if value != bson.ObjectId(m) {
				list = append(listAccount.Users, bson.ObjectIdHex(m))
			}
		}
	}

	update := bson.M{"users": list}

	_, err := a.Update(update, id)

	if err != nil {
		return nil, err
	}

	res, er := a.Find(id)

	if er != nil {
		return nil, er
	}

	return res, nil
}

func (a *chatUsecase) TrendingDiscussions() ([]map[string]interface{}, error) {
	resAccount, err := a.chatRepo.Trending()

	if err != nil {
		return nil, err
	}

	return resAccount, nil
}

func (a *chatUsecase) Subscribe(id string, discussionID string) (interface{}, error) {

	query := bson.M{"subscribers": bson.ObjectIdHex(id), "_id": bson.ObjectIdHex(discussionID)}

	interact, _ := a.chatRepo.FindWith(query)

	if interact != nil {
		return nil, errors.New("User Already Subscribed")
	}
	var listAccount *models.Chat

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(interact)
	bson.Unmarshal(bsonBytes, &listAccount)

	var list []bson.ObjectId

	if len(listAccount.Subscribers) == 0 {
		list = append(list, bson.ObjectIdHex(id))
	} else {
		for _, value := range listAccount.Subscribers {
			if value != bson.ObjectId(id) {
				list = append(listAccount.Subscribers, bson.ObjectIdHex(id))
			}
		}
	}

	update := bson.M{"subscribers": list}

	res, err := a.Update(update, discussionID)

	if err != nil {
		return nil, err
	}

	res, er := a.Find(discussionID)

	if er != nil {
		return nil, er
	}

	return res, nil
}

func (a *chatUsecase) UnSubscribe(id string, discussionID string) (interface{}, error) {

	query := bson.M{"subscribers": bson.ObjectIdHex(id), "_id": bson.ObjectIdHex(discussionID)}

	interact, _ := a.chatRepo.FindWith(query)

	if interact == nil {
		return nil, errors.New("User Already UnSubscribed")
	}
	var listAccount *models.Chat

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(interact)
	bson.Unmarshal(bsonBytes, &listAccount)

	s := listAccount.Subscribers

	for i, v := range s {
		if v == bson.ObjectIdHex(id) {
			s = append(s[:i], s[i+1:]...)
			break
		}
	}

	update := bson.M{"subscribers": s}

	_, err := a.Update(update, discussionID)

	if err != nil {
		return nil, err
	}

	res, er := a.Find(discussionID)

	if er != nil {
		return nil, er
	}

	return res, nil
}

func (a *chatUsecase) DeleteComment(id string, user string) (string, error) {

	match := bson.M{"_id": bson.ObjectIdHex(id), "created_by": bson.ObjectIdHex(user)}

	err := a.chatRepo.Delete(match, "comments")

	if err != nil {
		return "", err
	}
	query := bson.M{"comments": bson.ObjectIdHex(id)}

	interact, _ := a.chatRepo.FindWith(query)

	if interact == nil {
		return "", errors.New("Comment Already Deleted")
	}

	var listAccount *models.Chat

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(interact)
	bson.Unmarshal(bsonBytes, &listAccount)

	s := listAccount.Comments

	for i, v := range s {
		if v == bson.ObjectIdHex(id) {
			s = append(s[:i], s[i+1:]...)
			break
		}
	}

	update := bson.M{"comments": s}

	_, e := a.Update(update, listAccount.ID.Hex())

	if e != nil {
		return "", e
	}
	return "Comment Removed", nil
}

func (a *chatUsecase) DeleteDiscussion(id string, user string) (string, error) {

	match := bson.M{"_id": bson.ObjectIdHex(id), "created_by": bson.ObjectIdHex(user)}

	err := a.chatRepo.Delete(match, "chats")

	fmt.Println("here", err, id, user)
	if err != nil {
		return "", err
	}

	return "Discussion Removed", nil
}

func (a *chatUsecase) UpdateComment(m interface{}, id string) (interface{}, error) {
	existedAccount, _ := a.chatRepo.FindBy("_id", id, "comments")

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	_, err := a.chatRepo.UpdateComment(m, id)

	if err != nil {
		return nil, err
	}

	res, _ := a.chatRepo.FindBy("_id", id, "comments")

	return res, nil
}
