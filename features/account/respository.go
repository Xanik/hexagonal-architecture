package account

import (
	models "study/models"
)

//Repository interface represents account services
type Repository interface {
	Create(*models.Account) (*models.Account, error)
	Find(id string) (*models.Account, error)
	FindBy(key string, value string) (*models.Account, error)
	FindAll() ([]*models.Account, error)
	Update(user interface{}, id string) (interface{}, error)
	FindUserNotifications(string) ([]map[string]interface{}, error)
	IndexDocument(interface{}, string, string)
	GetIndexedAccount(id string) (*models.SearchAccount, error)
}
