// +build integration

package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	host, stop, err := startTestDb(ctx)
	if err != nil {
		log.Fatalf("error starting test database: %s", err)
	}
	defer stop()

	os.Setenv("DATASTORE_EMULATOR_HOST", host)

	time.Sleep(10 * time.Second)

	// handle the panic and clean up testdb if tests are run with -timeout=xxx
	defer func() {
		if err := recover(); err != nil {
			stop()
		}
	}()

	result := m.Run()
	stop()
	os.Exit(result)
}

func TestHandleReceiveIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	dsClient, err := datastore.NewClient(ctx, testProjectID)
	if err != nil {
		t.Fatalf("error connecting to testdb: %s", err)
	}
	defer dsClient.Close()

	app := NewApp(
		NewDB(dsClient),
		WithLogger(nullLogger()),
	)
	handler := app.Handler()

	t.Run("health check", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/healthz", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected HTTP response %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}
