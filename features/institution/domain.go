package institution

import (
	models "study/models"
)

//Usecase interface represents institution.go usecases
type Usecase interface {
	// Institution Usecases
	Create(*models.Institution) (*models.Institution, error)
	Find(id string) (*models.Institution, error)
	FindBy(key string, value string) (*models.Institution, error)
	FindAll(int, int) ([]*models.Institution, error)
	Update(user interface{}, id string) (interface{}, error)
}
