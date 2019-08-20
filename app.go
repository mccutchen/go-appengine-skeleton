package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Defaults
var (
	defaultLogger = logrus.New().WithFields(nil)
)

// App is the main application
type App struct {
	db  *DB
	log *logrus.Entry
}

// NewApp creates a new App, customized with the given options
func NewApp(db *DB, opts ...AppOption) *App {
	app := &App{
		db:  db,
		log: defaultLogger,
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// AppOption customizes an App
type AppOption func(a *App)

// WithLogger sets the logger for an App
func WithLogger(log *logrus.Entry) AppOption {
	return func(a *App) {
		a.log = log
	}
}

// Handler handles requests
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", a.HandleHealthz)
	return mux
}

// HandleHealthz receives health check requests
func (a *App) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK\n"))
}

func errResponse(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
