package blinky

import "net/http"

// The PackageNameToFileProvider interface describes what methods a struct needs
// to implement in order to be used as a way for Blinky to associate a package
// to the last file that was uploaded for it.
type PackageNameToFileProvider interface {
	PackageFile(packageName, packageArch string) (filePath string, err error)
	StorePackageFile(packageName, packageArch, filePath string) (err error)
	DeletePackageFileEntry(packageName, packageArch string) (err error)
}

// The Authenticator interface describes how a given API expects a struct that
// provides authentication capabilities to behave.
type Authenticator interface {
	CreateMiddleware(http.Handler) http.Handler
}
