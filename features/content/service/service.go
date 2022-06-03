package usecase

import (
	"fmt"
	"math/rand"
	"study/features/content"
	"study/models"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type contentUsecase struct {
	contentRepo content.Repository
}

// NewAccountUsecase will create new an contentUsecase object representation of content.Usecase interface
func NewAccountUsecase(a content.Repository) content.Usecase {
	return &contentUsecase{
		contentRepo: a,
	}
}

func GetTag() []string {
	Tags := []string{
		"Technology", "Photography", "Art", "Business", "Fashion", "Finance", "Health",
		"Literature", "Software", "Study Hacks", "Career", "Economy", "Education", "Employment",
		"Entertainment", "Family", "Food", "Games", "Investments", "Journalism", "Law",
		"Liquidity", "Military", "Money", "Music", "Physics", "Plants", "Schools",
		"Science", "Space", "Sport", "Transportation", "Tourism", "War", "Wedding", "Weather",
	}

	// randomly pick one beautiful girl(URL) from the URLs slice
	rand.Seed(time.Now().UnixNano())

	chosen := make([]string, 0)

	for i := 0; i < 5; i++ {
		chosen = append(chosen, Tags[rand.Intn(len(Tags)-1)])
	}

	fmt.Println("Randomly selected this Tag : ", chosen)
	return chosen
}

func Shuffle(vals []*models.Content) []*models.Content {
	content := make([]*models.Content, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(vals)) {
		val := vals[i]
		content = append(content, val)
		fmt.Println(val)
	}
	return content
}

func (a *contentUsecase) Find(id string) (interface{}, error) {
	resAccount, err := a.contentRepo.Find(id, "contents")

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *contentUsecase) FindBy(key string, value string) (*models.Content, error) {
	resAccount, err := a.contentRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *contentUsecase) FindAll(skip int, limit int) (*models.ArrayResponse, error) {
	resAccount, err := a.contentRepo.FindAll(skip, limit)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *contentUsecase) FindAllWithQuery(skip int, limit int, id string) (*models.ArrayResponse, error) {
	match := bson.M{"likes": bson.M{"$gte": 0},
		"created_at": bson.M{
			"$gt": time.Now().Add(-168 * time.Hour)}}

	resAccount, err := a.contentRepo.FindAllWithQuery(match, skip, limit)

	for key, x := range resAccount.Contents {
		// Check If User Has Liked Content And Save
		query1 := bson.M{"type": "like", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x.ID.Hex())}
		res, _ := a.contentRepo.FindWith(query1, "interactions")

		if res != nil {
			x.IsLiked = true
		} else {
			x.IsLiked = false
		}

		// Check If User Has Bookmarked Content And Save
		query2 := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x.ID.Hex())}
		resp, _ := a.contentRepo.FindWith(query2, "interactions")

		if resp != nil {
			x.IsSaved = true
		} else {
			x.IsSaved = false
		}
		resAccount.Contents[key] = x
	}

	if err != nil {
		return nil, err
	}

	resAccount.Contents = Shuffle(resAccount.Contents)

	return resAccount, nil
}

