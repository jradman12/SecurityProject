package startup

import (
	userGw "common/module/proto/user_service"
	"context"
	"fmt"
	"gateway/module/application/helpers"
	"gateway/module/application/services"
	"gateway/module/domain/model"
	"gateway/module/domain/repositories"
	"gateway/module/infrastructure/handlers"
	"gateway/module/infrastructure/persistance"
	cfg "gateway/module/startup/config"
	gorilla_handlers "github.com/gorilla/handlers"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux // Part of grpcGateway library
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

func (server *Server) initHandlers() {
	//Povezuje sa grpc generisanim fajlovima
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)

	err := userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	if err != nil {
		panic(err)
	}

}

//Gateway ima svoje endpointe
func (server *Server) initCustomHandlers() {

	l := log.New(os.Stdout, "products-api ", log.LstdFlags) // Logger koji dajemo handlerima
	db = server.SetupDatabase()
	userRepo := server.InitUserRepo(db)
	userService := server.InitUserService(l, userRepo)

	validator := validator.New()

	passwordUtil := &helpers.PasswordUtil{}
	authHandler := handlers.NewAuthenticationHandler(l, userService, validator, passwordUtil)
	authHandler.Init(server.mux)
}

func (server *Server) Start() {
	cors := gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"https://localhost:4200", "https://localhost:4200/**", "http://localhost:4200", "http://localhost:4200/**", "http://localhost:8080/**"}),
		gorilla_handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorilla_handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "Authorization", "Access-Control-Allow-Origin", "*"}),
		gorilla_handlers.AllowCredentials(),
	)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), cors(muxMiddleware(server))))
}
func muxMiddleware(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.mux.ServeHTTP(w, r)
	})
}

func (server *Server) InitUserService(l *log.Logger, repo repositories.UserRepository) *services.UserService {
	return services.NewUserService(l, repo)
}
func (server *Server) InitUserRepo(d *gorm.DB) repositories.UserRepository {
	return persistance.NewUserRepositoryImpl(db)
}

var db *gorm.DB

func (server *Server) SetupDatabase() *gorm.DB {

	host := os.Getenv("HOST")
	port := os.Getenv("PG_DBPORT")
	user := os.Getenv("PG_USER")
	dbname := os.Getenv("XML_DB_NAME")
	password := os.Getenv("PG_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	db.AutoMigrate(&model.User{}) //This will not remove columns
	//db.Create(users) // Use this only once to populate db with data

	return db
}