package api

import (
	"net/http"

	"cabhelp.ro/backend/internal/api/auth"
	v1 "cabhelp.ro/backend/internal/api/v1"
	"github.com/gorilla/mux"

	"cabhelp.ro/backend/internal/database"
)

// NewRouter returns a new router
func NewRouter(db database.Database, tokens auth.Tokens) (http.Handler, error) {

	router := mux.NewRouter()
	router.HandleFunc("/version", v1.VersionHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	userAPI := &v1.UserAPI{
		DB:     db,
		Tokens: tokens,
	}

	apiRouter.HandleFunc("/users", userAPI.Create).Methods("POST") // create user
	// apiRouter.HandleFunc("/users", userAPI.GetUsers).Methods("GET") // list all users
	// apiRouter.HandleFunc("/users/{userID}", userAPI.GetUserByID).Methods("GET") // get user by ID
	// apiRouter.HandleFunc("/users/{userID}", userAPI.GetUserByID).Methods("PATCH") // update user
	// apiRouter.HandleFunc("/users/{userID}", userAPI.GetUserByID).Methods("DELETE") // delete user

	apiRouter.HandleFunc("/login", userAPI.Login).Methods("POST") // create user

	return router, nil
}
