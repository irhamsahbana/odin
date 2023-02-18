package adapters

import (
	"fmt"

	// sqlite go lib.
	_ "github.com/mattn/go-sqlite3"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/rs/zerolog/log"
)

var HelloSqliteOpen = sql.Open // HelloSqliteOpen will invoke to test case.

// HelloSqlite is data of instances.
type HelloSqlite struct {
	driver *sql.Driver
}

// Open is open the connection of sqlite.
func (h *HelloSqlite) Open() (*sql.Driver, error) {
	if h.driver == nil {
		return nil, fmt.Errorf("driver was failed to connected")
	}
	return h.driver, nil
}

// Connect is connected the connection of sqlite.
func (h *HelloSqlite) Connect() (err error) {
	h.driver, err = HelloSqliteOpen(dialect.SQLite,
		"sqlite://../.data/hello.migration.db?_fk=1")
	if err != nil {
		log.Error().Err(err).Msg("HelloSqliteOpen is failed to open")
		return err
	}
	return nil
}

// Disconnect is disconnect the connection of sqlite.
func (h *HelloSqlite) Disconnect() error {
	return h.driver.Close()
}

// WithHelloSqlite option function to assign on adapters.
func WithHelloSqlite(driver Driver[*sql.Driver]) Option {
	return func(a *Adapter) {
		if err := driver.Connect(); err != nil {
			panic(err)
		}
		open, err := driver.Open()
		if err != nil {
			panic(err)
		}
		a.HelloSqlite = open
	}
}
