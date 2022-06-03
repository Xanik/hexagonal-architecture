package usecase

import (
	"errors"
	"fmt"
	"math/rand"
	interaction "study/features/interactions"
	"study/libs/notification"
	"study/models"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type interactionUsecase struct {
	interactionRepo interaction.Repository
}

// NewAccountUsecase will create new an accountUsecase object representation of account.Usecase interface
func NewAccountUsecase(a interaction.Repository) interaction.Usecase {
	return &interactionUsecase{
		interactionRepo: a,
	}
}

func GetColor() string {
	COLORs := []string{
		"#F5BC27", "#21607A", "#2D2118", "#D12F62", "#4A2046", "#E67A2B", "#14A749",
	}

	// randomly pick one beautiful girl(URL) from the URLs slice
	rand.Seed(time.Now().UnixNano())

	choosen := COLORs[rand.Intn(len(COLORs)-1)]

	fmt.Println("Randomly selected this Color : ", choosen)
	return choosen
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (a *interactionUsecase) Index(content string) {
	// Index Document In Elastic
	indexContent, _ := a.interactionRepo.GetIndexedContent(content)

	if indexContent != nil {
		a.interactionRepo.IndexDocument(indexContent, "contents", indexContent.ID.Hex())
		//Index Completed
	}

}

func (a *interactionUsecase) Find(id string) (interface{}, error) {
	resAccount, err := a.interactionRepo.Find(id, "interactions")

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *interactionUsecase) FindBy(key string, value string) (interface{}, error) {
	resAccount, err := a.interactionRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *interactionUsecase) FindByCollection(collection string, skip int, limit int, id string) (interface{}, error) {

	//Get Collections With Content
	resAccount, _ := a.interactionRepo.FindByCollection(collection, skip, limit, id)

	q := bson.M{"type": "collection", "collection_id": bson.ObjectIdHex(collection), "user_id": bson.ObjectIdHex(id), "content_id": bson.M{"$ne": nil}}

	existedCollection, _ := a.interactionRepo.FindWith(q)

	if existedCollection == nil {
		q := bson.M{"type": "collection", "collection_id": bson.ObjectIdHex(collection), "user_id": bson.ObjectIdHex(id)}

		Collection, _ := a.interactionRepo.FindWith(q)

		if Collection != nil {
			x := Collection.(*models.Interactions)

			// Return Custom Response And Add Empty Fields
			customResponse := make(map[string]interface{})
			customResponse["contents"] = make([]interface{}, 0)
			customResponse["count"] = 0
			customResponse["name"] = x.Name
			customResponse["color"] = x.Color
			customResponse["collection_id"] = x.CollectionID

			return customResponse, nil
		}
		return nil, errors.New("Invalid User CollectionID")

	}

	query := bson.M{"user_id": bson.ObjectIdHex(id), "type": "collection", "collection_id": bson.ObjectIdHex(collection), "content_id": bson.M{"$ne": nil}}

	count, _ := a.interactionRepo.CountContent(query)

	var listInteraction models.Collection

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(resAccount)
	bson.Unmarshal(bsonBytes, &listInteraction)

	var listContent []models.Content

	for _, value := range listInteraction.Content {
		//Check if user has bookmarked content
		isSave := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(value.ID.Hex())}

		save, _ := a.interactionRepo.FindWith(isSave)

		if save != nil {
			value.IsSaved = true
		}

		// Check If User Has Liked Content
		query := bson.M{"type": "like", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(value.ID.Hex())}

		resp, _ := a.interactionRepo.FindWith(query)

		if resp != nil {
			value.IsLiked = true
		}

		listContent = append(listContent, value)
	}

	content := make(map[string]interface{})

	content["contents"] = listContent

	content["name"] = listInteraction.Name
	content["color"] = listInteraction.Color
	content["collection_id"] = collection

	content["count"] = count

	return content, nil
}

func (a *interactionUsecase) FindAll(id string, skip int, limit int) (interface{}, error) {
	match := bson.M{"$match": bson.M{"user_id": bson.M{"$ne": bson.ObjectIdHex(id)}, "type": "like",
		"created_at": bson.M{
			"$gt": time.Now().Add(-24 * time.Hour)},
	}}

	resAccount, err := a.interactionRepo.FindAll(match, skip, limit)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *interactionUsecase) FindMyLikedContent(id string, user string, skip int, limit int) (interface{}, error) {
	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id), "type": "like"}}

	resAccount, err := a.interactionRepo.FindAll(match, skip, limit)

	if err != nil {
		customResponse := make(map[string]interface{})

		customResponse["contents"] = make([]bson.ObjectId, 0)

		customResponse["count"] = 0

		return customResponse, nil
	}

	x := resAccount.(*models.Data)

	var data []models.Content

	for _, x := range x.Content {

		// Check If User Has Content And Like
		query := bson.M{"type": "like", "user_id": bson.ObjectIdHex(user), "content_id": bson.ObjectIdHex(x.ID.Hex())}

		resp, _ := a.interactionRepo.FindWith(query)

		if resp != nil {
			x.IsLiked = true
		}
		// Check If User Has Bookmarked Content And Save
		query1 := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(user), "content_id": bson.ObjectIdHex(x.ID.Hex())}

		resp1, _ := a.interactionRepo.FindWith(query1)

		if resp1 != nil {
			x.IsSaved = true
		}
		data = append(data, x)
	}

	query := bson.M{"user_id": bson.ObjectIdHex(id), "type": "like"}

	count, _ := a.interactionRepo.CountContent(query)

	content := make(map[string]interface{})

	content["contents"] = data

	content["count"] = count

	return content, nil
}

