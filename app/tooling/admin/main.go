package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	err := genToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func genToken() error {
	claims := struct {
		jwt.RegisteredClaims `json:"registered"`
		Roles                []string `json:"roles"`
	}{
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "service Project",
			Subject:   "123456789",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
		[]string{"ADMIN"},
	}
	method := jwt.SigningMethodRS256
	token := jwt.NewWithClaims(method, claims.RegisteredClaims)
	token.Header["kid"] = "995fcf0b-73d9-4516-a32f-4e85d723f06d"

	name := "zarf/keys/995fcf0b-73d9-4516-a32f-4e85d723f06d.pem"

	file, err := os.Open(name)
	if err != nil {
		return err
	}

	privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))

	if err != nil {
		return fmt.Errorf("reading auth private key: %w", err)
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key: %w", err)
	}

	ss, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("====================TOKEN STARTED======================")
	fmt.Println(ss)

	fmt.Println("====================TOKEN ENDED======================")

	return nil
}

func genKey() error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return err
	}

	privateKeyFile, err := os.OpenFile("private.pem", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating private.pem file: %w", err)
	}

	defer privateKeyFile.Close()

	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateKeyFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private file: %v", err)
	}

	// ============================================================================================================================
	// Public Pem file creation

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	publicKeyFile, err := os.OpenFile("public.pem", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

	if err != nil {
		return fmt.Errorf("creating public.pem file: %w", err)
	}
	defer publicKeyFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	if err := pem.Encode(publicKeyFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %v", err)
	}

	fmt.Println("private and public keys generated successfully")
	return nil
}
