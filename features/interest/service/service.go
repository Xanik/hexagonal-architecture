package usecase

import (
	"study/features/interest"
	"study/models"

	"gopkg.in/mgo.v2/bson"
)

type interestUsecase struct {
	interestRepo interest.Repository
}

// NewAccountUsecase will create new an accountUsecase object representation of account.Usecase interface
func NewAccountUsecase(a interest.Repository) interest.Usecase {
	return &interestUsecase{
		interestRepo: a,
	}
}

func (a *interestUsecase) Find(id string) (interface{}, error) {
	resAccount, err := a.interestRepo.Find(id, "interests")

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *interestUsecase) FindBy(key string, value string) (*models.Interest, error) {
	resAccount, err := a.interestRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *interestUsecase) FindAll() ([]*models.Interest, error) {
	resAccount, err := a.interestRepo.FindAll()

	if err != nil {
		return nil, nil
	}
	return resAccount, nil
}

func (a *interestUsecase) Create(m *models.Interest) (*models.Interest, error) {

	m.ID = bson.NewObjectId()

	res, err := a.interestRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *interestUsecase) Update(m *models.Interest, id string) (*models.Interest, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount != nil {
		return nil, models.ErrNotFound
	}

	res, err := a.interestRepo.Update(m, id)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *interestUsecase) GetUsersSuggestedInterest(id string) (interface{}, error) {

	m, err := a.interestRepo.Find(id, "accounts")

	var listAccount *models.Account

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(m)
	bson.Unmarshal(bsonBytes, &listAccount)

	res, err := a.interestRepo.FindIn(listAccount.InterestID, "interests")

	if err != nil {
		return nil, nil
	}

	return res, nil
}
