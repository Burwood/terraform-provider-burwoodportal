package burwoodportal

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Portal URL
const HostURL string = "http://localhost:5000"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Token string `json:"token"`
}

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default burwood portal URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

// SignIn - Get a new token for user
func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/token", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	authString := fmt.Sprintf("%s:%s", c.Auth.Username, c.Auth.Password)
	encodedAuthString := b64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuthString))

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &ar, nil
}

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := c.Token

	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("x-access-token", token)
	req.Header.Set("content-type", "application/json")

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

// Reusable function to make a GET request on an API endpoint.
func (c *Client) getEndpointList(endpoint string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)

	// Unmarshal response JSON into a map data structure
	bodyMap := make([]map[string]interface{}, 0)
	err = json.Unmarshal(body, &bodyMap)
	if bodyMap == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return bodyMap, nil
}

// Reusable function to make a GET request on an API endpoint.
func (c *Client) getEndpointSingleItem(endpoint string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)

	// Unmarshal response JSON into a map data structure
	bodyMap := make(map[string]interface{}, 0)
	err = json.Unmarshal(body, &bodyMap)
	if bodyMap == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return bodyMap, nil
}


// Reusable function to make a GET request on an API endpoint.
func (c *Client) postEndpoint(endpoint string, postBody map[string]interface{}) ([]map[string]interface{}, error) {
	postBodyMarshaled, err := json.Marshal(postBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, endpoint), strings.NewReader(string(postBodyMarshaled)))
	if err != nil {
		return nil, err
	}

	responseBody, err := c.doRequest(req, nil)

	// Unmarshal response JSON into a map data structure
	responseBodyMap := make([]map[string]interface{}, 0)
	err = json.Unmarshal(responseBody, &responseBodyMap)
	if responseBodyMap == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return responseBodyMap, nil
}

func (c *Client) postGroups(endpoint string, postBody []Group) ([]Group, error) {
	postBodyMarshaled, err := json.Marshal(postBody)
	if err != nil {
		return nil, err
	}

	processedBody := strings.NewReader(string(postBodyMarshaled))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, endpoint), processedBody)
	if err != nil {
		return nil, err
	}
	

	responseBody, err := c.doRequest(req, nil)

	// Unmarshal response JSON into a map data structure
	responseBodyMap := make([]Group, 0)
	err = json.Unmarshal(responseBody, &responseBodyMap)
	if responseBodyMap == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return responseBodyMap, nil
}


