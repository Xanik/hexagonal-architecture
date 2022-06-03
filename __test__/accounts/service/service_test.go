package usecase

import (
	"errors"
	"study/__test__/accounts/mocks"
	_accountUsecase "study/features/account/service"
	"study/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

var (
	mockAccount = models.Account{
		FirstName:      "Sandra",
		LastName:       "Ike",
		Email:          "test@gmail.com",
		Phone:          "2349098617616",
		Code:           3209,
		Image:          "image.png",
		Step:           "SP",
		Validated:      false,
		AccountType:    "user",
		Institution:    "Babcock University",
		OrganizationID: bson.ObjectId("5c8a633e16733262ce8f3e5e"),
		CreatedAT:      time.Now(),
		UpdatedAt:      time.Now(),
		LastLogin:      time.Now(),
	}
)

func TestFind(t *testing.T) {
	mockAccountRepo := new(mocks.Respository)

	t.Run("success", func(t *testing.T) {
		tempMockAccount := mockAccount
		tempMockAccount.ID = bson.ObjectId("5cec00d5202969617bc791e1")
		mockAccountRepo.On("Find", mock.AnythingOfType("string")).Return(&mockAccount, nil).Once()
		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.Find(mockAccount.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockAccountRepo.AssertExpectations(t)

	})

	t.Run("success", func(t *testing.T) {
		tempMockAccount := mockAccount
		tempMockAccount.ID = bson.ObjectId("5cec00d5202969617bc791e1")
		mockAccountRepo.On("FindBy", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&mockAccount, nil).Once()
		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.FindBy("email", mockAccount.Email)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockAccountRepo.AssertExpectations(t)

	})

	t.Run("error-failed", func(t *testing.T) {
		mockAccountRepo.On("Find", mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.Find(mockAccount.ID.Hex())

		assert.Error(t, err)
		assert.Nil(t, a)

		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockAccountRepo.On("FindBy", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.FindBy("email", mockAccount.Email)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockAccountRepo.AssertExpectations(t)
	})

}

func TestCreate(t *testing.T) {
	mockAccountRepo := new(mocks.Respository)

	t.Run("success", func(t *testing.T) {
		tempMockAccount := mockAccount
		tempMockAccount.ID = bson.ObjectId("5cec00d5202969617bc791e1")
		mockAccountRepo.On("FindBy", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, nil).Once()
		mockAccountRepo.On("Create", mock.AnythingOfType("*models.Account")).Return(&mockAccount, nil).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.CreateAndSendMail(&mockAccount)

		assert.NoError(t, err)
		assert.NotNil(t, a)
		assert.Equal(t, mockAccount.Email, tempMockAccount.Email)

		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		tempMockAccount := mockAccount
		tempMockAccount.ID = bson.ObjectId("5cec00d5202969617bc791e1")
		mockAccountRepo.On("FindBy", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, nil).Once()
		mockAccountRepo.On("Create", mock.AnythingOfType("*models.Account")).Return(&mockAccount, nil).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.CreateAndSendText(&mockAccount)

		assert.NoError(t, err)
		assert.NotNil(t, a)
		assert.Equal(t, mockAccount.Email, tempMockAccount.Email)

		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockAccountRepo.On("FindBy", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()
		mockAccountRepo.On("Create", mock.AnythingOfType("*models.Account")).Return(nil, errors.New("Unexpected")).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.CreateAndSendMail(&mockAccount)

		assert.Error(t, err)
		assert.Nil(t, a)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockAccountRepo.On("FindBy", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()
		mockAccountRepo.On("Create", mock.AnythingOfType("*models.Account")).Return(nil, errors.New("Unexpected")).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.CreateAndSendText(&mockAccount)

		assert.Error(t, err)
		assert.Nil(t, a)
		mockAccountRepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockAccountRepo := new(mocks.Respository)

	t.Run("success", func(t *testing.T) {
		tempMockAccount := mockAccount
		tempMockAccount.ID = bson.ObjectId("5cec00d5202969617bc791e1")
		mockAccountRepo.On("Find", mock.AnythingOfType("string")).Return(&mockAccount, nil).Once()
		mockAccountRepo.On("Update", mock.AnythingOfType("*models.Account"), mock.AnythingOfType("string")).Return(&mockAccount, nil).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.Update(&mockAccount, tempMockAccount.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockAccountRepo.On("Find", mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()
		mockAccountRepo.On("Update", mock.AnythingOfType("*models.Account"), mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()

		u := _accountUsecase.NewAccountUsecase(mockAccountRepo)

		a, err := u.Update(&mockAccount, mockAccount.ID.Hex())

		assert.Error(t, err)
		assert.Nil(t, a)
		mockAccountRepo.AssertExpectations(t)
	})

}