func (a *interactionUsecase) FindMySharedContent(id string, user string, skip int, limit int) (interface{}, error) {
	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id), "type": "share"}}

	resAccount, err := a.interactionRepo.FindAll(match, skip, limit)

	if err != nil {
		customResponse := make(map[string]interface{})

		customResponse["contents"] = make([]bson.ObjectId, 0)

		customResponse["count"] = 0

		return customResponse, nil
	}

	x := resAccount.(*models.Data)

	var data []models.Content

	for _, x := range x.Content {

		// Check If User Has Liked Content And Like
		query1 := bson.M{"type": "like", "user_id": bson.ObjectIdHex(user), "content_id": bson.ObjectIdHex(x.ID.Hex())}

		resp1, _ := a.interactionRepo.FindWith(query1)

		if resp1 != nil {
			x.IsSaved = true
		} else {
			x.IsSaved = false
		}
		// Check If User Has Bookmarked Content And Save
		query := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(user), "content_id": bson.ObjectIdHex(x.ID.Hex())}

		resp, _ := a.interactionRepo.FindWith(query)

		if resp != nil {
			x.IsSaved = true
		} else {
			x.IsSaved = false
		}
		data = append(data, x)
	}

	query := bson.M{"user_id": bson.ObjectIdHex(id), "type": "share"}

	count, e := a.interactionRepo.CountContent(query)

	if e != nil {
		return nil, e
	}

	content := make(map[string]interface{})

	content["contents"] = data

	content["count"] = count

	return content, nil
}

func (a *interactionUsecase) FindMyBookmarkedContent(id string, user string, skip int, limit int) (interface{}, error) {
	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id), "type": "bookmark"}}

	resAccount, err := a.interactionRepo.FindAll(match, skip, limit)

	if err != nil {
		customResponse := make(map[string]interface{})

		customResponse["contents"] = make([]bson.ObjectId, 0)

		customResponse["count"] = 0

		return customResponse, nil
	}

	x := resAccount.(*models.Data)

	var data []models.Content

	for _, x := range x.Content {

		// Check If User Has Content And Like
		query := bson.M{"type": "like", "user_id": bson.ObjectIdHex(user), "content_id": bson.ObjectIdHex(x.ID.Hex())}

		resp, _ := a.interactionRepo.FindWith(query)

		if resp != nil {
			x.IsLiked = true
		}

		// Check If User Has Bookmarked Content And Save
		query1 := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(user), "content_id": bson.ObjectIdHex(x.ID.Hex())}

		resp1, _ := a.interactionRepo.FindWith(query1)

		if resp1 != nil {
			x.IsSaved = true
		}
		data = append(data, x)
	}

	query := bson.M{"user_id": bson.ObjectIdHex(id), "type": "bookmark"}

	count, _ := a.interactionRepo.CountContent(query)

	content := make(map[string]interface{})

	content["contents"] = data

	content["count"] = count

	return content, nil
}

