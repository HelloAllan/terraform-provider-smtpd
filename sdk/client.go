package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// APIURL - Default SMTPD API URL
const APIURL string = "https://api.smtpd.dev"
const APIVersion string = "v1"

// Client -
type Client struct {
	HTTPClient *http.Client
	Token      string
}

// AuthResponse -
type AuthResponse struct {
	TokenType    string   `json:"token_type"`
	AccessToken  string   `json:"access_token"`
	Expire       int      `json:"expires_at"`
	Scope        []string `json:"scope"`
	RefreshToken string   `json:"refresh_token"`
}

// NewClient -
func NewClient(host, key, secret, version *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	if (key != nil) && (secret != nil) {
		// authenticate
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/token?grant_type=password", APIURL), nil)
		if err != nil {
			return nil, err
		}

		body, err := c.doRequestBasicAuth(req, key, secret)
		if err != nil {
			return nil, err
		}

		// parse response body
		ar := AuthResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}

		c.Token = ar.AccessToken
	}

	return &c, nil
}

func (c *Client) doRequestBasicAuth(req *http.Request, key, secret *string) ([]byte, error) {
	req.SetBasicAuth(*key, *secret)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		return nil, err
	}

	return body, err
}

// CreateProfile creates a profile
func (c *Client) CreateProfile(profile *Profile) (*Profile, error) {
	data, err := profile.ConvertToJSON()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/%s/email/profile", APIURL, APIVersion), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	profile.LoadFromJSON([]byte(body))
	return profile, nil
}

// GetProfile retrieves a profile
func (c *Client) GetProfile(id string) (*Profile, error) {
	profile := &Profile{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%s/email/profile/%s/setup", APIURL, APIVersion, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	profile.LoadFromJSON([]byte(body))

	return profile, nil
}

// GetProfileByName retrieves a profile by name
func (c *Client) GetProfileByName(name string) (*Profile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%s/email/profile", APIURL, APIVersion), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	profiles := []Profile{}
	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}

	for _, profile := range profiles {
		if profile.ProfileName == name {
			return &profile, nil
		}
	}
	return nil, fmt.Errorf("Profile %s not found", name)
}

// NoOp retrieves a profile
func (c *Client) NoOp(id string) error {
	return nil
}
