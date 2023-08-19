package tests

import (
	"fmt"
	"net/http"
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
		userToken:  test.Token("user4@example.com", "helloworld"),
		adminToken: test.Token("user3@example.com", "helloservice"),
	}
}
