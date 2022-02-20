package apiunstable

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BrenekH/blinky"
	"github.com/gorilla/mux"
)

func New(storageProvider blinky.PackageNameToFileProvider, foundRepos map[string]string) API {
	return API{storage: storageProvider, repos: foundRepos}
}

type API struct {
	storage blinky.PackageNameToFileProvider
	repos   map[string]string // Map of repo name to repo path
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid repository"))
		return
	}

	// Populate r.FormFile by parsing the request body and loading up to 256 MB of the files into memory. The rest of the files are stored on disk.
	r.ParseMultipartForm(256_000_000)

	formPkgFile, formPkgHeader, err := r.FormFile("package")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	saveMultipartFile(formPkgFile, formPkgHeader, a.repos[targetRepo])

	formSigFile, _, err := r.FormFile("signature")
	if err != nil {
		log.Printf("Warning: No signature file ERROR: %v\n", err)
	} else {
		// Save the signature file, ensuring that it gets saved using the correct format no matter the filename sent.
		temp := formPkgHeader.Filename
		formPkgHeader.Filename = temp + ".sig"
		saveMultipartFile(formSigFile, formPkgHeader, a.repos[targetRepo])
		formPkgHeader.Filename = temp
	}

	cmd := exec.Command("repo-add", "-q", "-R", "--nocolor", a.repos[targetRepo]+"/x86_64/"+targetRepo+".db.tar.gz", a.repos[targetRepo]+"/x86_64/"+formPkgHeader.Filename)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ERROR running %s: %s", cmd.String(), string(out))
		http.Error(w, "Failed to add package to the database. Check the server logs for more information.", http.StatusInternalServerError)
		return
	}

	a.storage.StorePackageFile(fmt.Sprintf("%s/%s", targetRepo, targetPkgName), "")

	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteRepoPkg(w http.ResponseWriter, r *http.Request) {
	// Extract variables from URL
	vars := mux.Vars(r)
	targetRepo := vars["repo"]
	targetPkgName := vars["package_name"]
	_ = targetPkgName

	if !a.isValidRepo(targetRepo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid repository"))
		return
	}

	// TODO: Implement

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *API) isValidRepo(r string) bool {
	for repo := range a.repos {
		if r == repo {
			return true
		}
	}

	return false
}

func saveMultipartFile(mFile multipart.File, header *multipart.FileHeader, repoPath string) error {
	dest := repoPath + "/x86_64/" + filepath.Base(filepath.Clean(header.Filename))

	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, mFile); err != nil {
		return err
	}

	return nil
}
