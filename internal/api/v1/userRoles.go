package v1

import (
	"cabhelp.ro/backend/internal/api/auth"
	"cabhelp.ro/backend/internal/api/utils"
	"cabhelp.ro/backend/internal/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (api *UserAPI) GrantRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "userRoles.go:GrantRole")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	var userRole model.UserRole
	if err := json.NewDecoder(r.Body).Decode(&userRole); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.GrantRole(ctx, userID, userRole.Role); err != nil {
		logger.WithError(err).Warn("Error granting role to user")
		utils.WriteError(w, http.StatusInternalServerError, "Error granting role to user", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, &ActCreated{
		Created: true,
	})
}

func (api *UserAPI) RevokeRole(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "userRoles.go:RevokeRole")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	var userRole model.UserRole
	if err := json.NewDecoder(r.Body).Decode(&userRole); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.WriteError(w, http.StatusBadRequest, "could not decode parameters", map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx := r.Context()

	if err := api.DB.RevokeRole(ctx, userID, userRole.Role); err != nil {
		logger.WithError(err).Warn("Error revoking role to user")
		utils.WriteError(w, http.StatusInternalServerError, "Error revoking role to user", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, &ActDeleted{
		Deleted: true,
	})
}

func (api *UserAPI) GetRoleList(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "userRoles.go:GetRoleList")

	vars := mux.Vars(r)
	userID := model.UserID(vars["userID"])
	principal := auth.GetPrincipal(r)

	logger = logger.WithFields(logrus.Fields{
		"userID":    userID,
		"principal": principal,
	})

	ctx := r.Context()

	roles, err := api.DB.GetRolesByUser(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("Error getting roles for user")
		utils.WriteError(w, http.StatusInternalServerError, "Error getting roles for user", nil)
		return
	}

	utils.WriteJSON(w, http.StatusOK, &roles)
}
