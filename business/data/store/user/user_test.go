package user_test

import (
	"testing"

	"github.com/Avyukth/service3-clone/business/data/store/user"
	"github.com/Avyukth/service3-clone/business/data/tests"
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
		t.Logf("\t Test %D:\tWhen handling a single User.", testID)
		{

		}
	}
}
