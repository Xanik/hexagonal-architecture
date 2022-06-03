package interaction

import (
	"study/models"
)

//Usecase interface represents interaction usecases
type Usecase interface {
	// interaction Usecases
	Create(interface{}) (interface{}, error)
	CreateCollection(interactions *models.Interactions) (interface{}, error)
	SaveToCollection(interactions *models.Interactions) (interface{}, error)
	CreateEmptyCollection(interactions *models.Interactions) (interface{}, error)
	LikeContent(interactions *models.Interactions) (interface{}, error)
	UnLikeContent(string, string) (interface{}, error)
	ShareContent(interactions *models.Interactions) (interface{}, error)
	BookmarkContent(interactions *models.Interactions) (interface{}, error)
	DeleteBookmark(string, string) (interface{}, error)
	Find(id string) (interface{}, error)
	FindBy(key string, value string) (interface{}, error)
	FindByCollection(collection string, skip int, limit int, id string) (interface{}, error)
	FindAll(string, int, int) (interface{}, error)
	FindMyLikedContent(string, string, int, int) (interface{}, error)
	FindMySharedContent(string, string, int, int) (interface{}, error)
	FindMyBookmarkedContent(string, string, int, int) (interface{}, error)
	FindMyCollection(string, int, int) ([]interface{}, error)
	ListCollections(string, int, int) ([]interface{}, error)
	Update(interface{}, string, string) (interface{}, error)
	DeleteCollection(string, string) (string, error)
	DeleteContentFromCollection(string, string, string) (string, error)
	UpdateCollectionName(string, string, string) (interface{}, error)
}
