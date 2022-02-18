package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrenekH/blinky/apiunstable"
)

type handleFunc func(pattern string, handler http.Handler)

func main() {
	repoPaths, ok := os.LookupEnv("BLINKY_REPO_PATH")
	if !ok {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		repoPaths = cwd + "/repo"
	}

	registerRepoPaths(http.Handle, "/repo", strings.Split(repoPaths, ":"))

	apiunstable.Register()

	http.ListenAndServe(":9000", nil)
}

func registerRepoPaths(h handleFunc, base string, repoPaths []string) {
	for _, path := range repoPaths {
		repoName := filepath.Base(path)
		repoNameSlashed := "/" + repoName + "/"
		h(base+repoNameSlashed, http.StripPrefix(base+repoNameSlashed, http.FileServer(http.Dir(path))))
	}
}
