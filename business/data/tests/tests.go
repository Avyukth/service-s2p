package tests

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Avyukth/service3-clone/business/data/schema"
	"github.com/Avyukth/service3-clone/business/data/store/user"
	"github.com/Avyukth/service3-clone/business/sys/auth"
	"github.com/Avyukth/service3-clone/business/sys/database"
	"github.com/Avyukth/service3-clone/foundation/docker"
	"github.com/Avyukth/service3-clone/foundation/keystore"
	"github.com/Avyukth/service3-clone/foundation/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

type DBContainer struct {
	Image string
	Port  string
	Args  []string
}

func NewUnit(t *testing.T, dbc DBContainer) (*zap.SugaredLogger, *sqlx.DB, func()) {

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	c := docker.StartContainer(t, dbc.Image, dbc.Port, dbc.Args...)

	db, err := database.Open(database.Config{
		Host:       "localhost",
		User:       "postgres",
		Password:   "postgres",
		Name:       "postgres",
		DisableTLS: true,
	})

	if err != nil {
		t.Fatalf("Opening database connection %v", err)
	}
	t.Log("Waiting for database to be ready...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(t, c.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	if err := schema.Seed(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(t, c.ID)
		t.Fatalf("Seeding error: %s", err)
	}

	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("logger error: %s", err)
	}

	teardown := func() {

		t.Helper()
		docker.StopContainer(t, c.ID)
		log.Sync()

		w.Close()
		var buf bytes.Buffer

		io.Copy(&buf, r)
		os.Stdout = old
		fmt.Println("************************ LOGS ************************")
		fmt.Print(buf.String())
		fmt.Println("************************ LOGS ************************")
	}

	return log, db, teardown
}

type Test struct {
	DB       *sqlx.DB
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	t        *testing.T
	Teardown func()
}

func NewIntegration(t *testing.T, dbc DBContainer) *Test {

	log, db, teardown := NewUnit(t, dbc)

	keyID := "133d7df7-d74c-4802-985c-f4a64e696f47"

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Generating private key: %s", err)
	}

	auth, err := auth.New(keyID, keystore.NewMap(map[string]*rsa.PrivateKey{keyID: privateKey}))
	if err != nil {
		t.Fatalf("Auth error: %s", err)
	}

	test := Test{
		DB:       db,
		Log:      log,
		Auth:     auth,
		t:        t,
		Teardown: teardown,
	}
	return &test
}

func (test *Test) Token(email, pass string) string {

	test.t.Log("Generating token for tests ...")

	store := user.NewStore(test.Log, test.DB)
	claims, err := store.Authenticate(context.Background(), time.Now(), email, pass)
	if err != nil {
		test.t.Fatalf("Authenticating error: %s", err)
	}
	token, err := test.Auth.GenerateToken(claims)
	if err != nil {
		test.t.Fatalf("Generating token error: %s", err)
	}

	return token
}

func StringPointer(s string) *string {
	return &s
}

func IntPointer(i int) *int {
	return &i
}

// May be removed in the future
func deleteDB(t *testing.T, ctx context.Context, db *sqlx.DB) error {
	if err := schema.Seed(ctx, db); err != nil {
		t.Logf("Deleting error: %s", err)
		return err
	}
	return nil
}
