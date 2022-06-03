package developer

import (
	models "study/models"
)

//Repository interface represents account services
type Repository interface {
	RequestAccess(*models.Developer) (*models.Developer, error)
	GetDeveloperAccount(string) (*models.Developer, error)
	SearchContent(string, string, string) ([]map[string]interface{}, error)
	GetAllContents() ([]*models.Content, error)
	GetContentByInterest([]string) ([]*models.Content, error)
	GetAllCourses() ([]*models.Course, error)
}
