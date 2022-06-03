package usecase

import (
	"errors"
	"fmt"
	"study/features/network"
	"study/libs/notification"
	"study/models"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type networkUsecase struct {
	networkRepo network.Repository
}

// NewNetworkUsecase will create new an networkUsecase object representation of Network.go.Usecase interface
func NewNetworkUsecase(a network.Repository) network.Usecase {
	return &networkUsecase{
		networkRepo: a,
	}
}

func (a *networkUsecase) Find(id string) (interface{}, error) {
	resAccount, err := a.networkRepo.Find(id)

	var customResponse models.UserNetwork

	if err != nil {

		customResponse.Following = make([]bson.ObjectId, 0)
		customResponse.Followers = make([]bson.ObjectId, 0)

		customResponse.UserID = bson.ObjectIdHex(id)

		return customResponse, nil
	}

	if len(resAccount.Followers) == 0 {
		customResponse.ID = resAccount.ID
		customResponse.UserID = resAccount.UserID
		customResponse.Followers = make([]bson.ObjectId, 0)
		customResponse.Following = resAccount.Following
	} else if len(resAccount.Following) == 0 {
		customResponse.ID = resAccount.ID
		customResponse.UserID = resAccount.UserID
		customResponse.Followers = resAccount.Followers
		customResponse.Following = make([]bson.ObjectId, 0)
	} else {
		customResponse.ID = resAccount.ID
		customResponse.UserID = resAccount.UserID
		customResponse.Followers = resAccount.Followers
		customResponse.Following = resAccount.Following
	}

	return customResponse, nil
}

func (a *networkUsecase) FindBy(key string, value string) (*models.Network, error) {
	resAccount, err := a.networkRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *networkUsecase) FindAll() ([]*models.Network, error) {
	resAccount, err := a.networkRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *networkUsecase) Create(m *models.Network) (interface{}, error) {

	existedAccount, _ := a.FindBy("_id", m.UserID.Hex())

	if existedAccount != nil {

		return existedAccount, nil
	}

	m.ID = bson.NewObjectId()

	m.CreatedAT = time.Now()

	m.UpdatedAt = time.Now()

	res, err := a.networkRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *networkUsecase) Update(m interface{}, id string) (interface{}, error) {
	existedAccount, _ := a.FindBy("_id", id)

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	_, err := a.networkRepo.Update(m, id)

	if err != nil {
		return nil, err
	}

	newAccount, _ := a.FindBy("_id", id)

	return newAccount, nil
}

func (a *networkUsecase) FollowAUser(id string, user string) (interface{}, error) {
	existedAccount, _ := a.networkRepo.FindBy("_id", id)

	if existedAccount == nil {
		var m models.Network

		m.UserID = bson.ObjectIdHex(id)

		m.Following = append(m.Following, bson.ObjectIdHex(user))

		m.ID = bson.NewObjectId()

		m.CreatedAT = time.Now()

		m.UpdatedAt = time.Now()

		res, er := a.networkRepo.Create(m)

		if er != nil {
			return nil, er
		}

		// Update Users Account
		userAccount, _ := a.FindBy("_id", user)

		if userAccount == nil {
			var m models.Network

			m.UserID = bson.ObjectIdHex(user)

			m.Followers = append(m.Followers, bson.ObjectIdHex(id))

			m.ID = bson.NewObjectId()

			m.CreatedAT = time.Now()

			m.UpdatedAt = time.Now()

			_, er := a.networkRepo.Create(m)

			if er != nil {
				return nil, er
			}

			// Get User Detail To Send Push Notification
			r, _ := a.networkRepo.FindUser(id, "accounts")

			// Get User Detail To Send Push Notification
			u, _ := a.networkRepo.FindUser(user, "accounts")
			//Send Push Notification
			notification.SendPushNotification(r.FirstName+" "+r.LastName, "just followed you", u.FirebaseTokens, r.ID.Hex(), r, "network")

			//Construct Notify Json And Store In DB

			var send models.Notification

			send.ID = bson.NewObjectId()
			send.UserID = bson.ObjectIdHex(user)
			send.SenderID = bson.ObjectIdHex(id)
			send.Type = "network"
			send.Topic = "Someone Followed You"
			send.Message = fmt.Sprintf("You Have Been Followed By %s %s", r.FirstName, r.LastName)
			send.CreatedAT = time.Now()
			send.UpdatedAt = time.Now()

			notify, _ := a.networkRepo.Notify(send)

			if notify != nil {
				fmt.Println(notify)
			}

			return res, nil
		}

		List := append(userAccount.Followers, bson.ObjectIdHex(id))

		userUpdate := bson.M{"followers": List}

		_, errr := a.Update(userUpdate, user)

		if errr != nil {
			return nil, errr
		}

		// Get User Detail To Send Push Notification
		r, _ := a.networkRepo.FindUser(id, "accounts")

		// Get User Detail To Send Push Notification
		u, _ := a.networkRepo.FindUser(user, "accounts")
		//Send Push Notification
		notification.SendPushNotification(r.FirstName+" "+r.LastName, "just followed you", u.FirebaseTokens, r.ID.Hex(), r, "network")

		//Construct Notify Json And Store In DB

		var send models.Notification

		send.ID = bson.NewObjectId()
		send.UserID = bson.ObjectIdHex(user)
		send.SenderID = bson.ObjectIdHex(id)
		send.Type = "network"
		send.Topic = "Someone Followed You"
		send.Message = fmt.Sprintf("You Have Been Followed By %s %s", r.FirstName, r.LastName)
		send.CreatedAT = time.Now()
		send.UpdatedAt = time.Now()

		notify, _ := a.networkRepo.Notify(send)

		if notify != nil {
			fmt.Println(notify)
		}
		return res, nil
	}

	query := bson.M{"user_id": bson.ObjectIdHex(id), "following": bson.ObjectIdHex(user)}

	interact, _ := a.networkRepo.FindWith(query)

	if interact != nil {
		return nil, errors.New("You Have Already Followed This User")
	}

	newList := append(existedAccount.Following, bson.ObjectIdHex(user))

	update := bson.M{"following": newList}

	res, err := a.Update(update, id)

	if err != nil {
		return nil, err
	}

	// Update Users Account
	userAccount, _ := a.FindBy("_id", user)

	if userAccount == nil {
		var m models.Network

		m.UserID = bson.ObjectIdHex(user)

		m.Followers = append(m.Followers, bson.ObjectIdHex(id))

		m.ID = bson.NewObjectId()

		m.CreatedAT = time.Now()

		m.UpdatedAt = time.Now()

		_, er := a.networkRepo.Create(m)

		if er != nil {
			return nil, er
		}
		// Get User Detail To Send Push Notification
		r, _ := a.networkRepo.FindUser(id, "accounts")

		// Get User Detail To Send Push Notification
		u, _ := a.networkRepo.FindUser(user, "accounts")

		//Send Push Notification
		notification.SendPushNotification(r.FirstName+" "+r.LastName, "just followed you", u.FirebaseTokens, r.ID.Hex(), r, "network")

		//Construct Notify Json And Store In DB

		var send models.Notification

		send.ID = bson.NewObjectId()
		send.UserID = bson.ObjectIdHex(user)
		send.SenderID = bson.ObjectIdHex(id)
		send.Type = "network"
		send.Topic = "Someone Followed You"
		send.Message = fmt.Sprintf("You Have Been Followed By %s %s", r.FirstName, r.LastName)
		send.CreatedAT = time.Now()
		send.UpdatedAt = time.Now()

		notify, _ := a.networkRepo.Notify(send)

		if notify != nil {
			fmt.Println(notify)
		}

		return res, nil
	}

	List := append(userAccount.Followers, bson.ObjectIdHex(id))

	userUpdate := bson.M{"followers": List}

	_, er := a.networkRepo.Update(userUpdate, user)

	if er != nil {
		return nil, er
	}

	// Get User Detail To Send Push Notification
	r, _ := a.networkRepo.FindUser(id, "accounts")

	// Get User Detail To Send Push Notification
	u, _ := a.networkRepo.FindUser(user, "accounts")
	//Send Push Notification
	notification.SendPushNotification(r.FirstName+" "+r.LastName, "just followed you", u.FirebaseTokens, r.ID.Hex(), r, "network")

	//Construct Notify Json And Store In DB

	var send models.Notification

	send.ID = bson.NewObjectId()
	send.UserID = bson.ObjectIdHex(user)
	send.SenderID = bson.ObjectIdHex(id)
	send.Type = "network"
	send.Topic = "Someone Followed You"
	send.Message = fmt.Sprintf("You Have Been Followed By %s %s", r.FirstName, r.LastName)
	send.CreatedAT = time.Now()
	send.UpdatedAt = time.Now()

	notify, _ := a.networkRepo.Notify(send)

	if notify != nil {
		fmt.Println(notify)
	}

	//Returns Original User Payload
	newAccount, _ := a.FindBy("_id", id)

	return newAccount, nil
}

func (a *networkUsecase) UnFollowAUser(id string, user string) (interface{}, error) {
	existedAccount, _ := a.FindBy("_id", id)

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	s := existedAccount.Following

	for i, v := range s {
		if v == bson.ObjectIdHex(user) {
			s = append(s[:i], s[i+1:]...)
			break
		}
	}

	update := bson.M{"following": s}

	_, err := a.networkRepo.Update(update, id)

	if err != nil {
		return nil, err
	}

	//Update Other Users Followers Account
	Account, _ := a.FindBy("_id", user)

	if Account == nil {
		return nil, models.ErrNotFound
	}

	t := Account.Followers

	for i, v := range t {
		if v == bson.ObjectIdHex(id) {
			t = append(t[:i], t[i+1:]...)
			break
		}
	}

	userUpdate := bson.M{"followers": t}

	_, er := a.networkRepo.Update(userUpdate, user)

	if er != nil {
		return nil, er
	}

	//Send Updated Initial users Network
	newAccount, _ := a.FindBy("_id", id)

	return newAccount, nil
}

func (a *networkUsecase) FindFollowers(id string, skip int, limit int) (interface{}, error) {

	resAccount, err := a.networkRepo.FindFollowers(id, skip, limit)

	if err != nil {
		customResponse := make(map[string]interface{})

		customResponse["followers"] = make([]bson.ObjectId, 0)

		customResponse["user_id"] = bson.ObjectIdHex(id)

		customResponse["count"] = "0"

		return customResponse, nil
	}

	var listFollowers models.Followers

	var data []models.Follow

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(resAccount)
	bson.Unmarshal(bsonBytes, &listFollowers)

	for _, x := range listFollowers.Followers {

		// Check If User Has Followed Specified Follower
		query := bson.M{"user_id": bson.ObjectIdHex(id), "following": bson.ObjectIdHex(x.ID.Hex())}

		resp, _ := a.networkRepo.FindWith(query)

		if resp != nil {
			x.IsFollowed = true
		} else {
			x.IsFollowed = false
		}
		data = append(data, x)
	}

	listFollowers.Followers = data

	return listFollowers, nil
}

func (a *networkUsecase) FindFollowing(id string, skip int, limit int) (interface{}, error) {

	resAccount, err := a.networkRepo.FindFollowing(id, skip, limit)

	if err != nil {
		customResponse := make(map[string]interface{})

		customResponse["following"] = make([]bson.ObjectId, 0)

		customResponse["user_id"] = bson.ObjectIdHex(id)

		customResponse["count"] = "0"

		return customResponse, nil
	}

	var listFollowing models.Following

	var data []models.Follow

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(resAccount)
	bson.Unmarshal(bsonBytes, &listFollowing)

	for _, x := range listFollowing.Following {
		x.IsFollowed = true
		data = append(data, x)
	}

	listFollowing.Following = data

	return listFollowing, nil
}

func (a *networkUsecase) GetUsersSuggestedFollowers(id string, skip int, limit int) (interface{}, error) {

	resAccount, errr := a.networkRepo.FindUser(id, "accounts")

	if errr != nil {
		customResponse := make(map[string]interface{})

		customResponse["followers"] = make([]bson.ObjectId, 0)

		customResponse["count"] = "0"

		return customResponse, nil
	}

	resFollowers, err := a.networkRepo.GetUsersSuggestedFollowers(resAccount.InterestID, skip, limit)

	if err != nil {
		customResponse := make(map[string]interface{})

		customResponse["followers"] = make([]bson.ObjectId, 0)

		customResponse["count"] = "0"

		return customResponse, nil
	}

	var listFollowing models.Following

	var data []models.Follow

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(resFollowers)
	bson.Unmarshal(bsonBytes, &listFollowing)

	for _, x := range listFollowing.Following {
		// Check If User Has Followed Specified Follower
		query := bson.M{"user_id": bson.ObjectIdHex(id), "following": bson.ObjectIdHex(x.ID.Hex())}

		resp, _ := a.networkRepo.FindWith(query)

		if resp == nil && bson.ObjectIdHex(id) != x.ID {
			data = append(data, x)
		}

	}

	listFollowing.Following = data

	return listFollowing, nil
}
