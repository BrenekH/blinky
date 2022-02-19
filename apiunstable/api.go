package apiunstable

import (
	"net/http"

	"github.com/gorilla/mux"
)

func New() API {
	// TODO: Require a database struct to be passed for use by the API.
	return API{}
}

type API struct{}

// Register registers http handlers associated with the unstable API.
func (a *API) Register(router *mux.Router) {
	router.HandleFunc("/{repo}/package/{package_name}", a.putRepoPkg).Methods(http.MethodPut)
	router.HandleFunc("/{repo}/package/{package_name}", a.deleteRepoPkg).Methods(http.MethodDelete)
}

func (a *API) putRepoPkg(w http.ResponseWriter, r *http.Request) {
	// Extract variables from URL
	vars := mux.Vars(r)
	targetRepo := vars["repo"]
	targetPkgName := vars["package_name"]
	_, _ = targetRepo, targetPkgName

	// TODO: Implement

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *API) deleteRepoPkg(w http.ResponseWriter, r *http.Request) {
	// Extract variables from URL
	vars := mux.Vars(r)
	targetRepo := vars["repo"]
	targetPkgName := vars["package_name"]
	_, _ = targetRepo, targetPkgName

	// TODO: Implement

	w.WriteHeader(http.StatusNotImplemented)
}
