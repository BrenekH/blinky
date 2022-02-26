package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BrenekH/blinky/apikeyauth"
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

	fmt.Printf("Configuration: %+v\n", viper.AllSettings())

	repoPath := viper.GetString("RepoPath")
	jsonDBPath := viper.GetString("ConfigDir") + "/packageAssociations.json"
	requireSignedPkgs := viper.GetBool("RequireSignedPkgs")
	gpgDir := viper.GetString("GPGDir")
	signingKey := viper.GetString("SigningKeyFile")
	httpPort := viper.GetString("HTTPPort")
	apiKey := viper.GetString("APIKey")

	os.RemoveAll(gpgDir) // We don't care if this fails because of a missing dir, and if it's something else, we'll find out soon.

	var signDB bool
	if signingKey != "" {
		if _, err := os.Stat(signingKey); err == nil {
			signDB = true

			if err := os.MkdirAll(gpgDir, 0700); err != nil {
				panic(err)
			}

			cmd := exec.Command("gpg", "--allow-secret-key-import", "--import", signingKey)
			cmd.Env = append(cmd.Env, fmt.Sprintf("GNUPGHOME=%s", gpgDir))
			if b, err := cmd.CombinedOutput(); err != nil {
				log.Println(string(b))
				panic(err)
			}
		} else if errors.Is(err, os.ErrNotExist) {
			log.Printf("WARNING: The signing key %s does not exist\n", signingKey)
		}
	}

	repoPaths := strings.Split(repoPath, ":")

	for _, v := range repoPaths {
		if err := os.MkdirAll(v+"/x86_64", 0777); err != nil {
			log.Printf("WARNING: Unable to create %s because of the following error: %v", v+"/x86_64", err)
		}
	}

	registerHTTPHandlers(repoPaths, jsonDBPath, gpgDir, apiKey, requireSignedPkgs, signDB)

	fmt.Printf("Blinky is now listening for connections on port %s\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil)

	// This may or may not ever be reached :shrug:
	if signDB {
		if err := os.RemoveAll(gpgDir); err != nil {
			panic(err)
		}
	}
}

func registerHTTPHandlers(repoPaths []string, jsonDBPath, gpgDir, apiKey string, requireSignedPkgs, signDB bool) {
	rootRouter := mux.NewRouter()

	// The PathPrefix value and base string must be the same so that the file server can properly serve the files.
	registerRepoPaths(rootRouter.PathPrefix("/repo").Subrouter(), "/repo", repoPaths)

	ds, err := jsonds.New(jsonDBPath)
	if err != nil {
		panic(err)
	}

	apiAuth := apikeyauth.New(apiKey)

	apiUnstable := apiunstable.New(&ds, &apiAuth, correlateRepoNames(repoPaths), gpgDir, requireSignedPkgs, signDB)
	apiUnstable.Register(rootRouter.PathPrefix("/api/unstable/").Subrouter())

	http.Handle("/", rootRouter)

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
