package apiunstable

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Register registers http handlers associated with the unstable API.
func Register() {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/unstable/").Subrouter()

	s.HandleFunc("/{repo}/package/{package_name}", handleRepoPkg)

	http.Handle("/", r)
}

// Routes (all prefixed with /api/unstable):
// PUT /:repo/package/:package_name
// DELETE /:repo/package/:package_name

func handleRepoPkg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("repo: ", vars["repo"])
	fmt.Println("package: ", vars["package_name"])
	w.Write([]byte("Hello"))
}
