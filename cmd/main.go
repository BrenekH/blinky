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

	ds, err := jsonds.New("/var/lib/blinky/packageAssociations.json") // TODO: Allow user to override with env vars
	if err != nil {
		panic(err)
	}

	apiUnstable := apiunstable.New(&ds, correlateRepoNames(repoPaths))
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

func correlateRepoNames(repoPaths string) map[string]string {
	m := make(map[string]string)

	split := strings.Split(repoPaths, ":")

	for _, path := range split {
		m[filepath.Base(path)] = path
	}

	return m
}
