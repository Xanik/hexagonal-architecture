package course

import (
	models "study/models"
)

//Usecase interface represents Course usecases
type Usecase interface {
	// Course Usecases
	Create(interface{}) (interface{}, error)
	Find(id string) (interface{}, error)
	FindBy(key string, value string) (*models.Course, error)
	FindAll() ([]interface{}, error)
	Update(interface{}, string) (interface{}, error)
}
