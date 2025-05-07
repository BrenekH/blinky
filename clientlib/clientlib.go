package clientlib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

func New(url, username, password string) (*BlinkyClient, error) {
	if url == "" {
		return nil, errors.New("url must not be empty")
	}
	return &BlinkyClient{
		URL:      url,
		Username: username,
		Password: password,
	}, nil
}

type BlinkyClient struct {
	URL      string
	Username string
	Password string
}

// UploadPackageFiles uploads the provided packages, and any .sig files that match the
// package filepaths.
func (b *BlinkyClient) UploadPackageFiles(repo string, packageFilepaths ...string) error {
	for _, path := range packageFilepaths {
		if err := b.UploadPackageFile(repo, path); err != nil {
			return fmt.Errorf("UploadPackageFiles: %w", err)
		}
	}

	return nil
}

// UploadPackageFile reads and uploads the provided package file path as well as any
// matching .sig files.
func (b *BlinkyClient) UploadPackageFile(repo, packageFilepath string) error {
	pkgFile, err := os.Open(packageFilepath)
	if err != nil {
		return fmt.Errorf("UploadPackageFile open package file: %w", err)
	}
	defer pkgFile.Close()

	sigFile, err := os.Open(packageFilepath + ".sig")
	if err != nil {
		sigFile = nil
	} else {
		// Don't attempt to defer a nil pointer
		defer sigFile.Close()
	}

	err = b.UploadPackage(repo, packageFilepath, pkgFile, sigFile)
	if err != nil {
		return fmt.Errorf("UploadPackageFile upload package: %w", err)
	}

	return nil
}

// UploadPackage uploads a package to a repository on a Blinky server. packageFile is required to be non-nil,
// but if there is no signature file to upload, signatureFile may be nil. packageFileName is used to name the
// file on the remote server.
func (b *BlinkyClient) UploadPackage(repo, packageFileName string, packageFile, signatureFile io.Reader) error {
	if packageFile == nil {
		return errors.New("packageFile must not be nil")
	}

	r, w := io.Pipe()
	writer := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer writer.Close()

		pkgFormFile, err := writer.CreateFormFile("package", filepath.Base(packageFileName))
		if err != nil {
			fmt.Printf("UploadPackage create package form file: %v\n", err)
			return
		}

		_, err = io.Copy(pkgFormFile, packageFile)
		if err != nil {
			fmt.Printf("UploadPackage copy package file: %v\n", err)
			return
		}

		if !isNil(signatureFile) {
			sigFormFile, err := writer.CreateFormFile("signature", filepath.Base(packageFileName)+".sig")
			if err != nil {
				fmt.Printf("UploadPackage create signature form file: %v\n", err)
				return
			}

			_, err = io.Copy(sigFormFile, signatureFile)
			if err != nil {
				fmt.Printf("UploadPackage copy signature file: %v\n", err)
				return
			}
		}
	}()

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/unstable/%s/package", b.URL, repo), r)
	if err != nil {
		return fmt.Errorf("UploadPackage create request: %w", err)
	}

	request.SetBasicAuth(b.Username, b.Password)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("UploadPackage perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("UploadPackage read response body: %w", err)
		}

		// TODO: Translate status code into specific Go errors

		return fmt.Errorf("received a %s status code while uploading %s to %s: %s", resp.Status, packageFileName, repo, string(b))
	}

	return nil
}

// RemovePackages deletes the specified packages from the provided repository.
func (b *BlinkyClient) RemovePackages(repo string, packageNames ...string) error {
	for _, pkgName := range packageNames {
		if err := b.RemovePackage(repo, pkgName); err != nil {
			return fmt.Errorf("RemovePackages: %w", err)
		}
	}

	return nil
}

// RemovePackage deletes the specified package from the provided repository.
func (b *BlinkyClient) RemovePackage(repo string, packageName string) error {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/unstable/%s/package/%s", b.URL, repo, packageName), bytes.NewBufferString(""))
	if err != nil {
		return fmt.Errorf("RemovePackage create request: %w", err)
	}

	r.SetBasicAuth(b.Username, b.Password)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("RemovePackage perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("RemovePackage read response body: %w", err)
		}

		// TODO: Translate status code into specific Go errors

		return fmt.Errorf("received a non-200 status code while removing %s/%s: %s - %s", repo, packageName, resp.Status, string(b))
	}

	return nil
}

func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
