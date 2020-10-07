package jwthmac

import (
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/michaelk99/connectrn/internal/token"
)

// Claims the claims coming in from the `jwt`
type claims struct {
	// standard claims
	jwt.StandardClaims
	// session should always be encoded within the jwt
	token.Session
}

// validator is an implementation of the TokenValidator interface
type validator struct {
	secret []byte
	method string
}

// NewTokenStore create a new token validator
func NewTokenStore(secret []byte, method string) token.Validator {
	return &validator{
		secret: secret,
		method: method,
	}
}

// Validate an incoming jwt
// example token
// "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJp
// c3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
// Docs: https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
func (v validator) Validate(tokenString string) (*token.Session, error) {
	claims := &claims{}
	tok, err := jwt.ParseWithClaims(tokenString, claims, func(tok *jwt.Token) (interface{}, error) {
		return v.secret, nil
	})

	if err != nil {
		log.Printf("Error parsing jwt: %s", err)
		return nil, err
	}

	if !tok.Valid {
		log.Printf("Token is considered invalid")
		return nil, fmt.Errorf("Token is considered invalid")
	}

	// fill out the session to be used for the lifetime of the request
	session := token.Session{}
	session.Email = claims.Email
	session.FirstName = claims.FirstName
	session.LastName = claims.LastName
	session.AccountID = claims.AccountID
	session.ProfileID = claims.ProfileID
	session.CompanyID = claims.CompanyID
	session.ClubID = claims.ClubID
	session.AccountRoles = claims.AccountRoles

	return &session, nil
}
