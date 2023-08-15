package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/Avyukth/service3-clone/business/sys/auth"
	"github.com/golang-jwt/jwt/v5"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestAuth(t *testing.T) {

	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen Handling a single user", testID)
		{
			const KeyID = "12345678901234567890"
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to generate private key: %v", failed, testID, err)
			}
			t.Fatalf("\t%s\tTest %d:\tShould be able to generate private key.", success, testID)
			a, err := auth.New(KeyID, &keyStore{pk: privateKey})
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create authenticator: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create authenticator.", success, testID)
			claims := auth.Claims{

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
			token, err := a.GenerateToken(claims)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to generate JWT: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to generate JWT.", success, testID)

			parsedClaims, err := a.ValidateToken(token)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to parse the claims: %v ", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to parse the claims.", success, testID)

			if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
				t.Logf("\t\tTest %d:\texp: %d", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %d", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the same number of roles: %v", failed, testID, err)

			}
			t.Logf("\t%s\tTest %d:\tShould have the same number of roles.", success, testID)
			if exp, got := claims.Roles[0], parsedClaims.Roles[0]; exp != got {
				t.Logf("\t\tTest %d:\texp: %s", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %s", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected role: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected role.", success, testID)
		}
	}
}

type keyStore struct {
	pk *rsa.PrivateKey
}

func (ks *keyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	return ks.pk, nil
}

func (ks *keyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	return &ks.pk.PublicKey, nil
}
