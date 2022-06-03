package usecase

import (
	"study/features/upload"
	"study/models"
)

type uploadUsecase struct {
	uploadRepo upload.Repository
}

// NewUploadUsecase will create new an uploadUsecase object representation of Upload.go.Usecase interface
func NewUploadUsecase(a upload.Repository) upload.Usecase {
	return &uploadUsecase{
		uploadRepo: a,
	}
}

func (a *uploadUsecase) Find(id string) (*models.Upload, error) {
	resAccount, err := a.uploadRepo.Find(id)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *uploadUsecase) FindBy(key string, value string) (*models.Upload, error) {
	resAccount, err := a.uploadRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *uploadUsecase) FindAll() ([]*models.Upload, error) {
	resAccount, err := a.uploadRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *uploadUsecase) Create(m *models.Upload) (*models.Upload, error) {

	existedAccount, _ := a.Find(m.ID.Hex())

	if existedAccount != nil {

		return existedAccount, nil
	}

	res, err := a.uploadRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *uploadUsecase) Update(m interface{}, id string) (interface{}, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	_, err := a.uploadRepo.Update(m, id)

	if err != nil {
		return nil, err
	}

	newAccount, _ := a.Find(id)

	return newAccount, nil
}
