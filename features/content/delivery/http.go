package httplib

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"study/config/responses"
	"study/features/content"
	httplib "study/libs/http"
	mws "study/middlewares"
	"study/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/robfig/cron/v3"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase content.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us content.Usecase) *mux.Router {
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

	//contentRoute
	contentRoute := route.PathPrefix("/v1/contents").Subrouter()
	contentRoute.HandleFunc("", mwsWithAuth(handler.FindAll)).Methods("GET")
	// Use Authentication Token For Below routes
	contentRoute.HandleFunc("", mwsWithAuth(handler.Create)).Methods("POST")
	contentRoute.HandleFunc("/account/{accountID}", mwsWithAuth(handler.GetUsersSuggestedContent)).Methods("GET")
	contentRoute.HandleFunc("/trending", mwsWithAuth(handler.GetTrendingContent)).Methods("GET")
	contentRoute.HandleFunc("/tag", mwsWithAuth(handler.GetContentsByTag)).Methods("GET")
	contentRoute.HandleFunc("/{Key}/{Value}", mwsWithAuth(handler.FindBy)).Methods("GET")
	contentRoute.HandleFunc("/{contentID}", mwsWithAuth(handler.Find)).Methods("GET")
	contentRoute.HandleFunc("/{contentID}", mwsWithAuth(handler.Update)).Methods("PUT")

	cr := cron.New()
	cr.AddFunc("@hourly", handler.WebCrawler)
	cr.AddFunc("@hourly", func() { fmt.Println("Every day") })
	cr.Start()

	return route
}

// Create will store the account by given request body
func (a *AccountHandler) Create(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Content

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
	contentID := c.Params("contentID")

	listAccount, err := a.AUsecase.Find(contentID)

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

// GetTrendingContent will fetch accounts based on given params
func (a *AccountHandler) GetTrendingContent(res http.ResponseWriter, req *http.Request) {
	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.FindAllWithQuery(x, y, accountID.(string))

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
	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindAll(x, y)

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

	contentID := c.Params("contentID")

	var data *models.Content

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, contentID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// GetUsersSuggestedContent will fetch the account based on given params
func (a *AccountHandler) GetUsersSuggestedContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	accountID := c.Params("accountID")
	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.GetUsersSuggestedContent(accountID, x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// GetUsersSuggestedContent will fetch the account based on given params
func (a *AccountHandler) GetContentsByTag(res http.ResponseWriter, req *http.Request) {
	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	query := req.FormValue("query")

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.GetContentsByTag(accountID.(string), query, x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawler will store the account by given request body
func (a *AccountHandler) CrawlerRss(address string) {
	resp, err := http.Get(address)
	if err != nil {
		fmt.Printf("Error GET: %v\n", err)
		return
	}
	defer resp.Body.Close()

	rss := models.Rss{}

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		fmt.Printf("Error Decode: %v\n", err)
		return
	}
	var data models.Item
	for _, item := range rss.Channel.Items {
		data = item
		fmt.Println(data)
		data.Company = address
		a.AUsecase.Crawler(data)
	}
}

func (a *AccountHandler) WebCrawler() {
	for _, value := range models.Address {
		a.CrawlerRss(value)
	}
}
