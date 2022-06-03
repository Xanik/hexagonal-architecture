package usecase

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"study/config"
	"study/features/search"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type searchUsecase struct {
	searchRepo search.Repository
}

// NewAccountUsecase will create new an accountUsecase object representation of account.Usecase interface
func NewAccountUsecase(a search.Repository) search.Usecase {
	return &searchUsecase{
		searchRepo: a,
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (a *searchUsecase) SearchElastic(indexName string, key string, value string) (interface{}, error) {

	resSearch, err := a.searchRepo.SearchElastic(indexName, key, strings.ToLower(value))

	if err != nil {
		return nil, err
	}
	return resSearch, nil
}

func EmptyElastic(indexName string) {

	url := fmt.Sprintf("%s%s", config.Env.ElasticsearchURL, indexName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := http.Client{
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(data))
}

func (a *searchUsecase) SearchAccount(indexName string, value string, id string) (interface{}, error) {

	resSearch, err := a.searchRepo.SearchElastic(indexName, "email", strings.ToLower(value))

	if err != nil {
		res, er := a.searchRepo.SearchElastic(indexName, "first_name", strings.ToLower(value))
		if er != nil {
			resAccount, e := a.searchRepo.SearchElastic(indexName, "last_name", strings.ToLower(value))
			if e != nil {
				return nil, nil
			}
			for _, value := range resAccount {
				// Check If User Has Followed Specified Follower
				query := bson.M{"user_id": bson.ObjectIdHex(id), "following": bson.ObjectIdHex(value["id"].(string))}

				resp, _ := a.searchRepo.FindWith(query, "networks")

				if resp != nil {
					value["isfollowed"] = true
				} else {
					value["isfollowed"] = false
				}
				value["_id"] = value["id"]
			}
			fmt.Println(resAccount)
			return resAccount, nil
		}
		for _, value := range res {
			// Check If User Has Followed Specified Follower
			query := bson.M{"user_id": bson.ObjectIdHex(id), "following": bson.ObjectIdHex(value["id"].(string))}

			resp, _ := a.searchRepo.FindWith(query, "networks")

			if resp != nil {
				value["isfollowed"] = true
			} else {
				value["isfollowed"] = false
			}
			value["_id"] = value["id"]
		}

		return res, nil
	}

	for _, value := range resSearch {
		// Check If User Has Followed Specified Follower
		query := bson.M{"user_id": bson.ObjectIdHex(id), "following": bson.ObjectIdHex(value["id"].(string))}

		resp, _ := a.searchRepo.FindWith(query, "networks")

		if resp != nil {
			value["isfollowed"] = true
		} else {
			value["isfollowed"] = false
		}
		value["_id"] = value["id"]

	}

	return resSearch, nil
}

func (a *searchUsecase) SearchContent(indexName string, value string, id string) (interface{}, error) {

	resSearch, err := a.searchRepo.SearchElastic(indexName, "title", strings.ToLower(value))
	var data []map[string]interface{}

	if err != nil {
		res, e := a.searchRepo.SearchElastic(indexName, "tags", strings.ToLower(value))
		if e != nil {
			response, er := a.searchRepo.SearchElastic(indexName, "description", strings.ToLower(value))
			if er != nil {
				return nil, nil
			}
			for _, x := range response {

				// Check If User Has Liked Content And Like
				query1 := bson.M{"type": "like", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x["id"].(string))}

				resp1, _ := a.searchRepo.FindWith(query1, "interactions")

				if resp1 != nil {
					x["isliked"] = true
				} else {
					x["isliked"] = false
				}
				// Check If User Has Bookmarked Content And Save
				query := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x["id"].(string))}

				resp, _ := a.searchRepo.FindWith(query, "interactions")

				if resp != nil {
					x["issaved"] = true
				} else {
					x["issaved"] = false
				}
				x["_id"] = x["id"]
				data = append(data, x)
			}
			return data, nil
		}

		for _, x := range res {

			// Check If User Has Liked Content And Like
			query1 := bson.M{"type": "like", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x["id"].(string))}

			resp1, _ := a.searchRepo.FindWith(query1, "interactions")

			if resp1 != nil {
				x["isliked"] = true
			} else {
				x["isliked"] = false
			}
			// Check If User Has Bookmarked Content And Save
			query := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x["id"].(string))}

			resp, _ := a.searchRepo.FindWith(query, "interactions")

			if resp != nil {
				x["issaved"] = true
			} else {
				x["issaved"] = false
			}
			x["_id"] = x["id"]
			data = append(data, x)
		}
		return data, nil

	}

	for _, x := range resSearch {

		// Check If User Has Liked Content And Like
		query1 := bson.M{"type": "like", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x["id"].(string))}

		resp1, _ := a.searchRepo.FindWith(query1, "interactions")

		if resp1 != nil {
			x["isliked"] = true
		} else {
			x["isliked"] = false
		}
		// Check If User Has Bookmarked Content And Save
		query := bson.M{"type": "bookmark", "user_id": bson.ObjectIdHex(id), "content_id": bson.ObjectIdHex(x["id"].(string))}

		resp, _ := a.searchRepo.FindWith(query, "interactions")

		if resp != nil {
			x["issaved"] = true
		} else {
			x["issaved"] = false
		}
		x["_id"] = x["id"]
		data = append(data, x)
	}
	return data, nil
}

