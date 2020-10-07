package main

import (
	"fmt"
	"github.com/crgimenes/goconfig"
	"github.com/jmoiron/sqlx"
	"net/url"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	mw "gitlab.com/michaelk99/connectrn/internal/middleware"
	"gitlab.com/michaelk99/connectrn/internal/token/jwthmac"
	"gitlab.com/michaelk99/connectrn/internal/validator"
	"gitlab.com/michaelk99/connectrn/services/reserve"
	"gitlab.com/michaelk99/connectrn/services/reserve/handlers"
	"gitlab.com/michaelk99/connectrn/services/reserve/postgres"
	v9 "gopkg.in/go-playground/validator.v9"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Config this struct is using the goconfig library for simple flag and env var
// parsing. See: https://github.com/crgimenes/goconfig
type Config struct {
	HTTPListenAddr  string `cfgDefault:"0.0.0.0:3000" cfg:"MAIN_HTTP_LISTEN_ADDR"`
	JWTSecret       string `cfgDefault:"fde5247c0262798a9c" cfg:"JWT_SECRET"`
	PGConnString    string `cfgDefault:"host=localhost port=5432 user=postgres dbname=connectrn  sslmode=disable" cfg:"POSTGRES_CONN_STRING"`
	PGDriver        string `cfgDefault:"postgres" cfg:"POSTGRES_DRIVER"`
}

// root is the root route, used for k8s health checks
func root(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "healthy")
}

func main() {
	// parse our config
	conf := Config{}
	err := goconfig.Parse(&conf)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// create token store, which will be used to handle jwt authentication
	jwthmacStore := jwthmac.NewTokenStore([]byte(conf.JWTSecret), "HMAC")

	// create conn to db
	dbConn, err := sqlx.Connect(conf.PGDriver, conf.PGConnString)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// gcds = golf course data source
	gcds := postgres.NewGolfCourseStore(dbConn)

	v9Validator := v9.New()
	validator := validator.NewValidator(v9Validator)

	// create our reserve service
	reserveService := reserve.NewGolfCourseService(gcds, cache, nc, sc)

	// create our auth request middleware
	authRequest := mw.NewAuthRequest(jwthmacStore)

	// create our http mux and add routes [using gorilla for convenience]
	r := mux.NewRouter()

	/////////////////// Routes ///////////////////
	r.HandleFunc("/api/v2/golf-courses/{golf_course_id}/tee-times/{tee_time_id}/user-tee-times/{user_tee_time_id}/", authRequest.Auth(handlers.UpdateUserTeeTime(validator, reserveService))).Methods("PATCH")

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
		Addr:    conf.HTTPListenAddr,
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
