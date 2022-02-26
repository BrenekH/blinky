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

func New(storageProvider blinky.PackageNameToFileProvider, authProvider blinky.Authenticator, foundRepos map[string]string, gnupgDir string, requireSignedPackages, useSignedDB bool) API {
	return API{
		gnupgDir:              gnupgDir,
		repos:                 foundRepos,
		useSignedDB:           useSignedDB,
		auth:                  authProvider,
		storage:               storageProvider,
		requireSignedPackages: requireSignedPackages,
	}
}

type API struct {
	storage               blinky.PackageNameToFileProvider
	auth                  blinky.Authenticator
	repos                 map[string]string // Map of repo name to repo path
	requireSignedPackages bool
	useSignedDB           bool
	gnupgDir              string // The location to set GNUPGHOME to, when repo-add/repo-remove.
}

// Register registers http handlers associated with the unstable API.
func (a *API) Register(router *mux.Router) {
	router.Handle("/{repo}/package/{package_name}", a.auth.CreateMiddleware(http.HandlerFunc(a.putRepoPkg))).Methods(http.MethodPut)
	router.Handle("/{repo}/package/{package_name}", a.auth.CreateMiddleware(http.HandlerFunc(a.deleteRepoPkg))).Methods(http.MethodDelete)
}

func (a *API) putRepoPkg(w http.ResponseWriter, r *http.Request) {
	// Extract variables from URL
	vars := mux.Vars(r)
	targetRepo := vars["repo"]
	targetPkgName := vars["package_name"]

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
		log.Printf("Reading FormFile: %v\n", err)
		return
	}

	formSigFile, _, err := r.FormFile("signature")
	if err != nil {
		if a.requireSignedPackages {
			http.Error(w, "A signature file is required to upload to this server", http.StatusBadRequest)
			return
		} else {
			log.Printf("Warning: No signature file ERROR: %v\n", err)
		}
	} else {
		// Save the signature file, ensuring that it gets saved using the correct format no matter the filename sent.
		temp := formPkgHeader.Filename
		formPkgHeader.Filename = temp + ".sig"
		saveMultipartFile(formSigFile, formPkgHeader, a.repos[targetRepo])
		formPkgHeader.Filename = temp
	}

	// This is after the signature file so that if the server requires a signed package, the file doesn't get copied
	// until the request is known to have a .sig file.
	saveMultipartFile(formPkgFile, formPkgHeader, a.repos[targetRepo])

	// Build and run repo-add command, including the --sign arg if requested
	args := []string{"-q", "-R", "--nocolor"}
	if a.useSignedDB {
		args = append(args, "--sign")
	}
	args = append(args, a.repos[targetRepo]+"/x86_64/"+targetRepo+".db.tar.gz", a.repos[targetRepo]+"/x86_64/"+formPkgHeader.Filename)

	cmd := exec.Command("repo-add", args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("GNUPGHOME=%s", a.gnupgDir))
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ERROR running %s: %s", cmd.String(), string(out))
		http.Error(w, "Failed to add package to the database. Check the server logs for more information.", http.StatusInternalServerError)
		return
	}

	if err := a.storage.StorePackageFile(fmt.Sprintf("%s/%s", targetRepo, targetPkgName), formPkgHeader.Filename); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Got error while saving new Blinky db entry: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteRepoPkg(w http.ResponseWriter, r *http.Request) {
	// Extract variables from URL
	vars := mux.Vars(r)
	targetRepo := vars["repo"]
	targetPkgName := vars["package_name"]

	if !a.isValidRepo(targetRepo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid repository"))
		return
	}

	// Build and run repo-remove command, including the --sign arg if requested
	args := []string{"-q", "--nocolor"}
	if a.useSignedDB {
		args = append(args, "--sign")
	}
	args = append(args, a.repos[targetRepo]+"/x86_64/"+targetRepo+".db.tar.gz", targetPkgName)

	cmd := exec.Command("repo-remove", args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("GNUPGHOME=%s", a.gnupgDir))

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ERROR running %s: %s", cmd.String(), string(out))
		http.Error(w, "Failed to remove package from the database. Check the server logs for more information.", http.StatusInternalServerError)
		return
	}

	// Locate package file from Blinky database
	pkgFile, err := a.storage.PackageFile(fmt.Sprintf("%s/%s", targetRepo, targetPkgName))
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to find %s/%s in database.", targetRepo, targetPkgName), http.StatusInternalServerError)
		return
	}

	pkgFile = a.repos[targetRepo] + "/x86_64/" + pkgFile

	// Remove primary package file
	if err := os.Remove(pkgFile); err != nil {
		http.Error(w, fmt.Sprintf("Unable to remove %s because of error: %v", pkgFile, err), http.StatusInternalServerError)
		return
	}

	// Remove signature file
	if err := os.Remove(pkgFile + ".sig"); err != nil {
		log.Printf("Unable to remove %s because of error: %v\n", pkgFile+".sig", err)
	}

	if err := a.storage.DeletePackageFileEntry(fmt.Sprintf("%s/%s", targetRepo, targetPkgName)); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Got error while deleting package file from Blinky database: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
