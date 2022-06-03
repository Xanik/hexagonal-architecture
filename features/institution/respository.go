package institution

import (
	models "study/models"
)

//Repository interface represents interest services
type Repository interface {
	Create(*models.Institution) (*models.Institution, error)
	Find(id string) (*models.Institution, error)
	FindBy(key string, value string) (*models.Institution, error)
	FindAll(int, int) ([]*models.Institution, error)
	Update(user interface{}, id string) (interface{}, error)
}
