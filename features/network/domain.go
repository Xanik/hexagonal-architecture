package network

import (
	models "study/models"
)

//Usecase interface represents network.go usecases
type Usecase interface {
	// Network Usecases
	Create(*models.Network) (interface{}, error)
	Find(string) (interface{}, error)
	FindBy(key string, value string) (*models.Network, error)
	FindAll() ([]*models.Network, error)
	Update(user interface{}, id string) (interface{}, error)
	FindFollowers(string, int, int) (interface{}, error)
	FindFollowing(string, int, int) (interface{}, error)
	GetUsersSuggestedFollowers(string, int, int) (interface{}, error)
	FollowAUser(string, string) (interface{}, error)
	UnFollowAUser(string, string) (interface{}, error)
}
