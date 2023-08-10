package keystore

import (
	"crypto/rsa"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

type Keystore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

func New() *Keystore {
	return &Keystore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

func NewMap(store map[string]*rsa.PrivateKey) *Keystore {
	return &Keystore{
		store: store,
	}
}

func NewFs(fsys fs.FS) (*Keystore, error) {

	ks := Keystore{
		store: make(map[string]*rsa.PrivateKey),
	}
	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dirEntry.IsDir() {
			return nil
		}
		if path.Ext(fileName) != ".pem" {
			return nil
		}
		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", fileName, err)
		}
		defer file.Close()
		// =================================================================

		privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("failed to auth private key: %w", err)
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
		if err != nil {
			return fmt.Errorf("failed to parse auth private key: %w", err)
		}

		ks.store[strings.TrimSuffix(dirEntry.Name(), ".pem")] = privateKey

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return &ks, nil

}
