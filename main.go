package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/sirupsen/logrus"
)

const defaultPort = "9999"

func main() {
	log := logrus.New().WithFields(nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	dsClient, err := datastore.NewClient(ctx, "")
	if err != nil {
		log.WithError(err).Panic("error connecting to datastore")
	}

	db := NewDB(dsClient)

	app := NewApp(db, WithLogger(log))

	log.WithField("port", port).Info("listening ...")
	srv := &http.Server{
		Handler: app.Handler(),
		Addr:    fmt.Sprintf(":%s", port),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Error("error starting server")
	}
}
