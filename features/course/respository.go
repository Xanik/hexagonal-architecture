package course

import (
	models "study/models"
)

//Repository interface represents account services
type Repository interface {
	Create(interface{}) (interface{}, error)
	Find(id string) (interface{}, error)
	FindBy(key string, value string) (*models.Course, error)
	FindAll() ([]interface{}, error)
	Update(interface{}, string) (interface{}, error)
}
