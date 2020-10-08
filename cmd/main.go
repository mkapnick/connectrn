package main

import (
	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	mw "gitlab.com/michaelk99/connectrn/internal/middleware"
	"gitlab.com/michaelk99/connectrn/internal/token/jwthmac"
	// "gitlab.com/michaelk99/connectrn/internal/validator"
	"gitlab.com/michaelk99/connectrn/services/restaurant"
	"gitlab.com/michaelk99/connectrn/services/restaurant/handlers"
	"gitlab.com/michaelk99/connectrn/services/restaurant/postgres"

	// v9 "gopkg.in/go-playground/validator.v9"
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
	JWTSecret      string = "fde5247c0262798a9c"
	PGDriver       string = "postgres"
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

	// gcds = golf course data source
	rds := postgres.NewRestaurantStore(dbConn)

	// v9Validator := v9.New()
	// validator := validator.NewValidator(v9Validator)

	// create our restaurant service
	restaurantService := restaurant.NewService(rds)

	// create our auth request middleware
	authRequest := mw.NewAuthRequest(jwthmacStore)

	// create our http mux and add routes [using gorilla for convenience]
	r := mux.NewRouter()

	/////////////////// Routes ///////////////////
	r.HandleFunc("/api/v1/restaurants/{restaurant_id}/", authRequest.Auth(handlers.Fetch(restaurantService))).Methods("GET")

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
