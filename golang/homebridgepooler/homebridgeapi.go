package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HomeBridgeAPI struct {
	client      http.Client
	token       string
	tokenExpiry time.Time
	host        string
	user        string
	password    string
	debug       bool
}

func (api *HomeBridgeAPI) getToken(ctx context.Context) (string, error) {
	if !time.Now().After(api.tokenExpiry) {
		return api.token, nil
	}
	reqBody := TokenRequest{api.user, api.password}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	url := api.host + "auth/login"
	fmt.Println(url)
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
	if api.debug {
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
		api.tokenExpiry = time.Now().Add(time.Duration(r.ExpiresIn) * time.Millisecond)

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