func (a *searchUsecase) SearchCourse(indexName string, value string, id string) (interface{}, error) {

	resSearch, err := a.searchRepo.SearchElastic(indexName, "title", strings.ToLower(value))

	if err != nil {
		return nil, nil
	}

	for _, x := range resSearch {
		x["_id"] = x["id"]
	}

	return resSearch, nil
}

func (a *searchUsecase) SearchInterest(indexName string, value string) (interface{}, error) {

	resSearch, err := a.searchRepo.SearchElastic(indexName, "name", strings.ToLower(value))

	if err != nil {
		return nil, nil
	}
	return resSearch, nil
}

func (a *searchUsecase) IndexAccounts(collection string) (string, error) {
	//Delete Index To Bulk ReIndex
	EmptyElastic(collection)

	resAccount, err := a.searchRepo.FindAllAccounts(collection)

	if err != nil {
		return "", err
	}

	res, er := a.searchRepo.CreateIndex(collection)

	if er != nil {
		return "", er
	}

	fmt.Println("INDEX SUCCESSFUL", res)

	for _, value := range resAccount {
		resSearch, err := a.searchRepo.IndexDocument(value, collection, value["id"].(bson.ObjectId).Hex())

		if err != nil {
			return "", err
		}
		fmt.Println("INDEX RESPONSE", resSearch)
	}

	return "DOCUMENTS INDEXED SUCCESSFULLY", nil
}

func (a *searchUsecase) IndexContents(collection string) (string, error) {
	//Delete Index To Bulk ReIndex
	EmptyElastic(collection)

	resAccount, err := a.searchRepo.FindAllContents(collection)

	if err != nil {
		return "", err
	}

	res, er := a.searchRepo.CreateIndex(collection)

	if er != nil {
		return "", er
	}

	fmt.Println("INDEX SUCCESSFUL", res)

	for _, value := range resAccount {
		resSearch, err := a.searchRepo.IndexDocument(value, collection, value.ID.Hex())

		if err != nil {
			return "", err
		}
		fmt.Println("INDEX RESPONSE", resSearch)
	}

	return "DOCUMENTS INDEXED SUCCESSFULLY", nil
}

func (a *searchUsecase) IndexCourses(collection string) (string, error) {
	//Delete Index To Bulk ReIndex
	EmptyElastic(collection)

	resAccount, err := a.searchRepo.FindAllCourses(collection)

	if err != nil {
		return "", err
	}

	res, er := a.searchRepo.CreateIndex(collection)

	if er != nil {
		return "", er
	}

	fmt.Println("INDEX SUCCESSFUL", res)

	for _, value := range resAccount {
		value["id"] = value["_id"]
		delete(value, "_id")
		resSearch, err := a.searchRepo.IndexDocument(value, collection, value["id"].(bson.ObjectId).Hex())

		if err != nil {
			return "", err
		}
		fmt.Println("INDEX RESPONSE", resSearch)
	}

	return "DOCUMENTS INDEXED SUCCESSFULLY", nil
}

func (a *searchUsecase) IndexInterests(collection string) (string, error) {

	//Delete Index To Bulk ReIndex
	EmptyElastic(collection)

	resAccount, err := a.searchRepo.FindAllInterests(collection)

	if err != nil {
		return "", err
	}

	res, er := a.searchRepo.CreateIndex(collection)

	if er != nil {
		return "", er
	}

	fmt.Println("INDEX SUCCESSFUL", res)

	for _, value := range resAccount {
		id := randInt(0000, 9999)

		resSearch, err := a.searchRepo.IndexDocument(value, collection, strconv.Itoa(id))

		if err != nil {
			return "", err
		}
		fmt.Println("INDEX RESPONSE", resSearch)
	}

	return "DOCUMENTS INDEXED SUCCESSFULLY", nil
}
