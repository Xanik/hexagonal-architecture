package httplib

import (
	"net/http"
	"strings"
	"study/config/responses"
	"study/features/account"
	httplib "study/libs/http"
	mws "study/middlewares"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"

	// mws "study/middlewares"

	"study/models"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase account.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us account.Usecase) *mux.Router {
	handler := &AccountHandler{
		AUsecase: us,
	}
	//BASE ROUTE
	route.HandleFunc("/v1", func(res http.ResponseWriter, req *http.Request) {
		resp := responses.GeneralResponse{Success: true, Message: "Study  server running....", Data: "Study SERVER v1.0"}
		httplib.Response(res, resp)
	})

	//mwsWithAuth adds authorization token to endpoints
	mwsWithAuth := mws.AuthorizationSingle

	//authRoute
	route.HandleFunc("/v1/auth", handler.SignIn).Methods("POST")
	route.HandleFunc("/v1/auth/password/{accountID}", handler.ChangePassword).Methods("POST")
	route.HandleFunc("/v1/auth/{accountID}", handler.SetPassword).Methods("POST")
	route.HandleFunc("/v1/auth/{Key}/{Value}", handler.Verify).Methods("POST")

	//accountRoute
	route.HandleFunc("/v1/accounts/mail", handler.CreateAndSendMail).Methods("POST")
	route.HandleFunc("/v1/accounts/text", handler.CreateAndSendText).Methods("POST")
	route.HandleFunc("/v1/accounts/feedback", handler.SendFeedBack).Methods("POST")
	route.HandleFunc("/v1/accounts/mail/{Key}/{Value}", handler.SendNewVerificationCodeByMail).Methods("POST")
	route.HandleFunc("/v1/accounts/text/{Key}/{Value}", handler.SendNewVerificationCodeByText).Methods("POST")
	route.HandleFunc("/v1/accounts", mwsWithAuth(handler.FindAll)).Methods("GET")
	route.HandleFunc("/v1/accounts/notifications", mwsWithAuth(handler.FindUserNotifications)).Methods("GET")
	route.HandleFunc("/v1/accounts/{Key}/{Value}", handler.FindBy).Methods("GET")
	route.HandleFunc("/v1/accounts/{accountID}", mwsWithAuth(handler.Find)).Methods("GET")
	route.HandleFunc("/v1/accounts/{accountID}", mwsWithAuth(handler.Update)).Methods("PUT")
	route.HandleFunc("/v1/accounts/interest/{accountID}", mwsWithAuth(handler.UpdateUserInterest)).Methods("PUT")

	return route
}

// CreateAndSendMail will store the account by given request body and send a mail
func (a *AccountHandler) CreateAndSendMail(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Account

	c.BindJSON(&data)

	data.Email = strings.ToLower(data.Email)

	listAccount, err := a.AUsecase.CreateAndSendMail(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error creating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "account created"}
	httplib.Response(res, resp)
}

// CreateAndSendText will store the account by given request body and send a text
func (a *AccountHandler) CreateAndSendText(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Account

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.CreateAndSendText(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "Error creating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "account created successfully"}
	httplib.Response(res, resp)
}

// Find will fetch the account based on given params
func (a *AccountHandler) Find(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	accountID := c.Params("accountID")

	listAccount, err := a.AUsecase.Find(accountID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching account"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "successfully fetched account details"}
	httplib.Response(res, resp)
}

// FindUserNotifications will fetch the account based on given params
func (a *AccountHandler) FindUserNotifications(res http.ResponseWriter, req *http.Request) {
	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.FindUserNotifications(accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching account"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "successfully fetched account details"}
	httplib.Response(res, resp)
}

// FindBy will fetch the account based on given params
func (a *AccountHandler) FindBy(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	Key := c.Params("Key")
	Value := c.Params("Value")

	listAccount, err := a.AUsecase.FindBy(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching account"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "successfully fetched account details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	listAccount, err := a.AUsecase.FindAll()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching account"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]*models.Account, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "successfully fetched account details"}
	httplib.Response(res, resp)
}

// Update will change the account by given request body
func (a *AccountHandler) Update(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	var data interface{}

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, accountID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "Error creating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "successfully updated account detail"}
	httplib.Response(res, resp)
}

// UpdateUserInterest will change the account by given request body
func (a *AccountHandler) UpdateUserInterest(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	var data map[string][]bson.ObjectId

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, accountID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "Error creating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "successfully updated account detail"}
	httplib.Response(res, resp)
}

// SetPassword will change the account password by given request body
func (a *AccountHandler) SetPassword(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	var data *models.Account

	c.BindJSON(&data)

	listAccount, token, err := a.AUsecase.HashPasswordAndUpdate(data, accountID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error updating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	datum := struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}{User: listAccount, Token: token}

	resp := responses.GeneralResponse{Success: true, Data: datum, Message: "password set successfully"}
	httplib.Response(res, resp)
}

// SignIn authenticate account by given request body
func (a *AccountHandler) SignIn(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Login

	c.BindJSON(&data)

	data.Email = strings.ToLower(data.Email)

	listAccount, token, err := a.AUsecase.CheckHashAndUpdate(data, data.Password)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "Error authenticating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	datum := struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}{User: listAccount, Token: token}

	resp := responses.GeneralResponse{Success: true, Data: datum, Message: "Authentication successful"}
	httplib.Response(res, resp)
}

// Verify authenticate account by given request body
func (a *AccountHandler) Verify(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	Key := c.Params("Key")
	Value := c.Params("Value")

	var data *models.Account

	c.BindJSON(&data)

	listAccount, token, err := a.AUsecase.CompareCodeAndVerify(Key, Value, data.Code)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error updating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}
	datum := struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}{User: listAccount, Token: token}

	resp := responses.GeneralResponse{Success: true, Data: datum, Message: "Verification successful"}
	httplib.Response(res, resp)
}

// SendNewVerificationCodeByMail authenticate account by given request body
func (a *AccountHandler) SendNewVerificationCodeByMail(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	Key := c.Params("Key")
	Value := c.Params("Value")

	var data *models.Account

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.SendNewVerificationCodeByMail(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error updating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "Email Sent"}
	httplib.Response(res, resp)
}

// SendNewVerificationCodeByText authenticate account by given request body
func (a *AccountHandler) SendNewVerificationCodeByText(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	Key := c.Params("Key")
	Value := c.Params("Value")

	var data *models.Account

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.SendNewVerificationCodeByText(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error updating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "Text Sent"}
	httplib.Response(res, resp)
}

// ChangePassword will change the account password by given request body
func (a *AccountHandler) ChangePassword(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	var data *models.Password

	c.BindJSON(&data)

	listAccount, token, err := a.AUsecase.CheckHashAndUpdatePassword(accountID, data.OldPassword, data.NewPassword)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error updating account", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	datum := struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}{User: listAccount, Token: token}

	resp := responses.GeneralResponse{Success: true, Data: datum, Message: "password updated successfully"}
	httplib.Response(res, resp)
}

// SendFeedBack will send a mail
func (a *AccountHandler) SendFeedBack(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.FeedBack

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.SendFeedBack(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error sending feedback", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "Mail Sent"}
	httplib.Response(res, resp)
}
