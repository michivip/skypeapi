/*
MIT License

Copyright (c) 2017 MichiVIP

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */
package skypeapi

import (
	"net/http"
	"net/url"
	"encoding/json"
)

type MessageListener interface {
	handle()
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
