package clientlib

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/klauspost/compress/zstd"
)

type pkgNameFinder interface {
	FindPkgName(file string) (name string, err error)
}

var pkgNameFinders = map[string]pkgNameFinder{".zst": zstPkgNameFinder{}}

var pkgnameRegex = regexp.MustCompile(`pkgname = (.*)`)

// readPkgName attempts to read the package name stored in the Pacman
// package file.
func readPkgName(file string) (string, error) {
	finder, ok := pkgNameFinders[filepath.Ext(file)]
	if !ok {
		return "", fmt.Errorf("unknown extension %s", filepath.Ext(file))
	}

	return finder.FindPkgName(file)
}

type zstPkgNameFinder struct{}

func (z zstPkgNameFinder) FindPkgName(file string) (name string, err error) {
	// Create a temporary file to read the decompressed bytes into and schedule it's deletion
	tempFile, err := os.CreateTemp("/tmp", "blinky-cli-decompressed-*.pkg.tar")
	if err != nil {
		return "", fmt.Errorf("zst.FindPkgName create temporary file: %w", err)
	}
	defer func() {
		os.Remove(tempFile.Name())
	}()

	// Decompress into temp file
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("zst.FindPkgName open file: %w", err)
	}

	d, err := zstd.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("zst.FindPkgName create zst reader: %w", err)
	}
	defer d.Close()

	_, err = io.Copy(tempFile, d)
	if err != nil {
		return "", fmt.Errorf("zst.FindPkgName copy decompressed file: %w", err)
	}

	tempFile.Seek(0, 0)
	return findPkgNameFromTar(tempFile)
}

func findPkgNameFromTar(tarFile io.Reader) (string, error) {
	// Open and iterate through the files in the archive.
	tr := tar.NewReader(tarFile)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break // End of archive
			}
			return "", fmt.Errorf("read tar file: %w", err)
		}

		if hdr.Name == ".PKGINFO" {
			b, err := io.ReadAll(tr)
			if err != nil {
				return "", fmt.Errorf("read .PKGINFO: %w", err)
			}

			if s := pkgnameRegex.FindStringSubmatch(string(b)); len(s) < 2 {
				return "", errors.New("parse .PKGINFO: could not find pkgname line")
			} else {
				return s[1], nil
			}
		}
	}

	return "", errors.New(".PKGINFO not located in tar archive")
}
