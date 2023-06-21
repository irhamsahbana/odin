// Package adapters are the glue between components and external sources.
// # This manifest was generated by ymir. DO NOT EDIT.
package adapters

import (
	"database/sql"
	"testing"

	"entgo.io/ent/dialect"
	sqlEnt "entgo.io/ent/dialect/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/infrastructure"
)

func TestWithHelloSQLite(t *testing.T) {
	is := assert.New(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		is.Failf("failed to open stub db", "%v", err)
	}

	is.NotNil(db, "mock db is null")
	is.NotNil(mock, "sqlmock is null")

	HelloSQLiteOpen = func(d string, db *sql.DB) *sqlEnt.Driver {
		return sqlEnt.NewDriver(dialect.SQLite, sqlEnt.Conn{ExecQuerier: db})
	}

	infrastructure.Configuration(
		infrastructure.WithPath("../.."),
		infrastructure.WithFilename("config.yaml"),
	).Initialize()

	adapter := &Adapter{}
	adapter.Sync(
		WithHelloSQLite(&HelloSQLite{
			File: infrastructure.Envs.HelloSQLite.File,
		}),
	)
	mock.ExpectClose()
	_ = db.Close()

	// Asserts
	is.Nil(adapter.UnSync())
	is.Nil(mock.ExpectationsWereMet())
}
