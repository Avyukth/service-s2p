package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Avyukth/service3-clone/app/services/sales-api/handlers"
	"github.com/Avyukth/service3-clone/business/data/tests"
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

	t.Run("getToken200", tests.getToken200)
}

func (ut *UserTests) getToken200(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/users/token", nil)
	w := httptest.NewRecorder()

	r.SetBasicAuth("admin@example.com", "helloworld")
	ut.app.ServeHTTP(w, r)

}
