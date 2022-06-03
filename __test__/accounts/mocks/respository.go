package mocks

import (
	"study/models"

	"github.com/stretchr/testify/mock"
)

// Respository is an autogenerated __test__ type for the Usecase type
type Respository struct {
	mock.Mock
}

func (_m *Respository) FindUserNotifications(string) ([]map[string]interface{}, error) {
	panic("implement me")
}

func (_m *Respository) IndexDocument(interface{}, string, string) {
	panic("implement me")
}

func (_m *Respository) GetIndexedAccount(id string) (*models.SearchAccount, error) {
	panic("implement me")
}

// Create provides a mock function with given fields: a
func (_m *Respository) Create(a *models.Account) (*models.Account, error) {

	ret := _m.Called(a)

	var r0 *models.Account

	if rf, ok := ret.Get(0).(func(*models.Account) *models.Account); ok {
		r0 = rf(a)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error

	if rf, ok := ret.Get(1).(func(*models.Account) error); ok {
		r1 = rf(a)
	} else {
		r1 = ret.Error((1))
	}

	return r0, r1

}

// Find provides a __test__ function with given fields: id
func (_m *Respository) Find(id string) (*models.Account, error) {
	ret := _m.Called(id)

	var r0 *models.Account

	if rf, ok := ret.Get(0).(func(string) *models.Account); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error((1))
	}

	return r0, r1
}

func (_m *Respository) FindBy(key, value string) (*models.Account, error) {
	ret := _m.Called(key, value)

	var r0 *models.Account

	if rf, ok := ret.Get(0).(func(string, string) *models.Account); ok {
		r0 = rf(key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(key, value)
	} else {
		r1 = ret.Error((1))
	}

	return r0, r1
}

func (_m *Respository) FindAll() ([]*models.Account, error) {
	ret := _m.Called()

	var r0 []*models.Account

	if rf, ok := ret.Get(0).(func() []*models.Account); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Account)
		}
	}

	var r1 error

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error((1))
	}

	return r0, r1
}

func (_m *Respository) Update(a interface{}, id string) (interface{}, error) {
	ret := _m.Called(a, id)

	var r0 interface{}

	if rf, ok := ret.Get(0).(func(interface{}, string) interface{}); ok {
		r0 = rf(a, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error

	if rf, ok := ret.Get(1).(func(interface{}, string) error); ok {
		r1 = rf(a, id)
	} else {
		r1 = ret.Error((1))
	}

	return r0, r1
}