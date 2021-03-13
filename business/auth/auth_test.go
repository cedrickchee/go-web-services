package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cedrickchee/gowebservices/business/auth"
	"github.com/dgrijalva/jwt-go"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestAuth(t *testing.T) {
	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single user.", testID)
		{
			// privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateRSAKey))
			// if err != nil {
			// 	t.Fatalf("\t%s\tTest %d:\tShould be able to parse the private key from pem: %v", failed, testID, err)
			// }
			// t.Logf("\t%s\tTest %d:\tShould be able to parse the private key from pem.", success, testID)

			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				log.Fatalln(err)
			}

			// The key id we are stating represents the public key in the
			// public key store.
			const keyID = "91ff8662-82e2-10a6-cf5e-5ea7d0f0c1c3"
			lookup := func(kid string) (*rsa.PublicKey, error) {
				switch kid {
				case keyID:
					return &privateKey.PublicKey, nil
				}
				return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
			}

			a, err := auth.New("RS256", lookup, auth.Keys{keyID: privateKey})
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create an authenticator: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create an authenticator.", success, testID)

			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims{
					Issuer:    "gowebservice project",
					Subject:   "8ea09532-9245-8623-923d-3201212966b1",
					Audience:  "developer",
					ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				Roles: []string{auth.RoleAdmin},
			}

			token, err := a.GenerateToken(keyID, claims)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to generate a JWT: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to generate a JWT.", success, testID)

			parsedClaims, err := a.ValidateToken(token)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to parse the claims: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to parse the claims.", success, testID)

			if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
				t.Logf("t\tTest %d:\texp: %d", testID, exp)
				t.Logf("t\tTest %d:\tgot: %d", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected number of roles: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected number of roles.", success, testID)

			if exp, got := claims.Roles[0], parsedClaims.Roles[0]; exp != got {
				t.Logf("\t\tTest %d:\texp: %v", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %v", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected roles: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected roles.", success, testID)
		}
	}
}
