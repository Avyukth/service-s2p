package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	err := genKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
