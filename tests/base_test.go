package tests

import _ "github.com/lib/pq"

import (
	"chirpy/config"
	"chirpy/router"
	"database/sql"
	"errors"
	"github.com/pressly/goose/v3"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
)

type TestServer struct {
	Server  *http.Server
	BaseURL string
}

var testDB *sql.DB

//goland:noinspection HttpUrlsUsage
func Start(t *testing.T) *TestServer {
	t.Helper()

	handler := router.New(config.Init())

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	server := &http.Server{Handler: handler}
	go func() {
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("server error: %v", err)
			return
		}
	}()

	return &TestServer{
		Server:  server,
		BaseURL: "http://" + ln.Addr().String(),
	}
}

func TestMain(m *testing.M) {
	config.Init()

	dbName := "chirpy_test"
	dbURL := os.Getenv("DB_URL")
	adminURL := strings.Replace(dbURL, dbName, "postgres", 1)

	adminDB, err := sql.Open("postgres", adminURL)
	if err != nil {
		log.Fatalf("failed to connect to admin DB: %v", err)
	}

	defer func(adminDB *sql.DB) {
		err := adminDB.Close()
		if err != nil {
			log.Fatalf("failed to close admin DB: %v", err)
		}
	}(adminDB)

	_, _ = adminDB.Exec("DROP DATABASE IF EXISTS " + dbName)
	_, err = adminDB.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		log.Fatalf("failed to create test DB: %v", err)
	}

	testDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose dialect error: %v", err)
	}
	if err := goose.Up(testDB, "../sql/schema"); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	code := m.Run()

	_ = testDB.Close()
	_, _ = adminDB.Exec("DROP DATABASE IF EXISTS " + dbName)

	os.Exit(code)
}
