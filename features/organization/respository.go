package organization

import (
	models "study/models"
)

//Repository interface represents account services
type Repository interface {
	Create(*models.Organization) (*models.Organization, error)
	Find(id string) (*models.Organization, error)
	FindBy(key string, value string) (*models.Organization, error)
	FindAll() ([]*models.Organization, error)
	Update(*models.Organization, string) (*models.Organization, error)
}
