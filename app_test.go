package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewApp(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		app := NewApp(nil)
		if app.log != defaultLogger {
			t.Errorf("expected default logger, got %#v", app.log)
		}
	})

	t.Run("options", func(t *testing.T) {
		log := nullLogger()
		app := NewApp(nil, WithLogger(log))
		if app.log != log {
			t.Errorf("expected logger %#v, got %#v", log, app.log)
		}
	})
}

func TestHandleHealthz(t *testing.T) {
	app := NewApp(nil)

	r := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	app.Handler().ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP response %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func nullLogger() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithFields(nil)
}
