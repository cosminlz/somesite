package api

import (
	"net/http"

	v1 "cabhelp.ro/backend/internal/api/v1"
	"github.com/gorilla/mux"
)

// NewRouter returns a new router
func NewRouter() (http.Handler, error) {
	router := mux.NewRouter()

	router.HandleFunc("/version", v1.VersionHandler)

	return router, nil
}
