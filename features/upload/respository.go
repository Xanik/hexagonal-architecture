package upload

import (
	models "study/models"
)

//Repository interface represents network services
type Repository interface {
	Create(*models.Upload) (*models.Upload, error)
	Find(id string) (*models.Upload, error)
	FindBy(key string, value string) (*models.Upload, error)
	FindAll() ([]*models.Upload, error)
	Update(user interface{}, id string) (interface{}, error)
}
