package user_test

import (
	"testing"
	"time"

	"github.com/cedrickchee/gowebservices/business/auth"
	"github.com/cedrickchee/gowebservices/business/data/user"
	"github.com/cedrickchee/gowebservices/business/tests"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestUser(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	u := user.New(log, db)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := tests.Context()
			now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
			traceID := "00000000-0000-0000-0000-000000000000"

			nu := user.NewUser{
				Name:            "Ced Ric",
				Email:           "cedric@gideon.dev",
				Roles:           []string{auth.RoleAdmin},
				Password:        "sup3rS3cr3tGolang",
				PasswordConfirm: "sup3rS3cr3tGolang",
			}

			usr, err := u.Create(ctx, traceID, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims{
					Issuer:    "gowebservice project",
					Subject:   usr.ID,
					Audience:  "developers",
					ExpiresAt: now.Add(time.Hour).Unix(),
					IssuedAt:  now.Unix(),
				},
				Roles: []string{auth.RoleUser},
			}

			saved, err := u.QueryByID(ctx, traceID, claims, usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)

			if diff := cmp.Diff(usr, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", tests.Success, testID)

			upd := user.UpdateUser{
				Name:  tests.StringPointer("John Doe"),
				Email: tests.StringPointer("john@foo.bar"),
			}

			claims = auth.Claims{
				StandardClaims: jwt.StandardClaims{
					Issuer:    "gowebservice project",
					Audience:  "developers",
					ExpiresAt: now.Add(time.Hour).Unix(),
					IssuedAt:  now.Unix(),
				},
				Roles: []string{auth.RoleAdmin},
			}

			if err := u.Update(ctx, traceID, claims, usr.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", tests.Success, testID)

			saved, err = u.QueryByEmail(ctx, traceID, claims, *upd.Email)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by Email : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by Email.", tests.Success, testID)

			if saved.Name != *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Success, testID)
			}

			if saved.Email != *upd.Email {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Success, testID)
			}

			if err := u.Delete(ctx, traceID, claims, usr.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", tests.Success, testID)

			_, err = u.QueryByID(ctx, traceID, claims, usr.ID)
			if errors.Cause(err) != user.ErrNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve user.", tests.Success, testID)
		}
	}
}
