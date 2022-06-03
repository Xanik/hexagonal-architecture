package httplib

import (
	"net/http"
	"strconv"
	"study/config/responses"
	"study/features/institution"
	httplib "study/libs/http"
	mws "study/middlewares"

	// mws "study/middlewares"

	"study/models"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase institution.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us institution.Usecase) *mux.Router {
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

	//institutionRoute
	route.HandleFunc("/v1/institutions", mwsWithAuth(handler.Create)).Methods("POST")
	route.HandleFunc("/v1/institutions", handler.FindAll).Methods("GET")
	route.HandleFunc("/v1/institutions/{Key}/{Value}", mwsWithAuth(handler.FindBy)).Methods("GET")
	route.HandleFunc("/v1/institutions/{institutionID}", mwsWithAuth(handler.Find)).Methods("GET")
	route.HandleFunc("/v1/institutions/{institutionID}", mwsWithAuth(handler.Update)).Methods("PUT")
	return route
}

// CreateAndSendMail will store the account by given request body and send a mail
func (a *AccountHandler) Create(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Institution

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Create(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error creating institution"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "institution created"}
	httplib.Response(res, resp)
}

// Find will fetch the account based on given params
func (a *AccountHandler) Find(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	institutionID := c.Params("institutionID")

	listAccount, err := a.AUsecase.Find(institutionID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching institution"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "institution details"}
	httplib.Response(res, resp)
}

// FindBy will fetch the account based on given params
func (a *AccountHandler) FindBy(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	Key := c.Params("Key")
	Value := c.Params("Value")

	listAccount, err := a.AUsecase.FindBy(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching institution"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "institution details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindAll(x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching institution"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]*models.Institution, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "institution details"}
	httplib.Response(res, resp)
}

// Update will change the account by given request body
func (a *AccountHandler) Update(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	institutionID := c.Params("institutionID")

	var data interface{}

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, institutionID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating institution"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "institution details"}
	httplib.Response(res, resp)
}
