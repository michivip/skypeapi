package skypeapi

import (
	"net/http"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type MessageListener interface {
	handle()
}

type RequestEndpoint struct{}

func (RequestEndpoint) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) {
	if bytes, err := ioutil.ReadAll(req.Body); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
	} else {
		fmt.Println(string(bytes))
	}
	responseWriter.WriteHeader(http.StatusOK)
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int `json:"expires_in"`
	ExtExpiresIn int `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

const (
	requestTokenUrl = "https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token"
)

func RequestAccessToken(microsoftAppId string, microsoftAppPassword string) (TokenResponse, error) {
	var tokenResponse TokenResponse
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", microsoftAppId)
	values.Set("client_secret", microsoftAppPassword)
	values.Set("scope", "https://api.botframework.com/.default")
	if response, err := http.PostForm(requestTokenUrl, values); err != nil {
		return tokenResponse, err
	} else {
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&tokenResponse)
		return tokenResponse, err
	}
}
