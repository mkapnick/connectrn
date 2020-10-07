package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"gitlab.com/michaelk99/connectrn/internal/token"
)

const (
	// Session session
	Session = "Session"
	// Token token
	Token = "Token"
)

// AuthRequest is a public interface for implementing proper auth functionality
type AuthRequest interface {
	Auth(next http.HandlerFunc) http.HandlerFunc
}

type authRequest struct {
	tv token.Validator
}

// NewAuthRequest new authed request
func NewAuthRequest(tv token.Validator) AuthRequest {
	return &authRequest{
		tv: tv,
	}
}

// Auth is the main entry point into our application
// Step 1: Authenticate token [required]
// Step 2: If good, call the next middleware in the chain
func (service *authRequest) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1. authenticate
		ctx, err := service.authenticate(w, r)

		// authentication failed
		if err != nil {
			return
		}

		// Step 2. success, add context to the request and continue down chain
		next(w, r.WithContext(ctx))
		return
	}
}

// Authenticate token
func (service *authRequest) authenticate(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	// get the jwt token from the request
	bearer := r.Header.Get("Authorization")
	// no authorization header sent
	if bearer == "" {
		log.Printf("Authorization is required")
		http.Error(w, "no authorization", http.StatusUnauthorized)
		return nil, errors.New("")
	}
	splitToken := strings.Split(bearer, "Bearer ")
	jwt := splitToken[1]

	// confirm the token is sent in the request
	if len(jwt) == 0 {
		log.Printf("JWT not found")
		http.Error(w, "token not found", http.StatusUnauthorized)
		return r.Context(), errors.New("")
	}

	// validate the token and get the session
	var session *token.Session
	session, err := service.tv.Validate(jwt)

	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return r.Context(), errors.New("")
	}

	// store the session in the config
	ctx := context.WithValue(r.Context(), Session, session)
	// store the jwt in the context for service to service communication
	ctx = context.WithValue(ctx, Token, jwt)

	return ctx, nil
}
