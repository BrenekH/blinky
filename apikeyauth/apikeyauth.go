package apikeyauth

import "net/http"

func New(validKey string) APIKeyAuth {
	return APIKeyAuth{
		validKey: validKey,
	}
}

// APIKeyAuth is a basic implementation of the Authenticator interface.
// It expects requests to be sent with an Authorization header containing the key unless the key
// is an empty string, then no auth is required from the client.
type APIKeyAuth struct { // implements: github.com/BrenekH/blinky.Authenticator
	validKey string
}

func (a *APIKeyAuth) CreateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if a.validKey != "" && r.Header.Get("Authorization") != a.validKey {
			http.Error(rw, "Invalid API key", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(rw, r)
	})
}