func (a *interactionUsecase) FindMyCollection(id string, skip int, limit int) ([]interface{}, error) {
	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id), "type": "collection"}}

	resAccount, err := a.interactionRepo.FindByAndGroup(match, skip, limit)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *interactionUsecase) ListCollections(id string, skip int, limit int) ([]interface{}, error) {
	match := bson.M{"$match": bson.M{"user_id": bson.ObjectIdHex(id), "type": "collection"}}

	resAccount, err := a.interactionRepo.FindByAndGroupCount(match, skip, limit)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *interactionUsecase) UnLikeContent(account string, content string) (interface{}, error) {

	fmt.Println(account, content)

	query := bson.M{"type": "like", "user_id": bson.ObjectIdHex(account), "content_id": bson.ObjectIdHex(content)}

	interact, e := a.interactionRepo.FindWith(query)

	if e != nil {
		return "", e
	}

	var listInteraction models.Interactions

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(interact)
	bson.Unmarshal(bsonBytes, &listInteraction)

	resContent, er := a.interactionRepo.Find(listInteraction.ContentID.Hex(), "contents")

	if er != nil {
		return "", er
	}

	// Index Document In Elastic
	go a.Index(content)

	var listContent models.Content

	// convert bson to struct
	bsonByte, _ := bson.Marshal(resContent)
	bson.Unmarshal(bsonByte, &listContent)

	update := bson.M{"likes": listContent.Likes - 1}

	_, errr := a.interactionRepo.Update(update, listInteraction.ContentID.Hex(), "contents")

	if errr != nil {
		return "", errr
	}

	err := a.interactionRepo.Delete(account, content, "like")

	if err != nil {
		return "", err
	}

	// Check If User Has Bookmarked Content
	query1 := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(account), "content_id": bson.ObjectIdHex(content)}

	resp, _ := a.interactionRepo.FindWith(query1)

	var list models.Content

	// convert bson to struct
	bsonByte1, _ := bson.Marshal(resContent)
	bson.Unmarshal(bsonByte1, &list)

	if resp != nil {
		list.IsSaved = true
	} else {
		list.IsSaved = false
	}

	list.Likes = list.Likes - 1

	return list, nil
}

func (a *interactionUsecase) DeleteBookmark(account string, content string) (interface{}, error) {

	resContent, er := a.interactionRepo.Find(content, "contents")

	if er != nil {
		return "", er
	}

	err := a.interactionRepo.Delete(account, content, "bookmark")

	if err != nil {
		return "", err
	}

	var Date models.Content

	// convert bson to struct
	bsonByt, _ := bson.Marshal(resContent)
	bson.Unmarshal(bsonByt, &Date)

	update := bson.M{"saves": Date.Saves - 1}

	_, errr := a.interactionRepo.Update(update, content, "contents")

	if errr != nil {
		return "", errr
	}

	// Check If User Has Liked Content
	query := bson.M{"type": "like", "user_id": bson.ObjectIdHex(account), "content_id": bson.ObjectIdHex(content)}

	resp, _ := a.interactionRepo.FindWith(query)

	var listContent models.Content

	// convert bson to struct
	bsonByte, _ := bson.Marshal(resContent)
	bson.Unmarshal(bsonByte, &listContent)

	if resp != nil {
		listContent.IsSaved = true
	} else {
		listContent.IsSaved = false
	}

	// Index Document In Elastic
	go a.Index(listContent.ID.Hex())

	listContent.Saves = listContent.Saves - 1

	return listContent, nil
}

