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

type KeyStore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

func NewMap(store map[string]*rsa.PrivateKey) *KeyStore {
	return &KeyStore{
		store: store,
	}
}

func NewFs(fsys fs.FS) (*KeyStore, error) {

	ks := New()
	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk directory failure: %w", err)
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

	return ks, nil

}

func (ks *KeyStore) Add(privateKey *rsa.PrivateKey, kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.store[kid] = privateKey
}

func (ks *KeyStore) Remove(kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	delete(ks.store, kid)
}

func (ks *KeyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	privateKey, found := ks.store[kid]
	if !found {
		return nil, fmt.Errorf("key %s not found", kid)
	}

	return privateKey, nil
}

func (ks *KeyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	privateKey, found := ks.store[kid]
	if !found {
		return nil, fmt.Errorf("key %s  and corresponding public key not found", kid)
	}

	return &privateKey.PublicKey, nil

}
