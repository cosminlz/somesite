package database

import (
	"context"

	"cabhelp.ro/backend/internal/model"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type UsersDB interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID *model.UserID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

var ErrUserExists = errors.New("User with that email exists")

var createUserQuery = `
	INSERT INTO users (
		email, password_hash
	)
	VALUES (
		:email, :password_hash
	)
	RETURNING user_id;
`

func (db *database) CreateUser(ctx context.Context, user *model.User) error {

	rows, err := db.conn.NamedQueryContext(ctx, createUserQuery, user)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == UniqueViolation {
				if pqError.Constraint == "user_email" {
					return ErrUserExists // error user exists
				}
			}
		}
		return errors.Wrap(err, "could not insert user in db")
	}

	rows.Next()
	if err := rows.Scan(&user.ID); err != nil {
		return errors.Wrap(err, "Could not get ID for created user")
	}

	return nil
}

var getUserByIDQuery = `
	SELECT user_id, email, password_hash, created_at, deleted_at
	FROM users
	WHERE user_id = $1;
`

func (db *database) GetUserByID(ctx context.Context, userID *model.UserID) (*model.User, error) {
	var user model.User
	if err := db.conn.GetContext(ctx, &user, getUserByIDQuery, userID); err != nil {
		return nil, err
	}

	return &user, nil
}

var getUserByEmailQuery = `
	SELECT user_id, email, password_hash, created_at, deleted_at
	FROM users
	WHERE email = $1;
`

func (db *database) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := db.conn.GetContext(ctx, &user, getUserByEmailQuery, email); err != nil {
		return nil, err
	}

	return &user, nil
}
