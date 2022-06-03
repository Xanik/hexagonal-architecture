package httplib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"study/config"
	"study/config/responses"
	"study/features/search"
	httplib "study/libs/http"
	mws "study/middlewares"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/robfig/cron/v3"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase search.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us search.Usecase) *mux.Router {
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
	searchRoute := route.PathPrefix("/v1/search").Subrouter()
	// Use Authentication Token For Below routes
	searchRoute.HandleFunc("/interest", mwsWithAuth(handler.SearchInterest)).Methods("GET")
	searchRoute.HandleFunc("/account", mwsWithAuth(handler.SearchAccount)).Methods("GET")
	searchRoute.HandleFunc("/content", mwsWithAuth(handler.SearchContent)).Methods("GET")
	searchRoute.HandleFunc("/course", mwsWithAuth(handler.SearchCourse)).Methods("GET")
	searchRoute.HandleFunc("/{searchDoc}/{searchKey}/{searchValue}", mwsWithAuth(handler.Search)).Methods("GET")

	cr := cron.New()
	cr.AddFunc("@midnight", handler.Index)
	cr.AddFunc("@midnight", func() { fmt.Println("Every day") })
	cr.Start()

	return route
}

// Search will fetch accounts based on given params
func (a *AccountHandler) Search(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	doc := c.Params("searchDoc")
	key := c.Params("searchKey")
	value := c.Params("searchValue")
	listAccount, err := a.AUsecase.SearchElastic(doc, key, value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Search will fetch accounts based on given params
func (a *AccountHandler) SearchInterest(res http.ResponseWriter, req *http.Request) {

	query := req.FormValue("query")

	listAccount, err := a.AUsecase.SearchInterest("interests", query)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]bson.M, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "search details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Search will fetch accounts based on given params
func (a *AccountHandler) SearchAccount(res http.ResponseWriter, req *http.Request) {

	query := req.FormValue("query")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.SearchAccount("accounts", query, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]bson.M, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "search details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Search will fetch accounts based on given params
func (a *AccountHandler) SearchContent(res http.ResponseWriter, req *http.Request) {

	query := req.FormValue("query")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.SearchContent("contents", query, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]bson.M, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "search details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Search will fetch accounts based on given params
func (a *AccountHandler) SearchCourse(res http.ResponseWriter, req *http.Request) {

	query := req.FormValue("query")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.SearchCourse("courses", query, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]bson.M, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "search details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

func EmptyElastic(indexName string) {

	url := fmt.Sprintf("%s%s", config.Env.ElasticsearchURL, indexName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := http.Client{
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(data))
}

// IndexAccounts will fetch accounts based on given params
func (a *AccountHandler) IndexAccounts() {

	indexName := config.Env.Env + "-accounts"

	//Delete Index To Bulk ReIndex
	EmptyElastic(indexName)

	listAccount, err := a.AUsecase.IndexAccounts(indexName)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		fmt.Println(resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	fmt.Println(resp)
}

// IndexContents will fetch accounts based on given params
func (a *AccountHandler) IndexContents() {

	indexName := config.Env.Env + "-contents"

	//Delete Index To Bulk ReIndex
	EmptyElastic(indexName)

	listAccount, err := a.AUsecase.IndexContents(indexName)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		fmt.Println(resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	fmt.Println(resp)
}

// IndexCourses will fetch accounts based on given params
func (a *AccountHandler) IndexCourses() {

	indexName := config.Env.Env + "-courses"

	//Delete Index To Bulk ReIndex
	EmptyElastic(indexName)

	listAccount, err := a.AUsecase.IndexCourses(indexName)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		fmt.Println(resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	fmt.Println(resp)
}

// IndexInterest will fetch accounts based on given params
func (a *AccountHandler) IndexInterest() {

	indexName := config.Env.Env + "-interests"

	//Delete Index To Bulk ReIndex
	EmptyElastic(indexName)

	listAccount, err := a.AUsecase.IndexInterests(indexName)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		fmt.Println(resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	fmt.Println(resp)
}

func (a *AccountHandler) Index() {
	// ReIndex Account To Run Every Midnight
	fmt.Println("Ran Today")
}
