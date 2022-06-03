package httplib

import (
	"fmt"
	"net/http"
	"study/config/responses"
	"study/features/organization"
	httplib "study/libs/http"
	mws "study/middlewares"

	// mws "study/middlewares"

	"study/models"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase organization.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us organization.Usecase) *mux.Router {
	handler := &AccountHandler{
		AUsecase: us,
	}

	//BASE ROUTE
	route.HandleFunc("/v1", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("here")
		resp := responses.GeneralResponse{Success: true, Message: "Study  server running....", Data: "Study SERVER v1.0"}
		httplib.Response(res, resp)
	})

	//mwsWithAuth adds authorization token to endpoints
	mwsWithAuth := mws.AuthorizationSingle

	//organizationRoute
	organizationRoute := route.PathPrefix("/v1/organizations").Subrouter()
	// Use Authentication Token For Below routes
	organizationRoute.HandleFunc("", mwsWithAuth(handler.Create)).Methods("POST")
	organizationRoute.HandleFunc("", mwsWithAuth(handler.FindAll)).Methods("GET")
	organizationRoute.HandleFunc("/{organizationKey}/{organizationValue}", mwsWithAuth(handler.FindBy)).Methods("GET")
	organizationRoute.HandleFunc("/{organizationID}", mwsWithAuth(handler.Find)).Methods("GET")
	organizationRoute.HandleFunc("/{organizationID}", mwsWithAuth(handler.Update)).Methods("PUT")

	return route
}

// Create will store the account by given request body
func (a *AccountHandler) Create(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Organization

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
	organizationID := c.Params("organizationID")

	listAccount, err := a.AUsecase.Find(organizationID)

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
	organizationKey := c.Params("organizationKey")
	organizationValue := c.Params("organizationValue")

	listAccount, err := a.AUsecase.FindBy(organizationKey, organizationValue)

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

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Update will change the account by given request body
func (a *AccountHandler) Update(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	var data *models.Organization

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, accountID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}
