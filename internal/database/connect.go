package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	databaseURL     = flag.String("database-url", "postgres://data_processor:uber1234@localhost:5432/cabhelp?sslmode=disable", "Database URL")
	databaseTimeout = flag.Int64("database-timeout-ms", 2000, "")
)

func Connect() (*sqlx.DB, error) {

	dbURL := *databaseURL

	logrus.WithField("dbURL", dbURL).Info("Connecting to DBi...")
	conn, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to db")
	}

	conn.SetMaxOpenConns(32)

	if err := waitForDB(conn.DB); err != nil {
		return nil, err
	}

	if err := migrateDb(conn.DB); err != nil {
		return nil, errors.Wrap(err, "could not migrate")
	}

	return conn, nil
}

// New returns a new instance of a database
func New() (Database, error) {

	conn, err := Connect()
	if err != nil {
		return nil, err
	}

	d := &database{
		conn: conn,
	}
	return d, nil
}

func waitForDB(conn *sql.DB) error {
	ready := make(chan struct{})
	go func() {
		for {
			if err := conn.Ping(); err == nil {
				close(ready)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(time.Duration(*databaseTimeout) * time.Millisecond):
		return errors.New("database not ready")
	}
}
