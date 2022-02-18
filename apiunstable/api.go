package apiunstable

import "net/http"

func Register() {
	http.HandleFunc("/api/unstable/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}
