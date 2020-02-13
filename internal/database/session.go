package database

import (
	"context"

	"cabhelp.ro/backend/internal/model"
	"github.com/pkg/errors"
)

type SessionsDB interface {
	SaveRefreshToken(ctx context.Context, session *model.Session) error
}

var saveRefreshTokenQuery = `
	INSERT INTO user_sessions(user_id, device_id, refresh_token, expires_at)
	VALUES (:user_id, :device_id, :refresh_token, :expires_at)

	ON CONFLICT (user_id, device_id) 
	DO
		UPDATE
			SET refresh_token = :refresh_token,
				expires_at = :expires_at;
`

func (db *database) SaveRefreshToken(ctx context.Context, session *model.Session) error {
	if _, err := db.conn.NamedQueryContext(ctx, saveRefreshTokenQuery, session); err != nil {
		return errors.Wrap(err, "could not save refresh token")
	}
	return nil
}
