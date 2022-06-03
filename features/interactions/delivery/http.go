package httplib

import (
	"fmt"
	"net/http"
	"strconv"
	"study/config/responses"
	interaction "study/features/interactions"
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
	AUsecase interaction.Usecase
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us interaction.Usecase) *mux.Router {
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

	//interactionRoute
	interactionRoute := route.PathPrefix("/v1/interactions").Subrouter()
	// Use Authentication Token For Below routes
	interactionRoute.HandleFunc("/like", mwsWithAuth(handler.LikeContent)).Methods("POST")
	interactionRoute.HandleFunc("/share", mwsWithAuth(handler.ShareContent)).Methods("POST")
	interactionRoute.HandleFunc("/bookmark", mwsWithAuth(handler.BookmarkContent)).Methods("POST")
	interactionRoute.HandleFunc("/collection", mwsWithAuth(handler.CreateCollection)).Methods("POST")
	interactionRoute.HandleFunc("/savecollection", mwsWithAuth(handler.SaveToCollection)).Methods("POST")
	interactionRoute.HandleFunc("/newcollection", mwsWithAuth(handler.CreateEmptyCollection)).Methods("POST")
	interactionRoute.HandleFunc("/collection/{collectionID}", mwsWithAuth(handler.DeleteCollection)).Methods("DELETE")
	interactionRoute.HandleFunc("/content/{collectionName}/{contentID}", mwsWithAuth(handler.DeleteContentFromCollection)).Methods("DELETE")
	interactionRoute.HandleFunc("/like/{contentID}", mwsWithAuth(handler.UnLikeContent)).Methods("DELETE")
	interactionRoute.HandleFunc("/bookmark/{contentID}", mwsWithAuth(handler.DeleteBookmark)).Methods("DELETE")
	interactionRoute.HandleFunc("/trending/{interactionID}", mwsWithAuth(handler.FindAll)).Methods("GET")
	interactionRoute.HandleFunc("/like/{accountID}", mwsWithAuth(handler.FindMyLikedContent)).Methods("GET")
	interactionRoute.HandleFunc("/followers", mwsWithAuth(handler.FindMyFollowersContent)).Methods("GET")
	interactionRoute.HandleFunc("/share/{accountID}", mwsWithAuth(handler.FindMySharedContent)).Methods("GET")
	interactionRoute.HandleFunc("/bookmark/{accountID}", mwsWithAuth(handler.FindMyBookmarkedContent)).Methods("GET")
	interactionRoute.HandleFunc("/collection", mwsWithAuth(handler.FindMyCollection)).Methods("GET")
	interactionRoute.HandleFunc("/listcollection/{accountID}", mwsWithAuth(handler.ListCollections)).Methods("GET")
	interactionRoute.HandleFunc("/collection/{collectionID}", mwsWithAuth(handler.FindByCollection)).Methods("GET")
	interactionRoute.HandleFunc("/{Key}/{Value}", mwsWithAuth(handler.FindBy)).Methods("GET")
	interactionRoute.HandleFunc("/{interactionID}", mwsWithAuth(handler.Find)).Methods("GET")
	interactionRoute.HandleFunc("/collection/{collectionID}", mwsWithAuth(handler.UpdateCollectionName)).Methods("PUT")
	interactionRoute.HandleFunc("/{interactionID}", mwsWithAuth(handler.Update)).Methods("PUT")

	return route
}

