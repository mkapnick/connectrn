package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client gcs client
type Client interface {
	GetAccountByEmail(jwt string, email string) (*Account, error)
	SignUp(sc *SignupCredentials) (*Account, error)
}

// Client is an http client
type client struct {
	c          *http.Client
	accountURL *url.URL
}

// NewClient is a constructor for our client
func NewClient(accountURL *url.URL) Client {
	c := &client{
		c:          &http.Client{},
		accountURL: accountURL,
	}

	return c
}

func (c *client) SignUp(sc *SignupCredentials) (*Account, error) {
	b, err := json.Marshal(sc)
	if err != nil {
		return nil, fmt.Errorf("SignUp: failed to marshal signup creds: %s", err)
	}

	u := fmt.Sprintf("%s%s", c.accountURL.String(), "signup/")
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("SignUp: failed to create request: %s", err)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// look for `200` or `201`
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("SignUp: Code: %d Status: %s", resp.StatusCode, resp.Status)
	}

	var a *Account
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, fmt.Errorf("SignUp: failed to decode resp body: %s", resp.Body)
	}

	return a, nil
}

// GetAccountByEmail send notification request
func (c *client) GetAccountByEmail(jwt string, email string) (*Account, error) {
	req, err := http.NewRequest("GET", c.accountURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Account: failed to create request: %s", err)
	}

	// pass in email as a query param ?email=<email>
	q := req.URL.Query()
	q.Add("email", email)
	req.URL.RawQuery = q.Encode()

	// set jwt header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// look for `200` or `201`
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Account: Code: %d Status: %s", resp.StatusCode, resp.Status)
	}

	var acc *Account
	err = json.NewDecoder(resp.Body).Decode(&acc)
	if err != nil {
		return nil, fmt.Errorf("GetAccountProfile: failed to decode resp body: %s", resp.Body)
	}

	return acc, nil
}
