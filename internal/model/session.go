package model

type DeviceID string

var NilDeviceID DeviceID

type Session struct {
	UserID       UserID   `db:"user_id"`
	DeviceID     DeviceID `db:"device_id"`
	RefreshToken string   `db:"refresh_token"`
	ExpiresAt    int64    `db:"expires_at"`
}

type SessionData struct {
	DeviceID DeviceID `json:"deviceID"`
}
