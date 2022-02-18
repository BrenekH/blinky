package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	http.HandleFunc("/api/", api)

	http.ListenAndServe(":9000", nil)
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func registerRepoPaths(h handleFunc, base string, repoPaths []string) {
	for _, path := range repoPaths {
		repoName := filepath.Base(path)
		repoNameSlashed := "/" + repoName + "/"
		h(base+repoNameSlashed, http.StripPrefix(base+repoNameSlashed, http.FileServer(http.Dir(path))))
	}
}
