package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrenekH/blinky/apiunstable"
	"github.com/BrenekH/blinky/cmd/viperutils"
	"github.com/BrenekH/blinky/jsonds"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	if err := viperutils.Setup(); err != nil {
		panic(err)
	}

	repoPath := viper.GetString("RepoPath")
	jsonDBPath := viper.GetString("ConfigDir") + "/packageAssociations.json"
	requireSignedPkgs := viper.GetBool("RequireSignedPkgs")
	gpgDir := viper.GetString("GPGDir")
	httpPort := viper.GetString("HTTPPort")

	fmt.Printf("Configuration: %+v\n", viper.AllSettings())

	repoPaths := strings.Split(repoPath, ":")

	for _, v := range repoPaths {
		if err := os.MkdirAll(v+"/x86_64", 0777); err != nil {
			log.Printf("WARNING: Unable to create %s because of the following error: %v", v+"/x86_64", err)
		}
	}

	rootRouter := mux.NewRouter()

	// The PathPrefix value and base string must be the same so that the file server can properly serve the files.
	registerRepoPaths(rootRouter.PathPrefix("/repo").Subrouter(), "/repo", repoPaths)

	ds, err := jsonds.New(jsonDBPath) // TODO: Allow user to override with env vars
	if err != nil {
		panic(err)
	}

	apiUnstable := apiunstable.New(&ds, correlateRepoNames(repoPaths), gpgDir, requireSignedPkgs, false)
	apiUnstable.Register(rootRouter.PathPrefix("/api/unstable/").Subrouter())

	http.Handle("/", rootRouter)

	fmt.Printf("Blinky is now listening for connections on port %s\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil)
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
