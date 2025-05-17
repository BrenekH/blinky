package keyvaluestore

import (
	"bytes"
	"encoding/gob"
)

type pkgInfo struct {
	name  string
	files map[string]string // map[arch]filePath
}

func newPkgInfo(name string, files map[string]string) *pkgInfo {
	return &pkgInfo{
		name:  name,
		files: files,
	}
}

func (p *pkgInfo) encode() []byte {
	return encodePkgInfo(p)
}

func encodePkgInfo(pkg *pkgInfo) []byte {

	buf := bytes.NewBuffer(nil)
	_ = gob.NewEncoder(buf).Encode(&pkg)
	return buf.Bytes()
}

func decodePkgInfo(data []byte) (*pkgInfo, error) {
	buf := bytes.NewBuffer(data)
	var h pkgInfo
	err := gob.NewDecoder(buf).Decode(&h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
