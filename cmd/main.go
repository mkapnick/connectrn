package main

import (
	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
	"net/url"

	_ "github.com/lib/pq"
	mw "gitlab.com/michaelk99/connectrn/internal/middleware"
	"gitlab.com/michaelk99/connectrn/internal/token/jwthmac"
	"gitlab.com/michaelk99/connectrn/internal/validator"
	ajwthmac "gitlab.com/michaelk99/connectrn/services/account/jwthmac"
	"time"

	"gitlab.com/michaelk99/connectrn/services/account"
	"gitlab.com/michaelk99/connectrn/services/profile"
	"gitlab.com/michaelk99/connectrn/services/restaurant"

	ahandlers "gitlab.com/michaelk99/connectrn/services/account/handlers"
	phandlers "gitlab.com/michaelk99/connectrn/services/profile/handlers"
	rhandlers "gitlab.com/michaelk99/connectrn/services/restaurant/handlers"

	apostgres "gitlab.com/michaelk99/connectrn/services/account/postgres"
	ppostgres "gitlab.com/michaelk99/connectrn/services/profile/postgres"
	rpostgres "gitlab.com/michaelk99/connectrn/services/restaurant/postgres"

	v9 "gopkg.in/go-playground/validator.v9"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	HTTPListenAddr string = "0.0.0.0:3000"
	PGConnString   string = "host=localhost port=5432 user=postgres dbname=connectrn  sslmode=disable"
	PGDriver       string = "postgres"
	JWTExp         int64  = 65920000000
	JWTSecret      string = "fde5247c0262798a9c"
	JWTIssuer      string = "account"
	ProfileURL     string = "http://localhost:3000/api/v1/profile/"
)

// root is the root route, used for k8s health checks
func root(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "healthy")
}

func main() {
	// create token store, which will be used to handle jwt authentication
	jwthmacStore := jwthmac.NewTokenStore([]byte(JWTSecret), "HMAC")

	// create conn to db
	dbConn, err := sqlx.Connect(PGDriver, PGConnString)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// create profile url
	profURL, err := url.Parse(ProfileURL)
	if err != nil {
		log.Fatalf("Profile URL is not a valid url %s", profURL)
	}
	pc := profile.NewClient(profURL)

	// create data sources
	ads := apostgres.NewAccountStore(dbConn)
	pds := ppostgres.NewProfileStore(dbConn)
	rds := rpostgres.NewRestaurantStore(dbConn)

	v9Validator := v9.New()
	validator := validator.NewValidator(v9Validator)

	// create token creator we will use to isssue tokens via login requests
	exp := time.Duration(JWTExp) * time.Millisecond
	tc := ajwthmac.NewCreator(JWTSecret, JWTIssuer, exp)

	// create our services
	accountService := account.NewService(ads, tc, pc)
	profileService := profile.NewService(pds)
	restaurantService := restaurant.NewService(rds)

	// create our auth request middleware
	authRequest := mw.NewAuthRequest(jwthmacStore)

	// create our http mux and add routes [using gorilla for convenience]
	r := mux.NewRouter()

	/////////////////// Routes ///////////////////
	// acount routes
	r.HandleFunc("/api/v1/account/", authRequest.Auth(ahandlers.CRUD(accountService)))
	r.HandleFunc("/api/v1/account/signup/", ahandlers.SignUp(validator, accountService))
	r.HandleFunc("/api/v1/account/login/", ahandlers.Login(accountService))

	// profile routes
	r.HandleFunc("/api/v1/profile/", authRequest.Auth(phandlers.CRUD(validator, profileService)))

	// restaurant routes
	r.HandleFunc("/api/v1/restaurants/{restaurant_id}/", authRequest.Auth(rhandlers.Fetch(restaurantService))).Methods("GET")

	// not found route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Route not found\n")
			return
		}
		root(w, r)
	})

	// create server and launch in go routine
	s := http.Server{
		Addr:    HTTPListenAddr,
		Handler: r,
	}

	sigChan := make(chan os.Signal)
	errChan := make(chan error)

	go func(errChan chan error, s http.Server) {
		log.Printf("launching http server on %v", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}(errChan, s)

	// register signal handler
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// block on sigChan or errChan
	select {
	case sig := <-sigChan:
		log.Printf("received signal %v. attempting graceful shutdown of server", sig.String())

		err = s.Shutdown(nil)
		if err != nil {
			log.Printf("graceful shutdown of server unsuccessful: %v", err)
			os.Exit(1)
		}

		log.Printf("graceful shutdown of server successful")
		os.Exit(0)
	case e := <-errChan:
		log.Printf("failed to launched http server: %v", e)
		os.Exit(1)
	}
}
