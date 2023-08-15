package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/Avyukth/service3-clone/business/sys/auth"
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
				t.Fatalf("\t%s\tTest %d:\tShould be able to generate private key.", failed, testID, err)
			}
			t.Fatalf("\t%s\tTest %d:\tShould be able to generate private key.", success, testID)
			a, err := auth.New(KeyID, &keyStore{pk: privateKey})
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create authenticator.", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create authenticator.", success, testID)

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
