package interest

import (
	models "study/models"
)

//Usecase interface represents Organization usecases
type Usecase interface {
	// Organization Usecases
	Create(*models.Interest) (*models.Interest, error)
	Find(id string) (interface{}, error)
	FindBy(key string, value string) (*models.Interest, error)
	FindAll() ([]*models.Interest, error)
	Update(user *models.Interest, id string) (*models.Interest, error)
	GetUsersSuggestedInterest(id string) (interface{}, error)
}
