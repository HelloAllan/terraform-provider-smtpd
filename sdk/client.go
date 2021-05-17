package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default SMTPD API URL
const HostURL string = "https://api.smtpd.dev"
const APIVersion string = "v1"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	APIVersion string
	Token      string
}

// AuthStruct -
type AuthStruct struct {
	APIKey    string `json:"username"`
	APISecret string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	Expire       int    `json:"expires_in"`
	Scope        []int  `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// NewClient -
func NewClient(host, key, secret *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default SMPTD API URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if (key != nil) && (secret != nil) {
		// form request body
		data, err := json.Marshal(AuthStruct{
			APIKey:    *key,
			APISecret: *secret,
		})
		if err != nil {
			return nil, err
		}

		// authenticate
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/token?grant_type=password", c.HostURL), strings.NewReader(string(data)))
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)

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

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.Token)

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

// CreateProfile creates a profile
func (c *Client) CreateProfile(ctx context.Context, profile *Profile) (*Profile, error) {
	data, err := profile.ConvertToJSON()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/%s/email/profile", c.HostURL, c.APIVersion), strings.NewReader(string(data)))
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

// GetProfile retrieves a server
func (c *Client) GetProfile(ctx context.Context, id string) (*Profile, error) {

	profile := &Profile{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%s/email/profile/%s/setup", c.HostURL, c.APIVersion, id), strings.NewReader(string(data)))
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
