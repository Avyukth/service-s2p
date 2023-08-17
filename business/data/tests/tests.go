package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Avyukth/service3-clone/business/data/schema"
	"github.com/Avyukth/service3-clone/business/sys/database"
	"github.com/Avyukth/service3-clone/foundation/docker"
	"github.com/Avyukth/service3-clone/foundation/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	Success = "\u2713"
	Failure = "\u2717"
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

	if err := schema.Seed(ctx, db); err == nil {
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
