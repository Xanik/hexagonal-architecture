package usecase

import (
	organization "study/features/organization"
	models "study/models"

	"gopkg.in/mgo.v2/bson"
)

type organizationUsecase struct {
	organizationRepo organization.Repository
}

// NewAccountUsecase will create new an accountUsecase object representation of account.Usecase interface
func NewAccountUsecase(a organization.Repository) organization.Usecase {
	return &organizationUsecase{
		organizationRepo: a,
	}
}

func (a *organizationUsecase) Find(id string) (*models.Organization, error) {
	resAccount, err := a.organizationRepo.Find(id)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *organizationUsecase) FindBy(key string, value string) (*models.Organization, error) {
	resAccount, err := a.organizationRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *organizationUsecase) FindAll() ([]*models.Organization, error) {
	resAccount, err := a.organizationRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *organizationUsecase) Create(m *models.Organization) (*models.Organization, error) {

	m.ID = bson.NewObjectId()

	res, err := a.organizationRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *organizationUsecase) Update(m *models.Organization, id string) (*models.Organization, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount != nil {
		return nil, models.ErrNotFound
	}

	res, err := a.organizationRepo.Update(m, id)

	if err != nil {
		return nil, err
	}
	return res, nil
}