func (a *contentUsecase) Create(m *models.Content) (*models.Content, error) {

	m.ID = bson.NewObjectId()

	m.CreatedAT = time.Now()

	m.UpdatedAt = time.Now()

	// for index, value := range m.Media {
	// 	value.ID = bson.NewObjectId()
	// 	value.CreatedAT = time.Now()
	// 	value.UpdatedAt = time.Now()

	// 	m.Media[index] = value

	// }

	// fmt.Println(m.Media[0])

	res, err := a.contentRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *contentUsecase) Update(m *models.Content, id string) (*models.Content, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount != nil {
		return nil, models.ErrNotFound
	}

	res, err := a.contentRepo.Update(m, id)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *contentUsecase) GetUsersSuggestedContent(id string, skip int, limit int) (*models.ArrayResponse, error) {

	m, er := a.contentRepo.Find(id, "accounts")

	if er != nil {
		return nil, er
	}

	var listAccount *models.Account

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(m)
	bson.Unmarshal(bsonBytes, &listAccount)

	res, err := a.contentRepo.FindInInterest(listAccount.InterestID, "interests")

	if err != nil {
		return nil, err
	}
	x := res.([]models.Interest)

	var InterestNames []string

	for _, value := range x {
		InterestNames = append(InterestNames, value.Name)
	}

	resContent, errr := a.contentRepo.FindInContent(InterestNames, "contents", skip, limit)

	if errr != nil {
		return nil, errr
	}

	for key, x := range resContent.Contents {
		// Check If User Has Liked Content And Save
		query1 := bson.M{"type": "like", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x.ID.Hex())}
		res, _ := a.contentRepo.FindWith(query1, "interactions")

		fmt.Println("content", res)
		if res != nil {
			x.IsLiked = true
		}

		// Check If User Has Bookmarked Content And Save
		query2 := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x.ID.Hex())}
		resp, _ := a.contentRepo.FindWith(query2, "interactions")

		if resp != nil {
			x.IsSaved = true
		}
		resContent.Contents[key] = x
	}

	return resContent, nil
}

func (a *contentUsecase) GetContentsByTag(id string, name string, skip int, limit int) (interface{}, error) {

	resContent, errr := a.contentRepo.FindInContent([]string{name}, "contents", skip, limit)

	if errr != nil {
		return nil, errr
	}

	for key, x := range resContent.Contents {
		// Check If User Has Liked Content And Save
		query1 := bson.M{"type": "like", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x.ID.Hex())}
		res, _ := a.contentRepo.FindWith(query1, "interactions")

		fmt.Println("content", res)
		if res != nil {
			x.IsLiked = true
		}

		// Check If User Has Bookmarked Content And Save
		query2 := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x.ID.Hex())}
		resp, _ := a.contentRepo.FindWith(query2, "interactions")

		if resp != nil {
			x.IsSaved = true
		}
		resContent.Contents[key] = x
	}

	return resContent, nil
}

func (a *contentUsecase) CrawlerRss(m models.Item) {
	existedAccount, _ := a.FindBy("title", m.Title)

	if existedAccount != nil {
		fmt.Println("Duplicate Feed......")
		return
	}

	//Map Opportunity Data
	var data models.Content

	var listTags []string

	listTags = append(listTags, m.Category)

	data.ID = bson.NewObjectId()
	data.Title = m.Title
	data.Description = m.Desc
	data.Summary = m.Desc
	data.Library = m.Link
	data.Authors = []string{}
	data.UpdatedAt = time.Now()
	data.CreatedAT = time.Now()
	data.Image = ""
	data.Platform = m.Company
	data.Tags = GetTag()
	data.Likes = 0
	data.Saves = 0
	data.Shares = 0

	_, err := a.contentRepo.Create(&data)

	if err != nil {
		fmt.Println(err, "Unable To Create Content......")
	}
	fmt.Println("Creating Content......")

}

func (a *contentUsecase) Crawler(m models.Item) {
	existedAccount, _ := a.FindBy("title", m.Title)

	if existedAccount != nil {
		fmt.Println("Duplicate Feed......")
		return
	}

	//Map Opportunity Data
	var data models.Content

	var listTags []string

	listTags = append(listTags, m.Category)

	data.ID = bson.NewObjectId()
	data.Title = m.Title
	data.Description = m.Desc
	data.Summary = m.Desc
	data.Library = m.Link
	data.Authors = []string{}
	data.UpdatedAt = time.Now()
	data.CreatedAT = time.Now()
	data.Image = ""
	data.Platform = m.Company
	data.Tags = GetTag()
	data.Likes = 0
	data.Saves = 0
	data.Shares = 0

	_, err := a.contentRepo.Create(&data)

	if err != nil {
		fmt.Println(err, "Unable To Create Content......")
	}
	fmt.Println("Creating Content......")

}
