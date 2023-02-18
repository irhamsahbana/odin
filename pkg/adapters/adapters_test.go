package adapters

import (
	"database/sql"
	"testing"

	sqlEnt "entgo.io/ent/dialect/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAdapter(t *testing.T) {
	is := assert.New(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		is.Failf("failed to open stub db", "%v", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			is.Error(err)
		}
	}(db)

	is.NotNil(db, "mock db is null")
	is.NotNil(mock, "sqlmock is null")

	HelloMysqlOpen = func(dialect, source string) (*sqlEnt.Driver, error) {
		return sqlEnt.NewDriver(dialect, sqlEnt.Conn{ExecQuerier: db}), nil
	}
	HelloSqliteOpen = func(dialect, source string) (*sqlEnt.Driver, error) {
		return sqlEnt.NewDriver(dialect, sqlEnt.Conn{ExecQuerier: db}), nil
	}

	adapter := &Adapter{}
	adapter.Sync(
		WithHelloMysql(&HelloMysql{}),
		WithHelloSqlite(&HelloSqlite{}),
	)

	// Asserts
	mock.ExpectClose()

	is.Nil(adapter.UnSync())
	is.Nil(mock.ExpectationsWereMet())
}
