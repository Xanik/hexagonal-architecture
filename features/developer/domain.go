package developer

import (
	models "study/models"
)

//Usecase interface represents Developer usecases
type Usecase interface {
	// Developer Usecases
	AuthenticateAccount(*models.Developer) (string, error)
	RequestAccess(*models.Developer) (*models.Developer, error)
	SearchContent(string) ([]map[string]interface{}, error)
	GetAllContents() ([]*models.Content, error)
	GetContentByInterest(*models.DeveloperRequest) ([]*models.Content, error)
	GetAllCourses() ([]*models.Course, error)
}
