package usecase

import (
	"fmt"
	"math/rand"
	"study/features/account"
	"study/libs/emails"
	"study/libs/notification"
	"study/models"
	crypt "study/util/crypto"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type accountUsecase struct {
	accountRepo account.Repository
}

// NewAccountUsecase will create new an accountUsecase object representation of account.Usecase interface
func NewAccountUsecase(a account.Repository) account.Usecase {
	return &accountUsecase{
		accountRepo: a,
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (a *accountUsecase) Index(account string) {
	// Index Document In Elastic
	indexAccount, e := a.accountRepo.GetIndexedAccount(account)

	if indexAccount != nil {
		a.accountRepo.IndexDocument(indexAccount, "accounts", indexAccount.ID.Hex())
		//Index Completed
	}
	fmt.Println(e)

}

func (a *accountUsecase) Find(id string) (*models.Account, error) {
	resAccount, err := a.accountRepo.Find(id)

	if err != nil {
		return nil, err
	}

	resAccount.Password = ""

	return resAccount, nil
}

func (a *accountUsecase) FindBy(key string, value string) (*models.Account, error) {
	resAccount, err := a.accountRepo.FindBy(key, value)

	if err != nil {
		return nil, err
	}

	return resAccount, nil
}

func (a *accountUsecase) FindAll() ([]*models.Account, error) {
	resAccount, err := a.accountRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return resAccount, nil
}

func (a *accountUsecase) FindUserNotifications(id string) ([]map[string]interface{}, error) {
	resAccount, err := a.accountRepo.FindUserNotifications(id)

	if err != nil {
		return nil, err
	}

	return resAccount, nil
}

func (a *accountUsecase) CreateAndSendMail(m *models.Account) (*models.Account, error) {

	if !m.Valid() {
		return &models.Account{}, errors.New("invalid account type")
	}

	if !m.ValidFieldWithEmail() {
		return &models.Account{}, errors.New("No Empty Field Allowed")
	}

	if !m.ValidEmail() {
		return &models.Account{}, errors.New("Invalid Email")
	}

	existedAccount, _ := a.FindBy("email", m.Email)

	if existedAccount != nil {

		if existedAccount.Step == "SP" {
			return nil, errors.New("User Already Exist")
		}

		code := randInt(0000, 9999)

		update := bson.M{"code": code}

		a.Update(update, existedAccount.ID.Hex())

		go emails.SendTemplate("STUDY HALL", existedAccount.Email, existedAccount.FirstName, code)

		existedAccount.Code = code

		return existedAccount, nil
	}

	m.Code = randInt(1000, 9999)

	m.Validated = false

	m.CreatedWITH = "email"

	res, err := a.accountRepo.Create(m)

	if err != nil {
		return nil, err

	}

	go emails.SendTemplate("STUDY HALL", m.Email, m.FirstName, m.Code)

	return res, nil
}

func (a *accountUsecase) CreateAndSendText(m *models.Account) (*models.Account, error) {

	if !m.Valid() {
		return &models.Account{}, errors.New("invalid account type")
	}

	if !m.ValidFieldWithPhone() {
		return &models.Account{}, errors.New("No Empty Field Allowed")
	}

	if !m.ValidPhone() {
		return &models.Account{}, errors.New("Invalid Phone Number")
	}

	m.Phone = m.TransformPhone()

	existedAccount, _ := a.FindBy("phone", m.Phone)

	if existedAccount != nil {

		if existedAccount.Step == "SP" {
			return nil, errors.New("User Already Exist")
		}

		code := randInt(0000, 9999)

		update := bson.M{"code": code}

		a.Update(update, existedAccount.ID.Hex())

		go notification.SendTwilloNotification(existedAccount.Phone, code)

		existedAccount.Code = code

		return existedAccount, nil
	}

	m.Code = randInt(1000, 9999)

	m.Validated = false

	m.CreatedWITH = "phone"

	res, err := a.accountRepo.Create(m)

	if err != nil {
		return nil, err
	}

	go notification.SendTwilloNotification("+"+m.Phone, m.Code)

	return res, nil
}

func (a *accountUsecase) Update(m interface{}, id string) (interface{}, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount == nil {
		return nil, models.ErrNotFound
	}

	var Account models.Account

	// convert bson to struct
	bsonBytes, _ := bson.Marshal(m)
	bson.Unmarshal(bsonBytes, &Account)

	if Account.Email != "" {
		//Validate Email To Be Passed
		resAccount, _ := a.FindBy("email", Account.Email)

		if resAccount != nil {
			return nil, errors.New("Email Provided Already Exist")
		}
		if !Account.ValidEmail() {
			return &models.Account{}, errors.New("Invalid Email")
		}
	}

	if Account.Phone != "" {
		//Validate Phone To Be Passed
		res, _ := a.FindBy("phone", Account.Phone)

		if res != nil {
			return nil, errors.New("Number Provided Already Exist")
		}

		if !Account.ValidPhone() {
			return &models.Account{}, errors.New("Invalid Phone Number")
		}

		Account.Phone = Account.TransformPhone()
	}

	go a.Index(id)

	_, err := a.accountRepo.Update(m, id)

	if err != nil {
		return nil, err
	}

	newAccount, _ := a.Find(id)

	newAccount.Password = ""

	return newAccount, nil
}

func (a *accountUsecase) HashPasswordAndUpdate(m *models.Account, id string) (interface{}, string, error) {
	existedAccount, _ := a.Find(id)

	if existedAccount == nil {
		return nil, "", models.ErrNotFound
	}

	existedAccount.Password = crypt.HashText(m.Password)

	if !m.ValidPassword() {
		return &models.Account{}, "", errors.New("Weak Password")
	}

	token := crypt.Jwt(existedAccount.ID)

	update := bson.M{"password": existedAccount.Password, "step": "SP"}

	resAccount, err := a.Update(update, id)
	if err != nil {
		return nil, "", err
	}

	go a.Index(id)

	return resAccount, token, nil
}

func (a *accountUsecase) CheckHashAndUpdate(data *models.Login, password string) (interface{}, string, error) {

	var existPassword string

	var ID bson.ObjectId

	if data.Email == "" {
		if !data.ValidPhone() {
			return &models.Account{}, "", errors.New("Invalid Phone Number")
		}

		data.Phone = data.TransformPhone()

		resAccount, err := a.FindBy("phone", data.Phone)

		if err != nil {
			return nil, "", errors.New("Authentication Failed")
		}

		existPassword = resAccount.Password

		ID = resAccount.ID
	}

	if data.Phone == "" {
		resAccount, err := a.FindBy("email", data.Email)

		if err != nil {
			return nil, "", errors.New("Authentication Failed")
		}

		existPassword = resAccount.Password

		ID = resAccount.ID
	}

	hash := crypt.CheckTextHash(password, existPassword)

	if !hash {
		return nil, "", errors.New("Authentication Failed")
	}

	token := crypt.Jwt(ID)

	update := bson.M{"updated_at": time.Now(), "last_login": time.Now()}

	Account, err := a.Update(update, ID.Hex())

	if err != nil {
		return nil, "", err
	}

	return Account, token, nil
}

func (a *accountUsecase) CompareCodeAndVerify(key string, value string, code int) (interface{}, string, error) {
	var m models.Login

	if key == "phone" {
		m.Phone = value
		if !m.ValidPhone() {
			return &models.Account{}, "", errors.New("Invalid Phone Number")
		}

		m.Phone = m.TransformPhone()

		resAccount, _ := a.FindBy(key, m.Phone)

		if resAccount == nil {
			return nil, "", models.ErrNotFound
		}

		if code != resAccount.Code {
			return nil, "", errors.New("Invalid Code Provided")
		}

		token := crypt.Jwt(resAccount.ID)

		update := bson.M{"updated_at": time.Now(), "step": "VC", "validated": true}

		Account, err := a.Update(update, resAccount.ID.Hex())

		if err != nil {
			return nil, "", err
		}

		return Account, token, nil
	}
	resAccount, _ := a.FindBy(key, value)

	if resAccount == nil {
		return nil, "", models.ErrNotFound
	}

	if code != resAccount.Code {
		return nil, "", errors.New("Invalid Code Provided")
	}

	token := crypt.Jwt(resAccount.ID)

	update := bson.M{"updated_at": time.Now(), "step": "VC", "validated": true}

	Account, err := a.Update(update, resAccount.ID.Hex())

	if err != nil {
		return nil, "", err
	}

	return Account, token, nil
}

func (a *accountUsecase) SendNewVerificationCodeByMail(key string, value string) (interface{}, error) {
	resAccount, _ := a.FindBy(key, value)

	if resAccount == nil {
		return nil, models.ErrNotFound
	}

	resAccount.Code = randInt(1000, 9999)

	message := fmt.Sprintf("Your new verification code is %d", resAccount.Code)

	go notification.SendMail("test@gmail.com", "STUDY HALL", message, resAccount.Email)

	update := bson.M{"updated_at": time.Now(), "code": resAccount.Code}

	Account, err := a.Update(update, resAccount.ID.Hex())

	if err != nil {
		return nil, err
	}

	return Account, nil
}

func (a *accountUsecase) SendNewVerificationCodeByText(key string, value string) (interface{}, error) {
	resAccount, _ := a.FindBy(key, value)

	if resAccount == nil {
		return nil, models.ErrNotFound
	}

	resAccount.Code = randInt(1000, 9999)

	update := bson.M{"updated_at": time.Now(), "code": resAccount.Code}

	Account, err := a.Update(update, resAccount.ID.Hex())

	if err != nil {
		return nil, err
	}

	return Account, nil
}

func (a *accountUsecase) CheckHashAndUpdatePassword(id string, oldpassword string, newpassword string) (interface{}, string, error) {

	resAccount, err := a.accountRepo.Find(id)

	if err != nil {
		return nil, "", models.ErrNotFound
	}

	hash := crypt.CheckTextHash(oldpassword, resAccount.Password)

	if !hash {
		return nil, "", errors.New("Wrong Old Password")
	}

	resAccount.Password = newpassword

	if !resAccount.ValidPassword() {
		return &models.Account{}, "", errors.New("Weak Password")
	}

	newpassword = crypt.HashText(newpassword)

	token := crypt.Jwt(resAccount.ID)

	update := bson.M{"password": newpassword}

	Account, err := a.Update(update, id)

	if err != nil {
		return nil, "", err
	}

	return Account, token, nil
}

func (a *accountUsecase) SendFeedBack(m *models.FeedBack) (string, error) {

	existedAccount, _ := a.FindBy("email", m.Email)

	if existedAccount == nil {

		return "", errors.New("Mail Does Not Exist On The Platform")
	}

	go notification.SendMail(m.Email, "FEEDBACK", m.Description, "test@gmail.com")

	return "FeedBack Sent Successfully", nil
}
