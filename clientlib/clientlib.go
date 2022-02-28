package clientlib

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func New(url, username, password string) BlinkyClient {
	return BlinkyClient{
		URL:      url,
		Username: username,
		Password: password,
	}
}

type BlinkyClient struct {
	URL      string
	Username string
	Password string
}

func (b *BlinkyClient) UploadPackage() {}

// RemovePackages deletes the specified packages from the provided repository.
func (b *BlinkyClient) RemovePackages(repo string, packageNames ...string) error {
	for _, pkgName := range packageNames {
		if err := b.RemovePackage(repo, pkgName); err != nil {
			return fmt.Errorf("BlinkyClient.RemovePackages: %w", err)
		}
	}

	return nil
}

// RemovePackage deletes the specified package from the provided repository.
func (b *BlinkyClient) RemovePackage(repo string, packageName string) error {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/unstable/%s/package/%s", b.URL, repo, packageName), bytes.NewBufferString(""))
	if err != nil {
		return fmt.Errorf("BlinkyClient.RemovePackage create request: %w", err)
	}

	r.Header.Add("Authorization", b.Password) // TODO: Swap out for basic auth using both username and password

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("BlinkyClient.RemovePackage perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("BlinkyClient.RemovePackage read body: %w", err)
		}

		// TODO: Translate status code into specific Go errors

		return fmt.Errorf("received a non-200 status code while removing %s/%s: %s - %s", repo, packageName, resp.Status, string(b))
	}

	return nil
}
