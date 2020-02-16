package database

import (
	"io"

	"github.com/jmoiron/sqlx"
)

// UniqueViolation error is the postgres error string for
// unique index violation
const UniqueViolation = "unique_violation"

// Database interface for database
type Database interface {
	UsersDB
	SessionsDB
	UserRoleDB

	io.Closer
}

type database struct {
	conn *sqlx.DB
}

func (d *database) Close() error {
	return d.conn.Close()
}
