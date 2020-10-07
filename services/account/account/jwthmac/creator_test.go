package jwthmac_test

import (
	"fmt"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account/jwthmac"
	"github.com/stretchr/testify/assert"
)

// TestCreatorTestTable is a test table for the test TestCreator
var TestCreatorTestTable = []struct {
	secret     string
	issuer     string
	expiration time.Duration
	account    account.Account
	profID     string
}{
	{
		"?&M@4NmDM7#dQnUW#T?QfWc2",
		"jN_vch99wz%k56GKEPD%JP_G",
		25 * time.Minute,
		account.Account{
			ID:    "1",
			Email: "Z7Ryu28G57F_dT6hUa?-C_@8@g6rVsu7nWqHypnTNP8=?eTG!.com",
		},
		"1",
	},
	{
		"%jES39$YUZ&Z-x5y@Fu#!h7H2dnK@F86",
		"2$kahsARpLmKzHqHWEJWPb79?mRE?@^E",
		25 * time.Minute,
		account.Account{
			ID:    "2",
			Email: "QH5PPn@CyEn^cCaZ5&hc!jTNr@BG%GJr@=%Bpw8wGG8yMJW75==x9^!=%MMXWZB_L.com",
		},
		"2",
	},
	{
		"fR9Hr&&sPFpc4nFC?C6ZgxM8dr3DkB?TpuVBbLt=3?ywq7vy",
		"A_#%UWd3m5_V6cX3e^jp5@er@q2pH@N9JySQ-NjxnQnBzR$q",
		25 * time.Minute,
		account.Account{
			ID:    "3",
			Email: "u$&p#YK!H7c9kNf_cX=^w%!3W?VcG53n32Z4rB=YTU#%udQd@&nQ5_7V?s2vArCsswZ^ZDaznmFRnMNVN8A2@sgy&5+9aXHG9.com",
		},
		"3",
	},
	{
		"mJmbM5=Mhae?9JGuR#VXM43LR%DbCf4qqhTP+v=JFUV?_nv-bUQr#GuUAJ6s+w-F",
		"+%6#her2seba_unmA$b+yT+W6Zjcj%myD7s_4xpknA+h83sszk9aPHtm9776$-du",
		25 * time.Minute,
		account.Account{
			ID:    "4",
			Email: "#n@&dtyWUyFuF+tUJBH7cmc54NzUt$Uddxn-e%8YHPK8*y@q2qHB&Zwf6vQg6gqR@UE5@4HmUd!#qezzEhhs+Ptr4Sh2Hjs8z9u^5m$HWe^Ld5GB%H*JC?aUvY7LMd7FZ.com",
		},
		"4",
	},
}

// TestCreator tests that our jwthmac creator creates valid JWTs, and are parsed and validated correctly.
// TODO: do testing on the time values returned to confirm expiration and issued at are in some acceptable range.
func TestCreator(t *testing.T) {
	for _, tt := range TestCreatorTestTable {
		c := jwthmac.NewCreator(tt.secret, tt.issuer, tt.expiration)
		token, err := c.Create(&tt.account, tt.profID)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		// parse and validate token. will receive secret via closure
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("received token with wrong signing method %v", token.Header["alg"])
			}
			return []byte(tt.secret), nil
		}

		parsedToken, err := jwt.ParseWithClaims(token, &jwthmac.CustomClaim{}, keyFunc)
		if err != nil {
			t.Fatalf("failed to validate and parse token: %v", err)
		}

		assert.True(t, parsedToken.Valid, "token is not valid")
		assert.Nilf(t, parsedToken.Claims.Valid(), "claims are not valid: %v", parsedToken.Claims.Valid())

		var claims *jwthmac.CustomClaim
		var ok bool
		if claims, ok = parsedToken.Claims.(*jwthmac.CustomClaim); !ok {
			t.Fatal("type assertion to jwthmac.CustomClaim failed")
		}

		assert.Equal(t, tt.account.ID, claims.Subject)
		assert.Equal(t, tt.issuer, claims.Issuer)
		assert.NotEmpty(t, claims.ExpiresAt)
		assert.NotEmpty(t, claims.Id)
		assert.NotEmpty(t, claims.IssuedAt)
	}
}