// CreateEmptyCollection will store the account by given request body
func (a *AccountHandler) CreateEmptyCollection(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Interactions

	c.BindJSON(&data)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	data.Type = "collection"

	data.UserID = bson.ObjectIdHex(accountID.(string))

	listAccount, err := a.AUsecase.CreateEmptyCollection(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// CreateCollection will store the account by given request body
func (a *AccountHandler) CreateCollection(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Interactions

	c.BindJSON(&data)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	data.Type = "collection"

	data.UserID = bson.ObjectIdHex(accountID.(string))

	listAccount, err := a.AUsecase.CreateCollection(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// SaveToCollection will store the account by given request body
func (a *AccountHandler) SaveToCollection(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Interactions

	c.BindJSON(&data)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	data.Type = "collection"

	data.UserID = bson.ObjectIdHex(accountID.(string))

	listAccount, err := a.AUsecase.SaveToCollection(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// ShareContent will store the account by given request body
func (a *AccountHandler) ShareContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Interactions

	c.BindJSON(&data)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	data.Type = "share"

	data.UserID = bson.ObjectIdHex(accountID.(string))

	listAccount, err := a.AUsecase.ShareContent(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// LikeContent will store the account by given request body
func (a *AccountHandler) LikeContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Interactions

	c.BindJSON(&data)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	data.Type = "like"

	data.UserID = bson.ObjectIdHex(accountID.(string))

	listAccount, err := a.AUsecase.LikeContent(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// UnLikeContent will store the account by given request body
func (a *AccountHandler) UnLikeContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	contentID := c.Params("contentID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.UnLikeContent(accountID.(string), contentID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// BookmarkContent will store the account by given request body
func (a *AccountHandler) BookmarkContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Interactions

	c.BindJSON(&data)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	data.Type = "bookmark"

	data.UserID = bson.ObjectIdHex(accountID.(string))

	listAccount, err := a.AUsecase.BookmarkContent(data)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// DeleteBookmark will store the account by given request body
func (a *AccountHandler) DeleteBookmark(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	contentID := c.Params("contentID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.DeleteBookmark(accountID.(string), contentID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// FindByCollection will fetch the account based on given params
func (a *AccountHandler) FindByCollection(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	collectionID := c.Params("collectionID")

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.FindByCollection(collectionID, x, y, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Find will fetch the account based on given params
func (a *AccountHandler) Find(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	interactionID := c.Params("interactionID")

	listAccount, err := a.AUsecase.Find(interactionID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
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
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	interactionID := c.Params("interactionID")

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindAll(interactionID, x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindMyLikedContent will fetch accounts based on given params
func (a *AccountHandler) FindMyLikedContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	userID := user["_id"]

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindMyLikedContent(accountID, userID.(string), x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindMyFollowersContent will fetch accounts based on given params
func (a *AccountHandler) FindMyFollowersContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	userID := user["_id"]

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindMyLikedContent(accountID, userID.(string), x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindMySharedContent will fetch accounts based on given params
func (a *AccountHandler) FindMySharedContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	userID := user["_id"]

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindMySharedContent(accountID, userID.(string), x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindMyBookmarkedContent will fetch accounts based on given params
func (a *AccountHandler) FindMyBookmarkedContent(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	userID := user["_id"]

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindMyBookmarkedContent(accountID, userID.(string), x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindMyCollection will fetch accounts based on given params
func (a *AccountHandler) FindMyCollection(res http.ResponseWriter, req *http.Request) {
	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.FindMyCollection(accountID.(string), x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// ListCollections will fetch accounts based on given params
func (a *AccountHandler) ListCollections(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	accountID := c.Params("accountID")

	skip := req.FormValue("skip")

	limit := req.FormValue("limit")

	x, _ := strconv.Atoi(skip)

	y, _ := strconv.Atoi(limit)

	listAccount, err := a.AUsecase.ListCollections(accountID, x, y)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]interface{}, 0)
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

	interactionID := c.Params("interactionID")

	var data interface{}

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, interactionID, "interactions")

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// UpdateCollectionName will change the account by given request body
func (a *AccountHandler) UpdateCollectionName(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data models.Interactions

	collectionID := c.Params("collectionID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.UpdateCollectionName(collectionID, data.Name, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// DeleteCollection will store the account by given request body
func (a *AccountHandler) DeleteCollection(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	collectionID := c.Params("collectionID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.DeleteCollection(collectionID, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}

// DeleteContentFromCollection will store the account by given request body
func (a *AccountHandler) DeleteContentFromCollection(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	collectionName := c.Params("collectionName")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	contentID := c.Params("contentID")

	listAccount, err := a.AUsecase.DeleteContentFromCollection(collectionName, contentID, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "created"}
	httplib.Response(res, resp)
}
