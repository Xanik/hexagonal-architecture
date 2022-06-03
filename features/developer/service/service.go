package usecase

import (
	"study/features/developer"
	"study/models"

	"gopkg.in/mgo.v2/bson"

	"crypto/rand"
	"encoding/base32"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Usecase struct {
	Repo developer.Repository
}

// NewAccountUsecase will create new an accountUsecase object representation of account.Usecase interface
func NewAccountUsecase(a developer.Repository) developer.Usecase {
	return &Usecase{
		Repo: a,
	}
}

func genHex() string {
	buf := make([]byte, 20)
	rand.Read(buf)
	str := base32.StdEncoding.EncodeToString(buf)
	return str
}

func (a *Usecase) RequestAccess(developer *models.Developer) (*models.Developer, error) {
	developer.ID = bson.NewObjectId()
	// Use random string with timestamp for Refresh Token
	developer.RefreshToken = genHex()

	// To Expire every 48hrs
	developer.ExpiresIn = 48

	res, err := a.Repo.RequestAccess(developer)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *Usecase) AuthenticateAccount(developer *models.Developer) (string, error) {
	res, err := a.Repo.GetDeveloperAccount(developer.RefreshToken)

	if err != nil {
		return "", models.ErrNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": res,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		return "", error
	}
	return tokenString, error
}

func (a *Usecase) GetAllContents() ([]*models.Content, error) {

	res, err := a.Repo.GetAllContents()

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *Usecase) GetContentByInterest(request *models.DeveloperRequest) ([]*models.Content, error) {

	res, err := a.Repo.GetContentByInterest(request.Interests)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *Usecase) SearchContent(request string) ([]map[string]interface{}, error) {

	// Search Indexed Documents
	resSearch, err := a.Repo.SearchContent("contents", "title", strings.ToLower(request))
	var data []map[string]interface{}

	if err != nil {
		res, e := a.Repo.SearchContent("contents", "tags", strings.ToLower(request))
		if e != nil {
			response, er := a.Repo.SearchContent("contents", "description", strings.ToLower(request))
			if er != nil {
				return nil, nil
			}
			for _, x := range response {
				x["_id"] = x["id"]
				data = append(data, x)
			}
			return data, nil
		}

		for _, x := range res {
			x["_id"] = x["id"]
			data = append(data, x)
		}
		return data, nil
	}
	for _, x := range resSearch {
		x["_id"] = x["id"]
		data = append(data, x)
	}
	return data, nil
}

func (a *Usecase) GetAllCourses() ([]*models.Course, error) {

	res, err := a.Repo.GetAllCourses()

	if err != nil {
		return nil, err
	}

	return res, nil
}
