package apikeyauth

import "net/http"

func New(validKey string) APIKeyAuth {
	return APIKeyAuth{
		validKey: validKey,
	}
}

type APIKeyAuth struct { // implements: github.com/BrenekH/blinky.Authenticator
	validKey string
}

func (a *APIKeyAuth) CreateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != a.validKey {
			http.Error(rw, "Invalid API key", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(rw, r)
	})
}
