package skypeapi

import (
	"net/http"
	"crypto/tls"
	"encoding/json"
)

const (
	defaultPath           string = "/"
	defaultTlsHeaderValue        = "max-age=63072000; includeSubDomains" // max-age in seconds which matches 2 years
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
	endpoint := &Endpoint{
		Address:   address,
		Path:      defaultPath,
		TLSConfig: cfg,
	}
	return endpoint
}

type EndpointHandler struct {
	// The header value which will be sent to the client with the "Strict-Transport-Security" key
	TlsHeaderValue string
	// The function to handle incoming decoded Activity object
	ActivityReceivedHandleFunction func(activity *Activity)
}

// The activityReceivedHandleFunction will gets called on incoming Activity objects for example incoming skype messages
// Returns a new Endpoint struct object with the default Strict-Transport-Security Header "max-age=63072000; includeSubDomains".
func NewEndpointHandler(activityReceivedHandleFunction func(activity *Activity)) (*EndpointHandler) {
	endpointHandler := &EndpointHandler{
		TlsHeaderValue:                 defaultTlsHeaderValue,
		ActivityReceivedHandleFunction: activityReceivedHandleFunction,
	}
	return endpointHandler
}

// Internal method to hook skype actions.
func (endpointHandler EndpointHandler) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) {
	if len(endpointHandler.TlsHeaderValue) != 0 {
		responseWriter.Header().Add("Strict-Transport-Security", endpointHandler.TlsHeaderValue)
	}

	var activity Activity
	if err := json.NewDecoder(req.Body).Decode(&activity); err == nil {
		responseWriter.WriteHeader(http.StatusOK)
		endpointHandler.ActivityReceivedHandleFunction(&activity)
	} else {
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
}

// This method could be used on an Endpoint struct object to setup an own web server which
// handles skype actions. The returned http.Server can still be edited to
// TODO: editable tls configuration
func (endpoint Endpoint) SetupServer(handler EndpointHandler) (http.Server) {
	mux := http.NewServeMux()
	mux.Handle(endpoint.Path, handler)
	srv := &http.Server{
		Addr:         endpoint.Address,
		Handler:      mux,
		TLSConfig:    endpoint.TLSConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	return *srv
}
