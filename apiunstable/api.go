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

	s.HandleFunc("/{repo}/package/{package_name}", putRepoPkg).Methods(http.MethodPut)
	s.HandleFunc("/{repo}/package/{package_name}", deleteRepoPkg).Methods(http.MethodDelete)

	http.Handle("/", r) // I dislike handling the root path here just so we can use Gorilla Mux for URL path variables.
}

func putRepoPkg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("repo: ", vars["repo"])
	fmt.Println("package: ", vars["package_name"])

	// TODO: Implement

	w.WriteHeader(http.StatusNotImplemented)
}

func deleteRepoPkg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("repo: ", vars["repo"])
	fmt.Println("package: ", vars["package_name"])

	// TODO: Implement

	w.WriteHeader(http.StatusNotImplemented)
}
