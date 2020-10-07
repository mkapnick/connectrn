package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Client profile client
type Client interface {
	CreateProfile(jwt string, prof *Profile) (*Profile, error)
	GetProfileByAccountID(jwt string, accountID string) (*Profile, error)
	GetProfile(jwt string, profileID string) (*Profile, error)
}

// Client is an http client
type client struct {
	c       *http.Client
	profURL *url.URL
}

// NewClient is a constructor for our client
func NewClient(profURL *url.URL) Client {
	c := &client{
		c:       &http.Client{},
		profURL: profURL,
	}

	return c
}

// CreateProfile creates a profile [only used when creating an account]
func (c *client) CreateProfile(jwt string, prof *Profile) (*Profile, error) {
	b, err := json.Marshal(prof)
	if err != nil {
		return nil, fmt.Errorf("CreateProfile: failed to marshal profile: %s", err)
	}

	req, err := http.NewRequest("POST", c.profURL.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("CreateProfile: failed to create request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// look for `200` or `201`
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CreateProfile: Code: %d Status: %s", resp.StatusCode, resp.Status)
	}

	var p *Profile
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("CreateProfile: failed to decode resp body: %s", resp.Body)
	}

	return p, nil
}

// GetProfileByAccountID get profile by account id
func (c *client) GetProfileByAccountID(jwt string, accountID string) (*Profile, error) {
	req, err := http.NewRequest("GET", c.profURL.String(), nil)

	if err != nil {
		log.Printf("Error with new request %s", err)
		return nil, err
	}

	// pass in account_id as a query param ?account_id=<id>
	q := req.URL.Query()
	q.Add("account_id", accountID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	resp, err := c.c.Do(req)
	if err != nil {
		log.Printf("Error retrieving profile info %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GetProfile: Code: %d Status: %s", resp.StatusCode, resp.Status)
	}

	var prof *Profile
	err = json.NewDecoder(resp.Body).Decode(&prof)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: failed to decode resp body: %s", resp.Body)
	}

	return prof, nil
}

// GetProfile get profile by id
func (c *client) GetProfile(jwt string, profileID string) (*Profile, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", c.profURL.String(), profileID))
	if err != nil {
		fmt.Printf("Failed to parse url %s", err)
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Printf("Failed to create new request %s", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GetProfile: Code: %d Status: %s", resp.StatusCode, resp.Status)
	}

	var prof *Profile
	err = json.NewDecoder(resp.Body).Decode(&prof)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: failed to decode resp body: %s", resp.Body)
	}

	return prof, nil
}
