package usecase

import (
	"study/features/institution"
	"study/models"
)

type institutionUsecase struct {
	institutionRepo institution.Repository
}

// NewInstitutionUsecase will create new an institutionUsecase object representation of institution.go.Usecase interface
func NewInstitutionUsecase(a institution.Repository) institution.Usecase {
	return &institutionUsecase{
		institutionRepo: a,
	}
}

func (a *institutionUsecase) Find(id string) (*models.Institution, error) {
	resAccount, err := a.institutionRepo.Find(id)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *institutionUsecase) FindBy(key string, value string) (*models.Institution, error) {
	resAccount, err := a.institutionRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *institutionUsecase) FindAll(skip int, limit int) ([]*models.Institution, error) {
	resAccount, err := a.institutionRepo.FindAll(skip, limit)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *institutionUsecase) Create(m *models.Institution) (*models.Institution, error) {

	existedAccount, _ := a.FindBy("name", m.Name)

	if existedAccount != nil {

		return existedAccount, nil
	}

	res, err := a.institutionRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *institutionUsecase) Update(m interface{}, id string) (interface{}, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	_, err := a.institutionRepo.Update(m, id)

	if err != nil {
		return nil, err
	}

	newAccount, _ := a.Find(id)

	return newAccount, nil
}
