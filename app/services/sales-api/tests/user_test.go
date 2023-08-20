package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Avyukth/service3-clone/app/services/sales-api/handlers"
	"github.com/Avyukth/service3-clone/business/data/store/user"
	"github.com/Avyukth/service3-clone/business/data/tests"
	"github.com/Avyukth/service3-clone/business/sys/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type UserTests struct {
	app        http.Handler
	userToken  string
	adminToken string
}

func TestUsers(t *testing.T) {

	test := tests.NewIntegration(
		t, tests.DBContainer{
			Image: "postgres:latest",
			Port:  "5432",
			Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
		},
	)

	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)

	tests := UserTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken:  test.Token("user@example.com", "hellogopher"),
		adminToken: test.Token("admin@example.com", "hellogopher"),
	}

	t.Run("getToken404", tests.getToken404)
	t.Run("getToken200", tests.getToken200)
	// t.Run("postUser400", tests.postUser400)
	// t.Run("postUser401", tests.postUser401)
	// t.Run("postUser403", tests.postUser403)
	// t.Run("getUser400", tests.getUser400)
	// t.Run("getUser403", tests.getUser403)
	// t.Run("getUser404", tests.getUser404)
	// t.Run("deleteUserNotFound", tests.deleteUserNotFound)
	// t.Run("putUser404", tests.putUser404)
	// t.Run("crudUsers", tests.crudUser)
}

func (ut *UserTests) getToken404(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/users/token", nil)
	w := httptest.NewRecorder()

	r.SetBasicAuth("unknown@example.com", "some-password")
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to deny tokens to unknown users.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen fetching a token with an unrecognized email.", testID)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 404 for the response.", tests.Success, testID)
		}
	}
}

func (ut *UserTests) getToken200(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/users/token", nil)
	w := httptest.NewRecorder()

	r.SetBasicAuth("admin@example.com", "hellogopher")
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to issues tokens to known users")
	{
		testID := 0

		t.Logf("\tTest %d:\tWhen fetching a token with valid credentials", testID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of %d for response :%v", tests.Failed, testID, http.StatusOK, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of %d for response :%v", tests.Success, testID, http.StatusOK, w.Code)
			var got struct {
				Token string `json:"token"`
			}
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response :%v.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response.", tests.Success, testID)
		}
	}

}

func (ut *UserTests) postUser400(t *testing.T) {
	body, err := json.Marshal(&user.NewUser{})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ut.adminToken)
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new user can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete user value.", testID)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 400 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 400 for the response.", tests.Success, testID)

			var got validate.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type : %v", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type.", tests.Success, testID)

			fields := validate.FieldErrors{
				{Field: "name", Error: "name is a required field"},
				{Field: "email", Error: "email is a required field"},
				{Field: "roles", Error: "roles is a required field"},
				{Field: "password", Error: "password is a required field"},
			}
			exp := validate.ErrorResponse{
				Error:  "data validation error",
				Fields: fields.Error(),
			}

			// We can't rely on the order of the field errors so they have to be
			// sorted. Tell the cmp package how to sort them.
			sorter := cmpopts.SortSlices(func(a, b validate.FieldError) bool {
				return a.Field < b.Field
			})

			if diff := cmp.Diff(got, exp, sorter); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
		}
	}
}
