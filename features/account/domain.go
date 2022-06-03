package account

import (
	models "study/models"
)

//Usecase interface represents account usecases
type Usecase interface {
	// Account Usecases
	CreateAndSendMail(*models.Account) (*models.Account, error)
	CreateAndSendText(*models.Account) (*models.Account, error)
	Find(id string) (*models.Account, error)
	FindBy(key string, value string) (*models.Account, error)
	FindAll() ([]*models.Account, error)
	Update(user interface{}, id string) (interface{}, error)
	HashPasswordAndUpdate(*models.Account, string) (interface{}, string, error)
	CheckHashAndUpdate(data *models.Login, password string) (interface{}, string, error)
	CompareCodeAndVerify(key string, value string, code int) (interface{}, string, error)
	SendNewVerificationCodeByMail(key string, value string) (interface{}, error)
	SendNewVerificationCodeByText(key string, value string) (interface{}, error)
	CheckHashAndUpdatePassword(string, string, string) (interface{}, string, error)
	SendFeedBack(*models.FeedBack) (string, error)
	FindUserNotifications(id string) ([]map[string]interface{}, error)
}
