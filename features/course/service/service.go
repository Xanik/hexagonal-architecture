package usecase

import (
	"study/features/course"
	"study/models"
)

type courseUsecase struct {
	courseRepo course.Repository
}

// NewAccountUsecase will create new an accountUsecase object representation of course.Usecase interface
func NewAccountUsecase(a course.Repository) course.Usecase {
	return &courseUsecase{
		courseRepo: a,
	}
}

func (a *courseUsecase) Find(id string) (interface{}, error) {
	resAccount, err := a.courseRepo.Find(id)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *courseUsecase) FindBy(key string, value string) (*models.Course, error) {
	resAccount, err := a.courseRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *courseUsecase) FindAll() ([]interface{}, error) {
	resAccount, err := a.courseRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *courseUsecase) Create(m interface{}) (interface{}, error) {

	res, err := a.courseRepo.Create(m)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *courseUsecase) Update(m interface{}, id string) (interface{}, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount != nil {
		return nil, models.ErrNotFound
	}

	res, err := a.courseRepo.Update(m, id)

	if err != nil {
		return nil, err
	}
	return res, nil
}
