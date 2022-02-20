package apiunstable

import (
	"net/http"

	"github.com/BrenekH/blinky"
	"github.com/gorilla/mux"
)

func New(storageProvider *blinky.PackageNameToFileProvider, foundRepos []string) API {
	return API{storage: storageProvider, repos: foundRepos}
}

type API struct {
	storage *blinky.PackageNameToFileProvider
	repos   []string
}

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
	_ = targetPkgName

	if !a.isValidRepo(targetRepo) {
		w.Write([]byte("Invalid repository"))
		w.WriteHeader(http.StatusBadRequest)
	}

	// TODO: Implement

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *API) deleteRepoPkg(w http.ResponseWriter, r *http.Request) {
	// Extract variables from URL
	vars := mux.Vars(r)
	targetRepo := vars["repo"]
	targetPkgName := vars["package_name"]
	_ = targetPkgName

	if !a.isValidRepo(targetRepo) {
		w.Write([]byte("Invalid repository"))
		w.WriteHeader(http.StatusBadRequest)
	}

	// TODO: Implement

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *API) isValidRepo(r string) bool {
	for _, repo := range a.repos {
		if r == repo {
			return true
		}
	}

	return false
}
