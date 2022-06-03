package httplib

import (
	"net/http"
	"study/config/responses"
	"study/features/upload"
	"study/libs/files"
	httplib "study/libs/http"
	mws "study/middlewares"

	// mws "study/middlewares"

	"github.com/gorilla/mux"
)

// AccountHandler  represent the httphandler for account
type AccountHandler struct {
	AUsecase upload.Usecase
}

// NewAccountHandler will initialize the upload endpoint
func NewAccountHandler(route *mux.Router, us upload.Usecase) *mux.Router {
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

	//uploadRoute
	route.HandleFunc("/v1/uploads/image", mwsWithAuth(handler.ImageUpload)).Methods("POST")
	route.HandleFunc("/v1/uploads", mwsWithAuth(handler.FindAll)).Methods("GET")
	route.HandleFunc("/v1/uploads/{Key}/{Value}", mwsWithAuth(handler.FindBy)).Methods("GET")
	route.HandleFunc("/v1/uploads/{uploadID}", mwsWithAuth(handler.Find)).Methods("GET")
	route.HandleFunc("/v1/uploads/{uploadID}", mwsWithAuth(handler.Update)).Methods("PUT")

	return route
}

// ImageUpload will store the given file on an s3 bucket and return a key to access it
func (a *AccountHandler) ImageUpload(res http.ResponseWriter, req *http.Request) {
	url := files.UploadFile("url", req)

	resp := responses.GeneralResponse{Success: true, Data: url, Message: "upload completed"}
	httplib.Response(res, resp)
}

// Find will fetch the account based on given params
func (a *AccountHandler) Find(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	uploadID := c.Params("uploadID")

	listAccount, err := a.AUsecase.Find(uploadID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching upload"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "upload details"}
	httplib.Response(res, resp)
}

// FindBy will fetch the account based on given params
func (a *AccountHandler) FindBy(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}
	Key := c.Params("Key")
	Value := c.Params("Value")

	listAccount, err := a.AUsecase.FindBy(Key, Value)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching upload"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "upload details"}
	httplib.Response(res, resp)
}

// FindAll will fetch accounts based on given params
func (a *AccountHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	listAccount, err := a.AUsecase.FindAll()

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error fetching upload"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "upload details"}
	httplib.Response(res, resp)
}

// Update will change the account by given request body
func (a *AccountHandler) Update(res http.ResponseWriter, req *http.Request) {
	c := httplib.C{Res: res, Req: req}

	uploadID := c.Params("uploadID")

	var data interface{}

	c.BindJSON(&data)

	listAccount, err := a.AUsecase.Update(data, uploadID)

	if err != nil {
		resp := responses.GeneralResponse{Success: false, Error: err.Error(), Message: "error updating upload"}
		httplib.Response400(res, resp)
		return
	}

	resp := responses.GeneralResponse{Success: true, Data: listAccount, Message: "upload details"}
	httplib.Response(res, resp)
}
