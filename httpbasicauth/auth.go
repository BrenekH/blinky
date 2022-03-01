package httpbasicauth

import (
	"encoding/base64"
	"net/http"
)

func New(username, password string) BasicAuth {
	return BasicAuth{
		validKey: "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
	}
}

// BasicAuth is a basic implementation of the Authenticator interface.
// It expects requests to be sent with an Authorization header containing
// the HTTP Basic Authentication standard.
type BasicAuth struct { // implements: github.com/BrenekH/blinky.Authenticator
	validKey string
}

func (a *BasicAuth) CreateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != a.validKey {
			http.Error(rw, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(rw, r)
	})
}
