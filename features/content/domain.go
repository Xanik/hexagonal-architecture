package content

import (
	models "study/models"
)

//Usecase interface represents Content usecases
type Usecase interface {
	// Content Usecases
	Create(*models.Content) (*models.Content, error)
	Find(id string) (interface{}, error)
	FindBy(key string, value string) (*models.Content, error)
	FindAll(int, int) (*models.ArrayResponse, error)
	FindAllWithQuery(int, int, string) (*models.ArrayResponse, error)
	Update(user *models.Content, id string) (*models.Content, error)
	GetUsersSuggestedContent(string, int, int) (*models.ArrayResponse, error)
	GetContentsByTag(string, string, int, int) (interface{}, error)
	Crawler(models.Item)
	CrawlerRss(models.Item)
}
