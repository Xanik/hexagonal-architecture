package httplib

import (
	"fmt"
	"net/http"
	"strconv"
	"study/config/responses"
	"study/features/network"
	httplib "study/libs/http"
	mws "study/middlewares"

	"gopkg.in/mgo.v2/bson"

	// mws "study/middlewares"

	"study/models"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase network.Usecase
}

// NewAccountHandler will initialize the network endpoint
func NewAccountHandler(route *mux.Router, us network.Usecase) *mux.Router {
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

	//networkRoute
	//route.Use(mwsWithAuth)
	route.HandleFunc("/v1/networks", mwsWithAuth(handler.Create)).Methods("POST")
	route.HandleFunc("/v1/networks", mwsWithAuth(handler.FindAll)).Methods("GET")
	route.HandleFunc("/v1/networks/followers/{networkID}", mwsWithAuth(handler.FindFollowers)).Methods("GET")
	route.HandleFunc("/v1/networks/following/{networkID}", mwsWithAuth(handler.FindFollowing)).Methods("GET")
	route.HandleFunc("/v1/networks/suggested/{networkID}", mwsWithAuth(handler.GetUsersSuggestedFollowers)).Methods("GET")
	route.HandleFunc("/v1/networks/{Key}/{Value}", mwsWithAuth(handler.FindBy)).Methods("GET")
	route.HandleFunc("/v1/networks/{networkID}", mwsWithAuth(handler.Find)).Methods("GET")
	route.HandleFunc("/v1/networks/followers/{networkID}", mwsWithAuth(handler.UpdateUserFollowers)).Methods("PUT")
	route.HandleFunc("/v1/networks/following/{networkID}", mwsWithAuth(handler.UpdateUserFollowing)).Methods("PUT")
	route.HandleFunc("/v1/networks/follow/{accountID}/{followID}", mwsWithAuth(handler.FollowAUser)).Methods("PUT")
	route.HandleFunc("/v1/networks/unfollow/{accountID}/{followID}", mwsWithAuth(handler.UnFollowAUser)).Methods("PUT")
	route.HandleFunc("/v1/networks/{networkID}", mwsWithAuth(handler.Update)).Methods("PUT")

	return route
}

// CreateAndSendMail will store the account by given request body and send a mail
func (a *AccountHandler) Create(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Network

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Create(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error creating network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network created"}
	httplib.Response(res, resp)
}

// Find will fetch the account based on given params
func (a *AccountHandler) Find(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	networkID := c.Params("networkID")

	listAccount, err := a.AUsecase.Find(networkID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// FindBy will fetch the account based on given params
func (a *AccountHandler) FindBy(res http.ResponseWriter, req *http.Request) {

	c := httplib.C{Res: res, Req: req}
	Key := c.Params("Key")
	Value := c.Params("Value")

	listAccount, err := a.AUsecase.FindBy(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	listAccount, err := a.AUsecase.FindAll()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching network"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]*models.Network, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// Update will change the account by given request body
func (a *AccountHandler) Update(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	networkID := c.Params("networkID")

	var data interface{}

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, networkID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// UpdateUserFollowers will change the account by given request body
func (a *AccountHandler) UpdateUserFollowers(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	networkID := c.Params("networkID")

	var data map[string][]bson.ObjectId

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, networkID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// UpdateUserFollowing will change the account by given request body
func (a *AccountHandler) UpdateUserFollowing(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	networkID := c.Params("networkID")

	var data map[string][]bson.ObjectId

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, networkID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// FollowAUser will change the account by given request body
func (a *AccountHandler) FollowAUser(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	followID := c.Params("followID")

	listAccount, err := a.AUsecase.FollowAUser(accountID, followID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// UnFollowAUser will change the account by given request body
func (a *AccountHandler) UnFollowAUser(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	followID := c.Params("followID")

	listAccount, err := a.AUsecase.UnFollowAUser(accountID, followID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// FindFollowers will fetch the account based on given params
func (a *AccountHandler) FindFollowers(res http.ResponseWriter, req *http.Request) {

	c := httplib.C{Res: res, Req: req}

	networkID := c.Params("networkID")

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindFollowers(networkID, x, y)

	fmt.Println(x, y)
	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// FindFollowing will fetch the account based on given params
func (a *AccountHandler) FindFollowing(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	networkID := c.Params("networkID")

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	fmt.Println(x, y)

	listAccount, err := a.AUsecase.FindFollowing(networkID, x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}

// GetUsersSuggestedFollowers will fetch the account based on given params
func (a *AccountHandler) GetUsersSuggestedFollowers(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	networkID := c.Params("networkID")

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.GetUsersSuggestedFollowers(networkID, x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching network"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
	httplib.Response(res, resp)
}
