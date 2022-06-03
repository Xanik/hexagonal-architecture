package httplib

import (
	"fmt"
	"net/http"
	"study/config/responses"
	"study/features/developer"
	httplib "study/libs/http"
	mws "study/middlewares"
	"study/models"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase developer.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us developer.Usecase) *mux.Router {
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

	//developerRoute
	developerRoute := route.PathPrefix("/v1/developers").Subrouter()
	developerRoute.HandleFunc("/request", handler.RequestAccess).Methods("POST")
	developerRoute.HandleFunc("/auth", handler.AuthenticateAccount).Methods("POST")
	developerRoute.HandleFunc("/content", mwsWithAuth(handler.GetAllContents)).Methods("GET")
	developerRoute.HandleFunc("/course", mwsWithAuth(handler.GetAllCourses)).Methods("GET")
	developerRoute.HandleFunc("/interest", mwsWithAuth(handler.GetContentByInterest)).Methods("POST")
	developerRoute.HandleFunc("/search", mwsWithAuth(handler.SearchContent)).Methods("GET")

	return route
}

// RequestAccess will store the account by given request body
func (a *AccountHandler) RequestAccess(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Developer

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.RequestAccess(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error creating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// AuthenticateAccount will fetch the account based on given params
func (a *AccountHandler) AuthenticateAccount(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Developer

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.AuthenticateAccount(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) GetAllContents(res http.ResponseWriter, req *http.Request) {

	listAccount, err := a.AUsecase.GetAllContents()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]*models.Content, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "content details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// GetContentByInterest will fetch the account based on given params
func (a *AccountHandler) GetContentByInterest(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var developer *models.DeveloperRequest

	c.BindJSON(&developer)

	listAccount, err := a.AUsecase.GetContentByInterest(developer)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]*models.Content, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "content details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Search will fetch accounts based on given params
func (a *AccountHandler) SearchContent(res http.ResponseWriter, req *http.Request) {

	query := req.FormValue("query")

	listAccount, err := a.AUsecase.SearchContent(query)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]map[string]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "content details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) GetAllCourses(res http.ResponseWriter, req *http.Request) {

	listAccount, err := a.AUsecase.GetAllCourses()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]*models.Course, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "course details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}
