package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrenekH/blinky/apiunstable"
	"github.com/BrenekH/blinky/jsonds"
	"github.com/gorilla/mux"
)

func main() {
	repoPathStr, ok := os.LookupEnv("BLINKY_REPO_PATH")
	if !ok {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		repoPathStr = cwd + "/repo"
	}

	repoPaths := strings.Split(repoPathStr, ":")

	for _, v := range repoPaths {
		if err := os.MkdirAll(v+"/x86_64", 0777); err != nil {
			log.Printf("WARNING: Unable to create %s because of the following error: %v", v+"/x86_64", err)
		}
	}

	rootRouter := mux.NewRouter()

	// The PathPrefix value and base string must be the same so that the file server can properly serve the files.
	registerRepoPaths(rootRouter.PathPrefix("/repo").Subrouter(), "/repo", repoPaths)

	ds, err := jsonds.New("/var/lib/blinky/packageAssociations.json") // TODO: Allow user to override with env vars
	if err != nil {
		panic(err)
	}

	requireSignedPkgs := false
	if strings.ToLower(os.Getenv("BLINKY_REQUIRE_SIGNED_PKGS")) == "true" {
		requireSignedPkgs = true
	}

	apiUnstable := apiunstable.New(&ds, correlateRepoNames(repoPaths), requireSignedPkgs)
	apiUnstable.Register(rootRouter.PathPrefix("/api/unstable/").Subrouter())

	http.Handle("/", rootRouter)

	fmt.Println("Blinky is now listening for connections on port 9000")
	http.ListenAndServe(":9000", nil)
}

func registerRepoPaths(router *mux.Router, base string, repoPaths []string) {
	for _, path := range repoPaths {
		repoName := filepath.Base(path)
		repoNameSlashed := "/" + repoName + "/"
		router.Handle(repoNameSlashed, http.StripPrefix(base+repoNameSlashed, http.FileServer(http.Dir(path))))
	}
}

func correlateRepoNames(repoPaths []string) map[string]string {
	m := make(map[string]string)

	for _, path := range repoPaths {
		m[filepath.Base(path)] = path
	}

	return m
}
