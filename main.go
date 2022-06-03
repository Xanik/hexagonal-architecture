package main

import (
	"log"
	"net/http"
	"study/config"
	dbs "study/config/db"
	_accountHttpDeliver "study/features/account/delivery"
	_accountRepo "study/features/account/repository"
	_accountUsecase "study/features/account/service"

	loger "github.com/sirupsen/logrus"

	_interestHttpDeliver "study/features/interest/delivery"
	_interestRepo "study/features/interest/repository"
	_interestUsecase "study/features/interest/service"

	_organizationHttpDeliver "study/features/organization/delivery"
	_organizationRepo "study/features/organization/repository"
	_organizationUsecase "study/features/organization/service"

	_chatHttpDeliver "study/features/discuss/delivery"
	_chatRepo "study/features/discuss/repository"
	_chatUsecase "study/features/discuss/service"

	_institutionHttpDeliver "study/features/institution/delivery"
	_institutionRepo "study/features/institution/repository"
	_institutionUsecase "study/features/institution/service"

	_networkHttpDeliver "study/features/network/delivery"
	_networkRepo "study/features/network/repository"
	_networkUsecase "study/features/network/service"

	_uploadHttpDeliver "study/features/upload/delivery"
	_uploadRepo "study/features/upload/repository"
	_uploadUsecase "study/features/upload/service"

	_contentHttpDeliver "study/features/content/delivery"
	_contentRepo "study/features/content/repository"
	_contentUsecase "study/features/content/service"

	_interactionsHttpDeliver "study/features/interactions/delivery"
	_interactionsRepo "study/features/interactions/repository"
	_interactionsUsecase "study/features/interactions/service"

	_searchHttpDeliver "study/features/search/delivery"
	_searchRepo "study/features/search/repository"
	_searchUsecase "study/features/search/service"

	_courseHttpDeliver "study/features/course/delivery"
	_courseRepo "study/features/course/repository"
	_courseUsecase "study/features/course/service"

	_developerHttpDeliver "study/features/developer/delivery"
	_developerRepo "study/features/developer/repository"
	_developerUsecase "study/features/developer/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

func init() {

	loger.SetFormatter(&loger.TextFormatter{})

	log.Printf("%s environment started", config.Env.Env)

}

func main() {
	var db *mgo.Session

	// Start Database Connection
	db = dbs.NewClient()

	port := config.Env.Port

	route := mux.NewRouter()

	//	APPLY MIDDLEWARES
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	// account Service
	accountRepo := _accountRepo.NewMongoAccountRepository(db)
	accountUsecase := _accountUsecase.NewAccountUsecase(accountRepo)
	_accountHttpDeliver.NewAccountHandler(route, accountUsecase)

	// account Service
	interestRepo := _interestRepo.NewMongoAccountRepository(db)
	interestUsecae := _interestUsecase.NewAccountUsecase(interestRepo)
	_interestHttpDeliver.NewAccountHandler(route, interestUsecae)

	// organization Service
	organizationRepo := _organizationRepo.NewMongoAccountRepository(db)
	organizationUsecase := _organizationUsecase.NewAccountUsecase(organizationRepo)
	_organizationHttpDeliver.NewAccountHandler(route, organizationUsecase)

	// chat Service
	_chatRepo := _chatRepo.NewMongoAccountRepository(db)
	chatUsecase := _chatUsecase.NewAccountUsecase(_chatRepo)
	_chatHttpDeliver.NewAccountHandler(route, chatUsecase)

	// institution.go Service
	_institutionRepo := _institutionRepo.NewMongoInstitutionRepository(db)
	institutionUsecase := _institutionUsecase.NewInstitutionUsecase(_institutionRepo)
	_institutionHttpDeliver.NewAccountHandler(route, institutionUsecase)

	// network Service
	_networkRepo := _networkRepo.NewMongoNetworkRepository(db)
	networkUsecase := _networkUsecase.NewNetworkUsecase(_networkRepo)
	_networkHttpDeliver.NewAccountHandler(route, networkUsecase)

	// upload Service
	_uploadRepo := _uploadRepo.NewMongoUploadRepository(db)
	uploadUsecase := _uploadUsecase.NewUploadUsecase(_uploadRepo)
	_uploadHttpDeliver.NewAccountHandler(route, uploadUsecase)

	// content Service
	_contentRepo := _contentRepo.NewMongoAccountRepository(db)
	contentUsecase := _contentUsecase.NewAccountUsecase(_contentRepo)
	_contentHttpDeliver.NewAccountHandler(route, contentUsecase)

	// interactions Service
	_interactionsRepo := _interactionsRepo.NewMongoAccountRepository(db)
	interactionsUsecase := _interactionsUsecase.NewAccountUsecase(_interactionsRepo)
	_interactionsHttpDeliver.NewAccountHandler(route, interactionsUsecase)

	// search Service
	_searchRepo := _searchRepo.NewMongoAccountRepository(db)
	searchUsecase := _searchUsecase.NewAccountUsecase(_searchRepo)

	// Check Data
	//TODO: Remove Once Env Secret Has been fixed
	searchUsecase.IndexContents("contents")
	searchUsecase.IndexAccounts("accounts")
	searchUsecase.IndexCourses("courses")
	_searchHttpDeliver.NewAccountHandler(route, searchUsecase)

	// course Service
	_courseRepo := _courseRepo.NewMongoAccountRepository(db)
	courseUsecase := _courseUsecase.NewAccountUsecase(_courseRepo)
	_courseHttpDeliver.NewAccountHandler(route, courseUsecase)

	// developer Service
	_developerRepo := _developerRepo.NewMongoRepository(db)
	developerUsecase := _developerUsecase.NewAccountUsecase(_developerRepo)
	_developerHttpDeliver.NewAccountHandler(route, developerUsecase)

	log.Printf("server running at %s", port)
	//apmhttp.WrapClient(http.DefaultClient)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(route)))

}
