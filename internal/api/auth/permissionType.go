package auth

import (
	"cabhelp.ro/backend/internal/model"
)

type PermissionType string

const (
	Admin          PermissionType = "admin"
	Member         PermissionType = "member"
	MemberIsTarget PermissionType = "memberIsTarget"
)

var adminOnly = func(roles []*model.UserRole) bool {
	for _, role := range roles {
		switch role.Role {
		case model.RoleAdmin:
			return true
		}
	}
	return false
}

var member = func(principal model.Principal) bool {
	return principal.UserID != ""
}

var memberIsTarget = func(userID model.UserID, principal model.Principal) bool {
	if userID == "" || principal.UserID == "" {
		return false
	}

	if userID != principal.UserID {
		return false
	}

	return true
}
