package blinky

// The PackageNameToFileProvider interface describes what methods a struct needs
// to implement in order to be used as a way for Blinky to associate a package
// to the last file that was uploaded for it.
type PackageNameToFileProvider interface {
	PackageFile(packageName string) (filePath string, err error)
	StorePackageFile(packageName, filePath string) (err error)
}