func (a *interactionUsecase) Create(m interface{}) (interface{}, error) {
	res, err := a.interactionRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *interactionUsecase) CreateEmptyCollection(m *models.Interactions) (interface{}, error) {

	query := bson.M{"type": "collection", "name": m.Name, "user_id": bson.ObjectIdHex(m.UserID.Hex())}

	interact, _ := a.interactionRepo.FindWith(query)

	if interact != nil {
		return nil, errors.New("Collection Has Already Been Created")
	}

	m.ID = bson.NewObjectId()

	m.CollectionID = bson.NewObjectId()

	m.CreatedAT = time.Now()

	m.UpdatedAt = time.Now()

	m.Color = GetColor()

	res, err := a.interactionRepo.Create(m)

	if err != nil {
		return nil, err
	}

	y := res.(*models.Interactions)

	query1 := bson.M{"user_id": bson.ObjectIdHex(m.UserID.Hex()), "type": "collection", "collection_id": bson.ObjectIdHex(y.CollectionID.Hex()), "content_id": bson.M{"$ne": nil}}

	count, _ := a.interactionRepo.CountContent(query1)

	y.Count = count

	return y, nil
}

func (a *interactionUsecase) CreateCollection(m *models.Interactions) (interface{}, error) {

	query := bson.M{"type": "collection", "name": m.Name, "user_id": bson.ObjectIdHex(m.UserID.Hex()), "content_id": bson.ObjectIdHex(m.ContentID.Hex())}

	interact, _ := a.interactionRepo.FindWith(query)

	if interact != nil {
		return nil, errors.New("Content Has Already Been Added To Collection")
	}
	q := bson.M{"type": "collection", "name": m.Name, "user_id": bson.ObjectIdHex(m.UserID.Hex())}

	existedCollection, _ := a.interactionRepo.FindWith(q)

	if existedCollection != nil {

		x := existedCollection.(*models.Interactions)

		m.ID = bson.NewObjectId()

		m.CollectionID = x.CollectionID

		m.Name = x.Name

		m.CreatedAT = time.Now()

		m.UpdatedAt = time.Now()

		m.Color = x.Color

		_, er := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

		if er != nil {
			return nil, er
		}

		res, err := a.interactionRepo.Create(m)

		if err != nil {
			return nil, err
		}

		y := res.(*models.Interactions)

		query := bson.M{"user_id": bson.ObjectIdHex(m.UserID.Hex()), "type": "collection", "collection_id": bson.ObjectIdHex(y.CollectionID.Hex()), "content_id": bson.M{"$ne": nil}}

		count, _ := a.interactionRepo.CountContent(query)

		y.Count = count

		return y, nil
	}

	m.ID = bson.NewObjectId()

	m.CollectionID = bson.NewObjectId()

	m.CreatedAT = time.Now()

	m.UpdatedAt = time.Now()

	m.Color = GetColor()

	_, er := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

	if er != nil {
		return nil, er
	}

	res, err := a.interactionRepo.Create(m)

	if err != nil {
		return nil, err
	}

	y := res.(*models.Interactions)

	query1 := bson.M{"user_id": bson.ObjectIdHex(m.UserID.Hex()), "type": "collection", "collection_id": bson.ObjectIdHex(y.CollectionID.Hex()), "content_id": bson.M{"$ne": nil}}

	count, _ := a.interactionRepo.CountContent(query1)

	y.Count = count

	return y, nil
}

