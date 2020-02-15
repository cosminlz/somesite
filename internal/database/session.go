package database

import (
	"context"

	"cabhelp.ro/backend/internal/model"
	"github.com/pkg/errors"
)

type SessionsDB interface {
	SaveRefreshToken(ctx context.Context, session *model.Session) error
	GetSession(ctx context.Context, session model.Session) (*model.Session, error)
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

var getSessionQuery = `
	SELECT user_id, device_id, refresh_token, expires_at
	FROM user_sessions
	WHERE user_id = $1 AND
		  device_id = $2 AND
		  refresh_token = $3 AND
		  to_timestamp(expires_at) > NOW();
`

func (db *database) GetSession(ctx context.Context, inputSession model.Session) (*model.Session, error) {
	var session model.Session
	if err := db.conn.GetContext(ctx, &session, getSessionQuery, inputSession.UserID, inputSession.DeviceID, inputSession.RefreshToken); err != nil {
		return nil, err
	}
	return &session, nil
}
