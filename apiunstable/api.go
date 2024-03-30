package apiunstable

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BrenekH/blinky"
	"github.com/BrenekH/blinky/pacman"
	"github.com/gorilla/mux"
)

func New(storageProvider blinky.PackageNameToFileProvider, authProvider blinky.Authenticator, foundRepos map[string]string, repoArches []string, gnupgDir string, requireSignedPackages, useSignedDB bool) API {
	return API{
		gnupgDir:              gnupgDir,
		repos:                 foundRepos,
		repoArches:            repoArches,
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
	repoArches            []string
	requireSignedPackages bool
	useSignedDB           bool
	gnupgDir              string // The location to set GNUPGHOME to, when repo-add/repo-remove.
}

// Register registers http handlers associated with the unstable API.
func (a *API) Register(router *mux.Router) {
	router.Handle("/{repo}/package", a.auth.CreateMiddleware(http.HandlerFunc(a.putRepoPkg))).Methods(http.MethodPut)
	router.Handle("/{repo}/package/{package_name}", a.auth.CreateMiddleware(http.HandlerFunc(a.deleteRepoPkg))).Methods(http.MethodDelete)
}

func (a *API) putRepoPkg(w http.ResponseWriter, r *http.Request) {
	// Extract variables from URL
	vars := mux.Vars(r)
	targetRepo := vars["repo"]

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
		saveMultipartFile(formSigFile, formPkgHeader, a.repos[targetRepo]+"/tmp")
		formPkgHeader.Filename = temp
	}

	// This is after the signature file so that if the server requires a signed package, the file doesn't get copied
	// until the request is known to have a .sig file. This avoids unnecessary downloading.
	saveMultipartFile(formPkgFile, formPkgHeader, a.repos[targetRepo]+"/tmp")

	packageInfo, err := pkgInfoParseFile(a.repos[targetRepo] + "/tmp/" + formPkgHeader.Filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not read package info. Check that the file provided is a valid Pacman package", http.StatusBadRequest)
		return
	}

	var targetArch string
	for _, arch := range a.repoArches {
		if arch == packageInfo.Arch {
			targetArch = arch
		}
	}
	if targetArch == "" {
		http.Error(w, "Provided package is for an architecture not supported by this server", http.StatusBadRequest)
		return
	}

	// FIXME: Too verbose filepath
	if err := os.Rename(a.repos[targetRepo]+"/tmp/"+formPkgHeader.Filename, a.repos[targetRepo]+"/"+targetArch+"/"+formPkgHeader.Filename); err != nil {
		log.Println(err)
		http.Error(w, "Failed to move package from temp directory. Check the server logs for more information.", http.StatusInternalServerError)
		return
	}

	if targetArch == "any" {
		for _, arch := range a.repoArches {
			// Make symlink for "any" packages
			if err := os.Symlink(a.repos[targetRepo]+"/"+targetArch+"/"+formPkgHeader.Filename, a.repos[targetRepo]+"/"+arch+"/"+formPkgHeader.Filename); err != nil {
				log.Println(err)
				http.Error(w, "Failed to create symlink for package. Check the server logs for more information.", http.StatusInternalServerError)
				return
			}

			// Make symlink for signature file
			if err := os.Symlink(a.repos[targetRepo]+"/"+targetArch+"/"+formPkgHeader.Filename+".sig", a.repos[targetRepo]+"/"+arch+"/"+formPkgHeader.Filename+".sig"); err != nil {
				log.Println(err)
				http.Error(w, "Failed to create symlink for signature file. Check the server logs for more information.", http.StatusInternalServerError)
				return
			}

			// Run repo-add for each architecture
			if err := pacman.RepoAdd(a.repos[targetRepo]+"/"+arch+"/"+targetRepo+".db.tar.gz", a.repos[targetRepo]+"/"+arch+"/"+formPkgHeader.Filename, a.useSignedDB, &a.gnupgDir); err != nil {
				log.Println(err)
				http.Error(w, "Failed to add package to the database. Check the server logs for more information.", http.StatusInternalServerError)
				return
			}
		}
	} else {
		if err := pacman.RepoAdd(a.repos[targetRepo]+"/"+targetArch+"/"+targetRepo+".db.tar.gz", a.repos[targetRepo]+"/"+targetArch+"/"+formPkgHeader.Filename, a.useSignedDB, &a.gnupgDir); err != nil {
			log.Println(err)
			http.Error(w, "Failed to add package to the database. Check the server logs for more information.", http.StatusInternalServerError)
			return
		}

		if err := a.storage.StorePackageFile(fmt.Sprintf("%s/%s", targetRepo, packageInfo.Name), formPkgHeader.Filename); err != nil {
			log.Println(err)
			http.Error(w, "Got error while saving new Blinky db entry. Check server logs for more information.", http.StatusInternalServerError)
			return
		}
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

	// TODO: Handle "any" architectures
	for _, arch := range a.repoArches {
		if err := pacman.RepoRemove(a.repos[targetRepo]+"/"+arch+"/"+targetRepo+".db.tar.gz", targetPkgName, a.useSignedDB, &a.gnupgDir); err != nil {
			log.Printf("%s", err)
			http.Error(w, "Failed to remove package from the database. Check the server logs for more information.", http.StatusInternalServerError)
			return
		}

		// Locate package file from Blinky database
		pkgFile, err := a.storage.PackageFile(fmt.Sprintf("%s/%s", targetRepo, targetPkgName))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to find %s/%s in database.", targetRepo, targetPkgName), http.StatusInternalServerError)
			return
		}

		pkgFile = a.repos[targetRepo] + "/" + arch + "/" + pkgFile

		// Remove primary package file
		if err := os.Remove(pkgFile); err != nil {
			http.Error(w, fmt.Sprintf("Unable to remove %s because of error: %v", pkgFile, err), http.StatusInternalServerError)
			return
		}

		// Remove signature file
		if err := os.Remove(pkgFile + ".sig"); err != nil {
			log.Printf("Unable to remove %s because of error: %v\n", pkgFile+".sig", err)
		}
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

func saveMultipartFile(mFile multipart.File, header *multipart.FileHeader, targetDir string) error {
	cleanTargetDir := filepath.Clean(targetDir)
	dest := cleanTargetDir + "/" + filepath.Base(filepath.Clean(header.Filename))

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
