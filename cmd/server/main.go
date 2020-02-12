package main

import (
	"net/http"

	"cabhelp.ro/backend/internal/api"
	"cabhelp.ro/backend/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.WithField("version", config.Version).Info("Starting...")

	router, err := api.NewRouter()
	if err != nil {
		logrus.WithError(err).Error("Error creating router")
	}

	const addr = "0.0.0.0:8088"
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	if _, err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Cannot start server")
	}
}
