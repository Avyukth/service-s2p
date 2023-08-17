package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Avyukth/service3-clone/business/data/store/user"
	"github.com/Avyukth/service3-clone/business/data/tests"
	"github.com/Avyukth/service3-clone/business/sys/auth"
	"github.com/Avyukth/service3-clone/business/sys/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

var dbc = tests.DBContainer{
	Image: "postgres:latest",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestUser(t *testing.T) {

	log, db, teardown := tests.NewUnit(t, dbc)

	t.Cleanup(teardown)

	store := user.NewStore(log, db)

	t.Log("Given the need to work with the user records.")
	{
		testID := 0

		t.Logf("\t Test %d:\tWhen handling a single User.", testID)
		{
			ctx := context.Background()
			now := time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC)

			nu := user.NewUser{

				Name:            "Subhrajit",
				Email:           "subharajit@subhrajit.me",
				Roles:           []string{auth.RoleAdmin},
				Password:        "gophers",
				PasswordConfirm: "gophers",
			}

			usr, err := store.Create(ctx, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			jwtId := uuid.New().String()

			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					// A usual scenario is to set the expiration time relative to the current time
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					NotBefore: jwt.NewNumericDate(time.Now()),
					Issuer:    "service project",
					Subject:   usr.ID,
					ID:        jwtId,
					Audience:  []string{"https://subhrajit.eth"},
				},
				Roles: []string{auth.RoleUser},
			}

			saved, err := store.QueryByID(ctx, claims, usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)

			if diff := cmp.Diff(saved, usr); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s.", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", tests.Success, testID)

			upd := user.UpdateUser{
				Name:  tests.StringPointer("Subhrajit Makur"),
				Email: tests.StringPointer("subharajit@subhrajit.eth"),
			}

			claims = auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					// A usual scenario is to set the expiration time relative to the current time
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					Issuer:    "service project",
				},
				Roles: []string{auth.RoleAdmin},
			}
			if err := store.Update(ctx, claims, usr.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", tests.Success, testID)

			saved, err = store.QueryByEmail(ctx, claims, *upd.Email)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by email: %s.", tests.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by email.", tests.Success, testID)

			if saved.Name != *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to get back the updated user name. ", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould get back the updated user name.", tests.Success, testID)
			}

			if saved.Email != *upd.Email {
				t.Errorf("\t%s\tTest %d:\tShould be able to get back the updated user email. ", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould get back the updated user email.", tests.Success, testID)
			}

			if err := store.Delete(ctx, claims, usr.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", tests.Success, testID)

			_, err = store.QueryByID(ctx, claims, usr.ID)
			if !errors.Is(err, database.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to retrieve user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to retrieve user.", tests.Success, testID)
		}
	}
}
