package v1

import (
	"encoding/json"
	"net/http"

	"cabhelp.ro/backend/internal/api/utils"
	"cabhelp.ro/backend/internal/database"
	"cabhelp.ro/backend/internal/model"

	"github.com/sirupsen/logrus"
)

// UserAPI provides the handlers for user related routes
type UserAPI struct {
	DB database.Database
}

type UserParameters struct {
	model.User
	Password string `json:"password"`
}

func (api *UserAPI) Create(w http.ResponseWriter, r *http.Request) {

	logger := logrus.WithField("func", "user.go:Create()")

	// load parameters
	var userParams UserParameters
	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logrus.WithFields(logrus.Fields{
		"email": *userParams.Email,
	})

	// verify that input is valid
	if err := userParams.Verify(); err != nil {
		logger.WithError(err).Warn("Invalid fields")
		utils.WriteError(w, http.StatusBadRequest, "Invalid fields", map[string]string{
			"error": err.Error(),
		})
	}

	// hash password
	hashed, err := model.HashPassword(userParams.Password)
	if err != nil {
		logger.WithError(err).Warn("could not hash password")
		utils.WriteError(w, http.StatusInternalServerError, "could not hash password", nil)
		return
	}

	newUser := &model.User{
		Email:        userParams.Email,
		PasswordHash: &hashed,
	}

	// create user in user table
	ctx := r.Context()
	if err := api.DB.CreateUser(ctx, newUser); err != nil {
		if err == database.ErrUserExists {
			logger.WithError(err).Warn("User already exists")
			utils.WriteError(w, http.StatusConflict, "User already exists", nil)
		} else {
			logger.WithError(err).Warn("Error creating user")
			utils.WriteError(w, http.StatusInternalServerError, "Error creating user", map[string]string{
				"error": err.Error(),
			})
		}
		return
	}

	// retrieve the user
	createdUser, err := api.DB.GetUserByID(ctx, &newUser.ID)
	if err != nil {
		logger.WithError(err).Warn("Could not retrive user")
		utils.WriteError(w, http.StatusInternalServerError, "Could not retrieve user", nil)
		return
	}

	// return user info
	logger.Info("User created")
	utils.WriteJSON(w, http.StatusCreated, createdUser)
}

func (api *UserAPI) Login(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go:Login()")

	var credentials model.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		logger.WithError(err).Warn("could not decode credentials")
		utils.WriteError(w, http.StatusBadRequest, "could not decode credentials", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger = logrus.WithFields(logrus.Fields{
		"email": credentials.Email,
	})

	// Get user by email
	ctx := r.Context()
	user, err := api.DB.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		logger.WithError(err).Warn("Could not retrive user")
		utils.WriteError(w, http.StatusUnauthorized, "Invalid user or password", nil)
		return
	}

	// Check password
	if err := user.CheckPassword(credentials.Password); err != nil {
		logger.WithError(err).Warn("Invalid password")
		utils.WriteError(w, http.StatusUnauthorized, "Invalid user or password", nil)
		return
	}

	logger.WithField("userID", user.ID).Info("Logged in")
	utils.WriteJSON(w, http.StatusOK, user)
}
