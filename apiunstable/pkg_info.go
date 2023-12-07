package apiunstable

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/klauspost/compress/zstd"
)

var pkgnameRegex = regexp.MustCompile(`pkgname = (.*)`)

type pkgInfo struct {
	Name string
}

func pkgInfoParseFile(filepath string) (pkgInfo, error) {
	tarFile, err := decompressPackage(filepath)
	if err != nil {
		return pkgInfo{}, fmt.Errorf("pkgInfo.ParseFile: %w", err)
	}
	defer func() {
		os.Remove(tarFile.Name())
	}()

	pkgInfoStr, err := extractPkgInfoFromTar(tarFile)
	if err != nil {
		return pkgInfo{}, fmt.Errorf("pkgInfo.ParseFile: %w", err)
	}

	p := pkgInfo{}

	if s := pkgnameRegex.FindStringSubmatch(pkgInfoStr); len(s) < 2 {
		return pkgInfo{}, errors.New("pkg.ParseFile: parse .PKGINFO: could not find pkgname line")
	} else {
		p.Name = s[1]
	}

	// TODO: Parse for architecture information

	return p, nil
}

func decompressPackage(filepath string) (*os.File, error) {
	// Create a temporary file to read the decompressed bytes into and schedule it's deletion
	tempFile, err := os.CreateTemp("/tmp", "blinky-decompressed-*.pkg.tar")
	if err != nil {
		return nil, fmt.Errorf("zst.FindPkgName create temporary file: %w", err)
	}

	// Decompress into temp file
	f, err := os.Open(filepath)
	if err != nil {
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("decompressPackage open file: %w", err)
	}

	d, err := zstd.NewReader(f)
	if err != nil {
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("decompressPackage create zst reader: %w", err)
	}
	defer d.Close()

	_, err = io.Copy(tempFile, d)
	if err != nil {
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("decompressPackage copy decompressed file: %w", err)
	}

	tempFile.Seek(0, 0)
	return tempFile, nil
}

func extractPkgInfoFromTar(tarFile io.Reader) (string, error) {
	// Open and iterate through the files in the archive.
	tr := tar.NewReader(tarFile)
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break // End of archive
			}
			return "", fmt.Errorf("read tar file: %w", err)
		}

		if header.Name == ".PKGINFO" {
			b, err := io.ReadAll(tr)
			if err != nil {
				return "", fmt.Errorf("read .PKGINFO: %w", err)
			}

			return string(b), nil
		}
	}

	return "", errors.New(".PKGINFO not located in tar archive")
}
