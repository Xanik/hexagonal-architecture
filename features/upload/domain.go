package upload

import (
	models "study/models"
)

//Usecase interface represents Upload.go usecases
type Usecase interface {
	// Upload Usecases
	Create(*models.Upload) (*models.Upload, error)
	Find(id string) (*models.Upload, error)
	FindBy(key string, value string) (*models.Upload, error)
	FindAll() ([]*models.Upload, error)
	Update(user interface{}, id string) (interface{}, error)
}
