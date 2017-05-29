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
	"crypto/tls"
	"encoding/json"
)

const (
	defaultPath            string = "/"
	defaultTlsHeaderValue  string = "max-age=63072000; includeSubDomains" // max-age in seconds which matches 2 years
	authorizationHeaderKey string = "Authorization"
)

type Endpoint struct {
	// Explanation: The address the server should listen on. This declares the port and the ip.
	// Example: ":2345" The application would run on 0.0.0.0 with the port 2345
	Address string
	// Explanation: The path which the server receives its requests.
	// Example: If the value is set to "/skype/" than the server would listen on "https://domain.tld/skype/"
	Path string
	// Explanation: The TLSConfig which declares which values will be sent to a client
	TLSConfig *tls.Config
}

// Returns a new Endpoint struct object with the default request path "/".
func NewEndpoint(address string) (*Endpoint) {
	// the default TLS config
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	return &Endpoint{
		Address:   address,
		Path:      defaultPath,
		TLSConfig: cfg,
	}
}

type EndpointHandler struct {
	// The MicrosoftAppId is used to authorize incoming requests
	MicrosoftAppId string
	// The authorization token which is used to authorize incoming requests
	AuthorizationToken string
	// The header value which will be sent to the client with the "Strict-Transport-Security" key
	TlsHeaderValue string
	// The function to handle incoming decoded Activity object
	ActivityReceivedHandleFunction func(activity *Activity)
}

// The activityReceivedHandleFunction will gets called on incoming Activity objects for example incoming skype messages.
// The authorization token which is used to authenticate incoming requests by the microsoft servers.
// The microsoftAppId which is used to authorize incoming requests
// Returns a new Endpoint struct object with the default Strict-Transport-Security Header "max-age=63072000; includeSubDomains".
func NewEndpointHandler(activityReceivedHandleFunction func(activity *Activity), authorizationToken, microsoftAppId string) (*EndpointHandler) {
	endpointHandler := &EndpointHandler{
		AuthorizationToken:             authorizationToken,
		TlsHeaderValue:                 defaultTlsHeaderValue,
		ActivityReceivedHandleFunction: activityReceivedHandleFunction,
		MicrosoftAppId:                 microsoftAppId,
	}
	return endpointHandler
}

// This method does not cache the SigningKeys
// The req which should be proved
func (endpointHandler EndpointHandler) IsAuthorized(req *http.Request) bool {
	signingKeys, err := GetSigningKeys()
	if err != nil {
		return false
	} else {
		return endpointHandler.IsAuthorizedWithSigningKeys(req, signingKeys)
	}
}

// The req which should be proved
// The SigningKeys which can be used to authorize the request
func (endpointHandler EndpointHandler) IsAuthorizedWithSigningKeys(req *http.Request, signingKeys SigningKeys) bool {
	var authorizationValue string = req.Header.Get(authorizationHeaderKey)
	if microsoftJsonWebToken, err := ParseMicrosoftJsonWebToken(authorizationValue);
		err != nil {
		return false
	} else {
		return microsoftJsonWebToken.Verify(endpointHandler.MicrosoftAppId, signingKeys)
	}
}

// Internal method to hook skype actions.
func (endpointHandler EndpointHandler) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) {
	if len(endpointHandler.TlsHeaderValue) != 0 {
		responseWriter.Header().Add("Strict-Transport-Security", endpointHandler.TlsHeaderValue)
	}

	var activity Activity
	if !endpointHandler.IsAuthorized(req) {
		responseWriter.WriteHeader(http.StatusForbidden)
	} else if err := json.NewDecoder(req.Body).Decode(&activity); err == nil {
		responseWriter.WriteHeader(http.StatusOK)
		endpointHandler.ActivityReceivedHandleFunction(&activity)
	} else {
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
}

// This method could be used on an Endpoint struct object to setup an own web server which
// handles skype actions. The returned http.Server can still be edited to
func (endpoint Endpoint) SetupServer(handler EndpointHandler) (*http.Server) {
	mux := http.NewServeMux()
	mux.Handle(endpoint.Path, handler)
	srv := &http.Server{
		Addr:         endpoint.Address,
		Handler:      mux,
		TLSConfig:    endpoint.TLSConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	return srv
}
