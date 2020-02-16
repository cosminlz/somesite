package database

import (
	"cabhelp.ro/backend/internal/model"
	"context"
	"github.com/pkg/errors"
)

type UserRoleDB interface {
	GrantRole(ctx context.Context, userID model.UserID, role model.Role) error
	RevokeRole(ctx context.Context, userID model.UserID, role model.Role) error
	GetRolesByUser(ctx context.Context, userID model.UserID) ([]*model.UserRole, error)
}

var grantUserRoleQuery = `
	INSERT INTO user_roles (user_id, role)
	VALUES ($1, $2);
`

func (db *database) GrantRole(ctx context.Context, userID model.UserID, role model.Role) error {
	if _, err := db.conn.ExecContext(ctx, grantUserRoleQuery, userID, role); err != nil {
		return errors.Wrap(err, "could not grant user role")
	}
	return nil
}

var revokeUserRoleQuery = `
	DELETE FROM user_roles
	WHERE user_id = $1 AND role = $2;
`

func (db *database) RevokeRole(ctx context.Context, userID model.UserID, role model.Role) error {
	if _, err := db.conn.ExecContext(ctx, revokeUserRoleQuery, userID, role); err != nil {
		return errors.Wrap(err, "could not revoke user role")
	}
	return nil
}

var getRolesByUserIDQuery = `
	SELECT role
	FROM user_roles
	WHERE user_id = $1;
`

func (db *database) GetRolesByUser(ctx context.Context, userID model.UserID) ([]*model.UserRole, error) {
	var roles []*model.UserRole = make([]*model.UserRole, 0)
	if err := db.conn.SelectContext(ctx, &roles, getRolesByUserIDQuery, userID); err != nil {
		return roles, errors.Wrap(err, "could not get user roles")
	}
	return roles, nil
}
