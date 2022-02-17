package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	repoPath, ok := os.LookupEnv("BLINKY_REPO_PATH")
	if !ok {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		repoPath = cwd + "/repo"
	}

	repoName := filepath.Base(repoPath)
	repoNameSlashed := "/" + repoName + "/"

	http.Handle(repoNameSlashed, http.StripPrefix(repoNameSlashed, http.FileServer(http.Dir(repoPath))))
	http.HandleFunc("/api/", api)

	http.ListenAndServe(":9000", nil)
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
