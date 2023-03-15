package homebridge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	token  string
	opts   *ClientOpts
}

type ClientOpts struct {
	Host     string
	User     string
	Password string
	Debug    bool
}

func NewClient(httpClient *http.Client, opts *ClientOpts) *Client {
	return &Client{
		client: httpClient,
		opts:   opts,
	}
}

func (api *Client) GetAccessory(ctx context.Context, accessoryId string) ([]byte, error) {
	path := fmt.Sprintf("%saccessories/%s", api.opts.Host, accessoryId)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	token, err := api.getToken(ctx)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")
	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if api.opts.Debug {
		fmt.Println(string(responseBody))
	}
	switch response.StatusCode {
	case 200:
		return responseBody, nil
	default:
		return responseBody, fmt.Errorf("invalid http status: %d", response.StatusCode)
	}

}

func (api *Client) tokenIsValid() (bool, error) {
	claims := jwt.MapClaims{}
	token, _, err := jwt.NewParser().ParseUnverified(api.token, claims)
	if err != nil {
		return false, err
	}
	_ = token
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return false, err
	}
	return expirationTime.Before(time.Now()), nil
}

func (api *Client) getToken(ctx context.Context) (string, error) {
	tokenIsValid, err := api.tokenIsValid()
	if err != nil {
		return "", err
	}
	if tokenIsValid {
		return api.token, nil
	}

	reqBody := TokenRequest{api.opts.User, api.opts.Password}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := api.opts.Host + "auth/login"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := api.client.Do(request)
	if err != nil {
		return "", err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if api.opts.Debug {
		fmt.Println(string(responseBody))
	}
	switch response.StatusCode {
	case 201:
		r := TokenResponse{}
		err := json.Unmarshal(responseBody, &r)
		if err != nil {
			return "", err
		}
		api.token = r.AccessToken
		return api.token, nil

	default:
		return "", fmt.Errorf("invalid status code %d", response.StatusCode)
	}

}

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
