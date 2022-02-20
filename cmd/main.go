package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrenekH/blinky/apiunstable"
	"github.com/BrenekH/blinky/jsonds"
	"github.com/gorilla/mux"
)

func main() {
	repoPaths, ok := os.LookupEnv("BLINKY_REPO_PATH")
	if !ok {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		repoPaths = cwd + "/repo"
	}

	rootRouter := mux.NewRouter()

	// The PathPrefix value and base string must be the same so that the file server can properly serve the files.
	registerRepoPaths(rootRouter.PathPrefix("/repo").Subrouter(), "/repo", strings.Split(repoPaths, ":"))

	ds, err := jsonds.New("./db.json")
	if err != nil {
		panic(err)
	}

	apiUnstable := apiunstable.New(&ds, extractRepoNames(repoPaths))
	apiUnstable.Register(rootRouter.PathPrefix("/api/unstable/").Subrouter())

	http.Handle("/", rootRouter)

	http.ListenAndServe(":9000", nil)
}

func registerRepoPaths(router *mux.Router, base string, repoPaths []string) {
	for _, path := range repoPaths {
		repoName := filepath.Base(path)
		repoNameSlashed := "/" + repoName + "/"
		router.Handle(repoNameSlashed, http.StripPrefix(base+repoNameSlashed, http.FileServer(http.Dir(path))))
	}
}

func extractRepoNames(repoPaths string) []string {
	s := []string{}

	split := strings.Split(repoPaths, ":")

	for _, path := range split {
		s = append(s, filepath.Base(path))
	}

	return s
}
