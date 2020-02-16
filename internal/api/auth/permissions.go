package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/bluele/gcache"
	"github.com/gorilla/mux"

	"cabhelp.ro/backend/internal/api/utils"
	"cabhelp.ro/backend/internal/database"
	"cabhelp.ro/backend/internal/model"
)

type Permissions interface {
	Wrap(next http.HandlerFunc, permissionTypes ...PermissionType) http.HandlerFunc
	Check(r *http.Request, permissionTypes ...PermissionType) bool
}

type permissions struct {
	DB    database.Database
	cache gcache.Cache
}

func NewPermissions(db database.Database) Permissions {
	p := &permissions{
		DB: db,
	}

	p.cache = gcache.New(200).LRU().LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
		userID := key.(model.UserID)
		roles, err := p.DB.GetRolesByUser(context.Background(), userID)
		if err != nil {
			return nil, nil, err
		}
		expire := 1 * time.Minute

		return roles, &expire, nil
	}).Build()

	return p
}

func (p *permissions) getRoles(userID model.UserID) ([]*model.UserRole, error) {
	roles, err := p.cache.Get(userID)
	if err != nil {
		return nil, err
	}
	return roles.([]*model.UserRole), nil
}

func (p *permissions) withRoles(principal model.Principal, roleFunc func([]*model.UserRole) bool) (bool, error) {
	if principal.UserID == model.NilUserID {
		return false, nil
	}

	roles, err := p.getRoles(principal.UserID)
	if err != nil {
		return false, err
	}

	return roleFunc(roles), nil
}

func (p *permissions) Wrap(next http.HandlerFunc, permissionTypes ...PermissionType) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowed := p.Check(r, permissionTypes...); !allowed {
			utils.WriteError(w, http.StatusUnauthorized, "permission denied", nil)
			return
		}
		next.ServeHTTP(w, r)

	})
}

func (p *permissions) Check(r *http.Request, permissionTypes ...PermissionType) bool {
	principal := GetPrincipal(r)

	for _, permissionType := range permissionTypes {
		switch permissionType {
		case Admin:
			if allowed, _ := p.withRoles(principal, adminOnly); allowed {
				return true
			}
		case Member:
			if allowed := member(principal); allowed {
				return true
			}
		case MemberIsTarget:
			targetUserID := model.UserID(mux.Vars(r)["userID"])
			if allowed := memberIsTarget(targetUserID, principal); allowed {
				return true
			}
		}
	}

	return false
}
