package organization

import (
	models "study/models"
)

//Usecase interface represents Organization usecases
type Usecase interface {
	// Organization Usecases
	Create(*models.Organization) (*models.Organization, error)
	Find(id string) (*models.Organization, error)
	FindBy(key string, value string) (*models.Organization, error)
	FindAll() ([]*models.Organization, error)
	Update(user *models.Organization, id string) (*models.Organization, error)
}