func (a *interactionUsecase) SaveToCollection(m *models.Interactions) (interface{}, error) {

	query := bson.M{"type": "collection", "collection_id": bson.ObjectIdHex(m.CollectionID.Hex()), "user_id": bson.ObjectIdHex(m.UserID.Hex()), "content_id": bson.ObjectIdHex(m.ContentID.Hex())}

	interact, _ := a.interactionRepo.FindWith(query)

	if interact != nil {
		return nil, errors.New("Content Has Already Been Added To Collection")
	}

	q := bson.M{"type": "collection", "collection_id": bson.ObjectIdHex(m.CollectionID.Hex()), "user_id": bson.ObjectIdHex(m.UserID.Hex())}

	existedCollection, _ := a.interactionRepo.FindWith(q)

	if existedCollection != nil {

		x := existedCollection.(*models.Interactions)

		fmt.Println("color", x)

		m.ID = bson.NewObjectId()

		m.CollectionID = x.CollectionID

		m.Name = x.Name

		m.CreatedAT = time.Now()

		m.UpdatedAt = time.Now()

		m.Color = x.Color

		_, er := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

		if er != nil {
			return nil, er
		}

		res, err := a.interactionRepo.Create(m)

		if err != nil {
			return nil, err
		}

		y := res.(*models.Interactions)

		query := bson.M{"user_id": bson.ObjectIdHex(m.UserID.Hex()), "type": "collection", "collection_id": bson.ObjectIdHex(y.CollectionID.Hex()), "content_id": bson.M{"$ne": nil}}

		count, _ := a.interactionRepo.CountContent(query)

		y.Count = count

		return y, nil
	}
	return nil, errors.New("Collection Does Not Exist")
}

func (a *interactionUsecase) LikeContent(m *models.Interactions) (interface{}, error) {

	query := bson.M{"type": "like", "user_id": bson.ObjectIdHex(m.UserID.Hex()), "content_id": bson.ObjectIdHex(m.ContentID.Hex())}

	interact, _ := a.interactionRepo.FindWith(query)

	if interact != nil {
		return nil, errors.New("You Have Already Liked This Content")
	}

	m.ID = bson.NewObjectId()

	m.CreatedAT = time.Now()

	m.UpdatedAt = time.Now()

	resContent, er := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

	if er != nil {
		return nil, er
	}

	var listContent models.Content

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(resContent)
	bson.Unmarshal(bsonBytes, &listContent)

	update := bson.M{"likes": listContent.Likes + 1, "updated_at": time.Now()}

	_, errr := a.interactionRepo.Update(update, m.ContentID.Hex(), "contents")

	if errr != nil {
		return nil, errr
	}

	_, err := a.interactionRepo.Create(m)

	if err != nil {
		return nil, err
	}

	// Index Document In Elastic
	go a.Index(m.ContentID.Hex())

	Content, _ := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

	var x models.Content

	// convert bson to struct
	bsonByte, _ := bson.Marshal(Content)
	bson.Unmarshal(bsonByte, &x)

	x.IsLiked = true

	isSave := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(m.UserID.Hex()), "content_id": bson.ObjectIdHex(m.ContentID.Hex())}

	save, _ := a.interactionRepo.FindWith(isSave)

	if save != nil {
		x.IsSaved = true
	}

	return x, nil
}

func (a *interactionUsecase) ShareContent(m *models.Interactions) (interface{}, error) {

	m.ID = bson.NewObjectId()

	m.CreatedAT = time.Now()

	m.UpdatedAt = time.Now()

	resContent, er := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

	if er != nil {
		return nil, er
	}

	// Get User Detail To Send Push Notification
	r, _ := a.interactionRepo.FindUser(m.Recipient.Hex())
	u, _ := a.interactionRepo.FindUser(m.UserID.Hex())

	//Send Push Notification
	notification.SendPushNotification(u.FirstName+" "+u.LastName, "shared a content with you.", r.FirebaseTokens, m.ContentID.Hex(), resContent, "content")

	//Construct Notify Json And Store In DB
	var send models.Notification

	send.ID = bson.NewObjectId()
	send.UserID = bson.ObjectIdHex(m.Recipient.Hex())
	send.SenderID = bson.ObjectIdHex(m.UserID.Hex())
	send.Type = "share"
	send.Topic = "Someone Shared A Content With You"
	send.Message = m.Comment
	send.ContentID = m.ContentID
	send.CreatedAT = time.Now()
	send.UpdatedAt = time.Now()

	notify, _ := a.interactionRepo.Notify(send)

	if notify != nil {
		fmt.Println(notify)
	}

	var listContent models.Content

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(resContent)
	bson.Unmarshal(bsonBytes, &listContent)

	update := bson.M{"shares": listContent.Shares + 1, "updated_at": time.Now()}

	_, errr := a.interactionRepo.Update(update, m.ContentID.Hex(), "contents")

	if errr != nil {
		return nil, errr
	}

	_, err := a.interactionRepo.Create(m)

	if err != nil {
		return nil, err
	}

	Content, _ := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

	return Content, nil
}

