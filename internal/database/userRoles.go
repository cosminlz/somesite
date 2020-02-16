package database

import (
	"cabhelp.ro/backend/internal/model"
	"context"
)

type UserRoleDB interface {
	GrantRole(ctx context.Context, userID model.UserID, role model.UserRole)
}
