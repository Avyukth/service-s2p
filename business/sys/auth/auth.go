package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// private and public keys for JWT use
type KeyLookup interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	PublicKey(kid string) (*rsa.PublicKey, error)
}

// Auth is to used for authenticate clients
type Auth struct {
	activeKID string
	ketLookup KeyLookup
	method    jwt.SigningMethod
	keyFunc   func(t *jwt.Token) (interface{}, error)
	parser    jwt.Parser
}

func New(activeKID string, ketLookup KeyLookup) (*Auth, error) {
	_, err := ketLookup.PrivateKey(activeKID)
	if err != nil {
		return nil, errors.New("active key does not exist in store")
	}
	method := jwt.GetSigningMethod("RS256")
	if method == nil {
		return nil, errors.New("unsupported signing method RS256")
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing kid in token header")
		}
		kidID, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id (kid) in token header must be string")
		}
		return ketLookup.PrivateKey(kidID)

	}

	parser := jwt.Parser{}

	a := Auth{
		activeKID: activeKID,
		ketLookup: ketLookup,
		method:    method,
		keyFunc:   keyFunc,
		parser:    parser,
	}

	return &a, nil

}

func (a *Auth) GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["kid"] = a.activeKID

	privateKey, err := a.ketLookup.PrivateKey(a.activeKID)
	if err != nil {
		return "", errors.New("kid lookup Failed")
	}

	tokenSString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("token signing Failed: %v", err)
	}
	return tokenSString, nil
}

func (a *Auth) ValidateToken(tokenString string) (Claims, error) {
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenString, claims, a.keyFunc)

	if err != nil {
		return claims, fmt.Errorf("token parsing Failed: %w", err)
	}

	if !token.Valid {

		return claims, errors.New("token is not valid")
	}
	return claims, nil
}
