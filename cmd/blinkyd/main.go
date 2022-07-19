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

	"github.com/BrenekH/blinky/apiunstable"
	"github.com/BrenekH/blinky/cmd/blinkyd/viperutils"
	"github.com/BrenekH/blinky/httpbasicauth"
	"github.com/BrenekH/blinky/keyvaluestore"
	"github.com/BrenekH/blinky/vars"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	// Print out the version when requested
	for _, v := range os.Args {
		switch strings.ToLower(v) {
		case "--version", "-v":
			fmt.Printf("blinkyd version %s\n", vars.Version)
			os.Exit(0)
		}
	}

	// TODO: Print out a custom help message that better explains blinkyd's usage

	if err := viperutils.Setup(); err != nil {
		panic(err)
	}

	repoPath := viper.GetString("RepoPath")
	dbPath := viper.GetString("ConfigDir") + "/kv-db"
	requireSignedPkgs := viper.GetBool("RequireSignedPkgs")
	gpgDir := viper.GetString("GPGDir")
	signingKey := viper.GetString("SigningKeyFile")
	httpPort := viper.GetString("HTTPPort")
	apiUname := viper.GetString("APIUsername")
	apiPasswd := viper.GetString("APIPassword")

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

	for _, repoPath := range repoPaths {
		if err := os.MkdirAll(repoPath+"/x86_64", 0777); err != nil {
			log.Printf("WARNING: Unable to create %s because of the following error: %v", repoPath+"/x86_64", err)
		}

		// Build and run repo-add command, including the --sign arg if requested
		repoAddArgs := []string{"-q", "-R", "--nocolor"}
		if signDB {
			repoAddArgs = append(repoAddArgs, "--sign")
		}
		repoAddArgs = append(repoAddArgs, repoPath+"/x86_64/"+filepath.Base(repoPath)+".db.tar.gz")

		cmd := exec.Command("repo-add", repoAddArgs...)
		cmd.Env = append(cmd.Env, fmt.Sprintf("GNUPGHOME=%s", gpgDir))
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Printf("WARNING: Unable to create repository database because of following error: could not run %s, received %s: %v", cmd.String(), string(out), err)
		}
	}

	registerHTTPHandlers(repoPaths, dbPath, gpgDir, apiUname, apiPasswd, requireSignedPkgs, signDB)

	fmt.Printf("Blinky is now listening for connections on port %s\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil)

	// This may or may not ever be reached :shrug:
	if signDB {
		if err := os.RemoveAll(gpgDir); err != nil {
			panic(err)
		}
	}
}

func registerHTTPHandlers(repoPaths []string, dbPath, gpgDir, apiUname, apiPasswd string, requireSignedPkgs, signDB bool) {
	registerRepoPaths("/repo", repoPaths)

	ds, err := keyvaluestore.New(dbPath)
	if err != nil {
		panic(err)
	}

	apiAuth := httpbasicauth.New(apiUname, apiPasswd)

	apiRouter := mux.NewRouter()
	apiUnstable := apiunstable.New(&ds, &apiAuth, correlateRepoNames(repoPaths), gpgDir, requireSignedPkgs, signDB)
	apiUnstable.Register(apiRouter.PathPrefix("/api/unstable/").Subrouter())

	http.Handle("/api/unstable/", apiRouter)

	http.Handle("/", http.HandlerFunc(indexPageHandler))
}

func registerRepoPaths(base string, repoPaths []string) {
	for _, path := range repoPaths {
		repoName := filepath.Base(path)
		repoNameSlashed := "/" + repoName + "/"
		http.Handle(base+repoNameSlashed, logRequestMiddleware(http.StripPrefix(base+repoNameSlashed, http.FileServer(http.Dir(path)))))
	}
}

// logRequestMiddleware wraps a http.Handler and logs the request path if
// the environment variable BLINKY_LOG_LEVEL is set to 'debug'.
func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(os.Getenv("BLINKY_LOG_LEVEL")) == "debug" {
			log.Printf("method=%v path=%v\n", r.Method, r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}

func correlateRepoNames(repoPaths []string) map[string]string {
	m := make(map[string]string)

	for _, path := range repoPaths {
		m[filepath.Base(path)] = path
	}

	return m
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(`<html>
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Blinky - Pacman Repo Server</title>

		<style>
			body {
				text-align: center;
				font-family: Arial;
			}

			a:visited {
				color: blue;
			}
		</style>
	</head>
	<body>
		<h1 style="margin-bottom: 0.25rem;">Blinky</h1>
		<p style="margin-top: 0.25rem;">Simple, all in one Pacman repository hosting server software.</p>
		<hr>
		<a href="https://github.com/BrenekH/blinky#README" target="_blank" rel="noopener noreferrer">GitHub Project Link</a>
	</body>
</html>`))
}
