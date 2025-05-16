package keyvaluestore

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

func New(dbDirPath string) (BadgerAdapter, error) {
	db, err := badger.Open(badger.DefaultOptions(dbDirPath))
	if err != nil {
		return BadgerAdapter{}, err
	}

	return BadgerAdapter{db: db}, nil
}

type BadgerAdapter struct { // implements: github.com/BrenekH/blinky.PackageNameToFileProvider
	db *badger.DB
}

func (b *BadgerAdapter) getTxn(key []byte, txn *badger.Txn) (*pkgInfo, error) {
	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	dstBuf, err := item.ValueCopy(nil)
	if err != nil || dstBuf == nil {
		return nil, fmt.Errorf("badger.get view: %w", err)
	}
	pi, err := decodePkgInfo(dstBuf)
	if err != nil {
		// TODO: How handle legacy data?
		return nil, fmt.Errorf("badger.PackageFile decode: %w", err)
	}
	if pi == nil {
		return nil, fmt.Errorf("badger.PackageFile decode: nil pkgInfo")
	}

	return pi, nil
}

func (b *BadgerAdapter) get(packageName string) (*pkgInfo, error) {
	var pi *pkgInfo

	// Convert to bytes outside the txn to reduce time spent in txn.
	keyBytes := []byte(packageName)

	err := b.db.View(func(txn *badger.Txn) error {
		var err error
		pi, err = b.getTxn(keyBytes, txn)
		return err
	})
	return pi, err
}

func (b *BadgerAdapter) PackageFile(packageName string, packageArch string) (string, error) {
	pi, err := b.get(packageName)
	if err != nil {
		return "", fmt.Errorf("badger.PackageFile get: %w", err)
	}

	allowKeys := []string{packageArch, "any", ""}
	for _, k := range allowKeys {
		if _, ok := pi.files[k]; ok {
			return pi.files[k], nil
		}
	}

	return "", fmt.Errorf("badger.PackageFile: no file found for %s", packageName)
}

func (b *BadgerAdapter) StorePackageFile(packageName, packageArch, filePath string) error {
	// Convert to bytes outside the txn to reduce time spent in txn.
	key := []byte(packageName)
	val := []byte(filePath)

	err := b.db.Update(func(txn *badger.Txn) error {

		return txn.Set(key, val)
	})
	if err != nil {
		return fmt.Errorf("badger.StorePackageFile update: %w", err)
	}

	return nil
}

func (b *BadgerAdapter) DeletePackageFileEntry(packageName, packageArch string) error {
	// Convert to bytes outside the txn to reduce time spent in txn.
	key := []byte(packageName)

	err := b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	if err != nil {
		return fmt.Errorf("badger.DeletePackageFileEntry update: %w", err)
	}

	return nil
}
