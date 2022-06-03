package httplib

import (
	"fmt"
	"net/http"
	"study/config/responses"
	"study/features/interest"
	httplib "study/libs/http"
	mws "study/middlewares"
	"study/models"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase interest.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us interest.Usecase) *mux.Router {
	handler := &AccountHandler{
		AUsecase: us,
	}

	//BASE ROUTE
	route.HandleFunc("/v1", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("here")
		resp := responses.GeneralResponse{Success: true, Message: "Study  server running....", Data: "Study SERVER v1.0"}
		httplib.Response(res, resp)
	})

	mwsWithAuth := mws.AuthorizationSingle

	//interestRoute
	interestRoute := route.PathPrefix("/v1/interests").Subrouter()
	interestRoute.HandleFunc("", mwsWithAuth(handler.FindAll)).Methods("GET")
	// Use Authentication Token For Below routes
	interestRoute.HandleFunc("", mwsWithAuth(handler.Create)).Methods("POST")
	interestRoute.HandleFunc("/account/{accountID}", mwsWithAuth(handler.GetUsersSuggestedInterest)).Methods("GET")
	interestRoute.HandleFunc("/{Key}/{Value}", mwsWithAuth(handler.FindBy)).Methods("GET")
	interestRoute.HandleFunc("/{interestID}", mwsWithAuth(handler.Find)).Methods("GET")
	interestRoute.HandleFunc("/{interestID}", mwsWithAuth(handler.Update)).Methods("PUT")

	return route
}

// Create will store the account by given request body
func (a *AccountHandler) Create(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Interest

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Create(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error creating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// Find will fetch the account based on given params
func (a *AccountHandler) Find(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	interestID := c.Params("interestID")

	listAccount, err := a.AUsecase.Find(interestID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindBy will fetch the account based on given params
func (a *AccountHandler) FindBy(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	Key := c.Params("Key")
	Value := c.Params("Value")

	listAccount, err := a.AUsecase.FindBy(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	listAccount, err := a.AUsecase.FindAll()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]*models.Interest, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Update will change the account by given request body
func (a *AccountHandler) Update(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	interestID := c.Params("interestID")

	var data *models.Interest

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, interestID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// GetUsersSuggestedInterest will fetch the account based on given params
func (a *AccountHandler) GetUsersSuggestedInterest(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	accountID := c.Params("accountID")

	listAccount, err := a.AUsecase.GetUsersSuggestedInterest(accountID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]bson.M, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}
