package delivery

import (
	"log"
	"net/http"
	"strings"
	"study/config/responses"
	chat "study/features/discuss"
	httplib "study/libs/http"
	mws "study/middlewares"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"

	// mws "study/middlewares"

	"study/models"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase chat.Usecase
}

var store map[string]interface{}

func init() {

	store = make(map[string]interface{})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// NewAccountHandler will initialize the account/ resources endpoint
func NewAccountHandler(route *mux.Router, us chat.Usecase) *mux.Router {
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

	//chatRoute
	chatRoute := route.PathPrefix("/v1/chats").Subrouter()
	//socketRoute

	chatRoute.HandleFunc("/socket/{userID}", handler.chatHandler)
	// mwsWithAuth(authRoute)
	chatRoute.HandleFunc("", mwsWithAuth(handler.CreateDiscussion)).Methods("POST")
	chatRoute.HandleFunc("", mwsWithAuth(handler.FindAll)).Methods("GET")
	chatRoute.HandleFunc("/tag", mwsWithAuth(handler.FindAllAndGroup)).Methods("GET")
	chatRoute.HandleFunc("/userID/{Value}", mwsWithAuth(handler.FindDiscussionByUserID)).Methods("GET")
	chatRoute.HandleFunc("/course/{CourseID}", mwsWithAuth(handler.FindByCourse)).Methods("GET")
	chatRoute.HandleFunc("/{Key}/{Value}", mwsWithAuth(handler.FindBy)).Methods("GET")
	chatRoute.HandleFunc("/comment", mwsWithAuth(handler.Comment)).Methods("PUT")
	chatRoute.HandleFunc("/comment/{Key}/{Value}", mwsWithAuth(handler.CommentByID)).Methods("GET")
	chatRoute.HandleFunc("/trending", mwsWithAuth(handler.TrendingDiscussions)).Methods("GET")
	chatRoute.HandleFunc("/{chatID}", mwsWithAuth(handler.Find)).Methods("GET")
	chatRoute.HandleFunc("/user/{chatID}", mwsWithAuth(handler.UpdateDiscussionUsers)).Methods("PUT")
	chatRoute.HandleFunc("/unsubscribe/{chatID}", mwsWithAuth(handler.UnSubscribe)).Methods("PUT")
	chatRoute.HandleFunc("/subscribe/{chatID}", mwsWithAuth(handler.Subscribe)).Methods("PUT")
	chatRoute.HandleFunc("/{chatID}", mwsWithAuth(handler.Update)).Methods("PUT")
	chatRoute.HandleFunc("/comment/{commentID}", mwsWithAuth(handler.DeleteComment)).Methods("DELETE")
	chatRoute.HandleFunc("/comment/{commentID}", mwsWithAuth(handler.UpdateComment)).Methods("PUT")
	chatRoute.HandleFunc("/chat/{chatID}", mwsWithAuth(handler.DeleteDiscussion)).Methods("DELETE")

	return chatRoute
}

func (a *AccountHandler) chatHandler(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["userID"]

	conn, err := upgrader.Upgrade(w, r, nil)

	println(conn.RemoteAddr().String())
	if err != nil {
		log.Println(err.Error())
		return
	}

	println(store)

	go a.AUsecase.ReadData(conn)

	go a.AUsecase.BroadcastMessage(store)

	//saves user's connection to global store
	store[userID] = conn

}

// Create will store the account by given request body
func (a *AccountHandler) CreateDiscussion(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Chat

	c.BindJSON(&data)

	data.ID = bson.NewObjectId()

	data.CreatedAT = time.Now()

	data.UpdatedAt = time.Now()

	listAccount, err := a.AUsecase.CreateDiscussion(data)

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
	chatID := c.Params("chatID")

	listAccount, err := a.AUsecase.Find(chatID)

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

// FindByCourse will fetch the account based on given params
func (a *AccountHandler) FindByCourse(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	CourseID := c.Params("CourseID")

	listAccount, err := a.AUsecase.FindByCourse(CourseID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]map[string]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) FindAll(res http.ResponseWriter, req *http.Request) {

	var match bson.M

	filter := req.FormValue("filter")

	if filter == "" {
		match = bson.M{"$match": bson.M{"type": "public", "course_id": nil}}
	}

	if filter != "" {
		x := strings.Split(filter, "|")
		key := x[0]
		value := x[1]
		match = bson.M{"$match": bson.M{key: value, "course_id": nil}}
	}

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	match2 := bson.M{"$match": bson.M{"type": "private", "users": bson.ObjectIdHex(accountID.(string)), "course_id": nil}}

	listAccount, err := a.AUsecase.FindAll(match, match2)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]map[string]interface{}, 0)
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

	chatID := c.Params("chatID")

	var data interface{}

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, chatID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// UpdateDiscussionUsers will change the account by given request body
func (a *AccountHandler) UpdateDiscussionUsers(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	chatID := c.Params("chatID")

	var data *models.Chat

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.UpdateDiscussionUsers(data.ID.Hex(), chatID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// Comment will change the account by given request body
func (a *AccountHandler) Comment(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	var data *models.Comment

	c.BindJSON(&data)

	data.ID = bson.NewObjectId()

	data.CreatedAT = time.Now()

	listAccount, err := a.AUsecase.AddComment(data, data.ID, data.DiscussionID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// TrendingDiscussions will fetch accounts based on given params
func (a *AccountHandler) TrendingDiscussions(res http.ResponseWriter, req *http.Request) {

	listAccount, err := a.AUsecase.TrendingDiscussions()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]map[string]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// CommentByID will fetch accounts based on given params
func (a *AccountHandler) CommentByID(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	Key := c.Params("Key")
	Value := c.Params("Value")

	listAccount, err := a.AUsecase.FindCommentsBy(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]map[string]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindDiscussionByUserID will fetch accounts based on given params
func (a *AccountHandler) FindDiscussionByUserID(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	Key := "comments"
	Value := c.Params("Value")

	match := bson.M{"$match": bson.M{"$or": []interface{}{
		bson.M{"created_by": bson.ObjectIdHex(Value)},
		bson.M{Key: bson.ObjectIdHex(Value)}}}}

	listAccount, err := a.AUsecase.FindDiscussionByUserID(match)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
		httplib.Response400(res, resp)
		return
	}

	if listAccount == nil {
		listAccount = make([]map[string]interface{}, 0)
		resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "network details"}
		httplib.Response(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// FindAllAndGroup will fetch accounts based on given params
func (a *AccountHandler) FindAllAndGroup(res http.ResponseWriter, req *http.Request) {

	listAccount, err := a.AUsecase.FindAllAndGroup()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching"}
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

// Subscribe will change the account by given request body
func (a *AccountHandler) Subscribe(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	chatID := c.Params("chatID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.Subscribe(accountID.(string), chatID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// UnSubscribe will change the account by given request body
func (a *AccountHandler) UnSubscribe(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	chatID := c.Params("chatID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.UnSubscribe(accountID.(string), chatID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}

// DeleteComment will store the account by given request body
func (a *AccountHandler) DeleteComment(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	commentID := c.Params("commentID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.DeleteComment(commentID, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "deleted"}
	httplib.Response(res, resp)
}

// DeleteDiscussion will store the account by given request body
func (a *AccountHandler) DeleteDiscussion(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	chatID := c.Params("chatID")

	user := context.Get(req, "decoded").(jwt.MapClaims)

	accountID := user["_id"]

	listAccount, err := a.AUsecase.DeleteDiscussion(chatID, accountID.(string))

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: "error fetching", Message: err.Error()}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "deleted"}
	httplib.Response(res, resp)
}

// UpdateComment will change the account by given request body
func (a *AccountHandler) UpdateComment(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	commentID := c.Params("commentID")

	var data interface{}

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.UpdateComment(data, commentID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "details"}
	httplib.Response(res, resp)
}
