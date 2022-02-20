package jsonds

import (
	"encoding/json"
	"fmt"
	"os"
)

func New(saveLocation string) (JSONDataStore, error) {
	j := JSONDataStore{filePath: saveLocation, data: make(map[string]string)}

	if err := j.loadFile(); err != nil {
		if err := j.saveFile(); err != nil {
			return j, err
		}
	}

	return j, nil
}

type JSONDataStore struct { // implements: github.com/BrenekH/blinky.PackageNameToFileProvider
	filePath string
	data     map[string]string
}

func (j *JSONDataStore) PackageFile(packageName string) (filePath string, err error) {
	path, ok := j.data[packageName]
	if !ok {
		return "", fmt.Errorf("JSONDataStore.PackageFile: package %s not known in database", packageName)
	}

	return path, nil
}

func (j *JSONDataStore) StorePackageFile(packageName, filePath string) (err error) {
	j.data[packageName] = filePath

	return j.saveFile()
}

func (j *JSONDataStore) loadFile() error {
	b, err := os.ReadFile(j.filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &j.data); err != nil {
		return err
	}

	return nil
}

func (j *JSONDataStore) saveFile() error {
	b, err := json.Marshal(j.data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(j.filePath, b, 0777); err != nil {
		return err
	}

	return nil
}
