package main

import (
	"net/http"

	"cabhelp.ro/backend/internal/api"
	"cabhelp.ro/backend/internal/api/auth"
	"cabhelp.ro/backend/internal/config"
	"cabhelp.ro/backend/internal/database"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
)

func main() {

	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	logrus.WithField("version", config.Version).Info("Starting...")

	// Create auth module
	tokens := auth.NewTokens()

	// Create DB
	db, err := database.New()
	if err != nil {
		logrus.WithError(err).Fatal("Could not connect to db")
	}

	// Create router
	router, err := api.NewRouter(db, tokens)
	if err != nil {
		logrus.WithError(err).Fatal("Error creating router")
	}

	const addr = "0.0.0.0:8088"
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Cannot start server")
	}
}