func (a *interactionUsecase) BookmarkContent(m *models.Interactions) (interface{}, error) {

	query := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(m.UserID.Hex()), "content_id": bson.ObjectIdHex(m.ContentID.Hex())}

	interact, _ := a.interactionRepo.FindWith(query)

	if interact != nil {
		return nil, errors.New("You Have Already Bookmarked This Content")
	}

	m.ID = bson.NewObjectId()

	m.CreatedAT = time.Now()

	m.UpdatedAt = time.Now()

	resContent, er := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

	if er != nil {
		return nil, er
	}

	var listContent models.Content

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(resContent)
	bson.Unmarshal(bsonBytes, &listContent)

	update := bson.M{"saves": listContent.Saves + 1, "updated_at": time.Now()}

	_, errr := a.interactionRepo.Update(update, m.ContentID.Hex(), "contents")

	if errr != nil {
		return nil, errr
	}

	_, err := a.interactionRepo.Create(m)

	if err != nil {
		return nil, err
	}

	Content, _ := a.interactionRepo.Find(m.ContentID.Hex(), "contents")

	// Index Document In Elastic
	go a.Index(m.ContentID.Hex())

	var x models.Content

	// convert bson to struct
	bsonByte, _ := bson.Marshal(Content)
	bson.Unmarshal(bsonByte, &x)

	x.IsSaved = true

	isLike := bson.M{"type": "like", "user_id": bson.ObjectIdHex(m.UserID.Hex()), "content_id": bson.ObjectIdHex(m.ContentID.Hex())}

	like, _ := a.interactionRepo.FindWith(isLike)

	if like != nil {
		x.IsLiked = true
	}

	return x, nil
}

func (a *interactionUsecase) Update(m interface{}, id string, collection string) (interface{}, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount != nil {
		return nil, models.ErrNotFound
	}

	res, err := a.interactionRepo.Update(m, id, collection)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *interactionUsecase) DeleteCollection(id string, user string) (string, error) {

	match := bson.M{"type": "collection", "collection_id": bson.ObjectIdHex(id), "user_id": bson.ObjectIdHex(user)}

	err := a.interactionRepo.DeleteMany(match)

	fmt.Println("here", err, id, user)
	if err != nil {
		return "", err
	}

	return "Collection Removed", nil
}

func (a *interactionUsecase) DeleteContentFromCollection(name string, content string, user string) (string, error) {

	match := bson.M{"type": "collection", "name": name, "content_id": bson.ObjectIdHex(content), "user_id": bson.ObjectIdHex(user)}

	err := a.interactionRepo.DeleteCollection(match)

	if err != nil {
		return "", err
	}

	return "Content Removed From Collection", nil
}

func (a *interactionUsecase) UpdateCollectionName(collectionID string, newName string, id string) (interface{}, error) {
	// Check If User Has Collection
	query := bson.M{"type": "collection", "user_id": bson.ObjectIdHex(id), "collection_id": bson.ObjectIdHex(collectionID)}

	existedAccount, _ := a.interactionRepo.FindWith(query)

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	update := bson.M{"name": newName}

	_, err := a.interactionRepo.UpdateCollectionName(update, query)

	if err != nil {
		return nil, err
	}
	//Get Updated Payload
	res, _ := a.interactionRepo.FindWith(query)

	if res == nil {
		return nil, models.ErrNotFound
	}

	y := res.(*models.Interactions)

	query1 := bson.M{"user_id": bson.ObjectIdHex(id), "type": "collection", "collection_id": bson.ObjectIdHex(y.CollectionID.Hex()), "content_id": bson.M{"$ne": nil}}

	count, _ := a.interactionRepo.CountContent(query1)

	y.Count = count

	return y, nil
}
