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
package skypeapiexamples

import (
	"net/http"
	"crypto/tls"
	"github.com/michivip/skypeapi"
	"fmt"
)

// some basic constants
const (
	actionHookPath     string = "/skype/actionhook"
	address                   = ":9443"
	someOtherStuffPath string = "/"
	// bad practice. In real production you should better request the token via skypeapi.RequestAccessToken
	authorizationBearerToken string = "YOUR-AUTH-TOKEN"
)

// this handles our skype activity
func handleActivity(activity *skypeapi.Activity) {
	if activity.Type == "message" {
		if err := skypeapi.SendReplyMessage(activity, "Good evening. Nice to meet you!", authorizationBearerToken);
			err != nil {
			panic(err)
		} else {
			fmt.Println("Successfully sent response message to skype user: " + activity.From.Name)
		}
	}
}

// our custom application handler function
func handleMainPath(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("This content is hilarious."))
}

func startCustomServerEndpoint() {
	mux := http.NewServeMux()
	// here we setup an own activity handler which listens to the path "/skype/actionhook"
	mux.Handle(actionHookPath, skypeapi.NewEndpointHandler(handleActivity))
	// here we could probably just handle our main application
	mux.HandleFunc(someOtherStuffPath, handleMainPath)
	// here you could provide your own TLS configuration
	customTlsConfig := &tls.Config{
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
	// custom server setup
	srv := &http.Server{
		Addr:         address,
		Handler:      mux,
		TLSConfig:    customTlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	// finally we just use the default method to start the server
	panic(srv.ListenAndServeTLS("keys/fullchain.pem", "keys/privkey.pem"))
}
