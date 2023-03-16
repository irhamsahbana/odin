package adapters

import (
	"testing"

	sqlEnt "entgo.io/ent/dialect/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/infrastructure"
)

func TestWithHelloMysql(t *testing.T) {
	is := assert.New(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		is.Failf("failed to open stub db", "%v", err)
	}

	is.NotNil(db, "mock db is null")
	is.NotNil(mock, "sqlmock is null")

	HelloMysqlOpen = func(dialect, source string) (*sqlEnt.Driver, error) {
		return sqlEnt.NewDriver(dialect, sqlEnt.Conn{ExecQuerier: db}), nil
	}

	infrastructure.Configuration(
		infrastructure.WithPath("../.."),
		infrastructure.WithFilename("config.yaml"),
	).Initialize()

	adapter := &Adapter{}
	adapter.Sync(
		WithHelloMysql(&HelloMysql{
			NetworkDB: NetworkDB{
				Database:    infrastructure.Envs.Mysql.Database,
				User:        infrastructure.Envs.Mysql.User,
				Password:    infrastructure.Envs.Mysql.Password,
				Host:        infrastructure.Envs.Mysql.Host,
				Port:        infrastructure.Envs.Mysql.Port,
				MaxIdleCons: infrastructure.Envs.DB.MaxIdleCons,
			},
		}),
	)

	mock.ExpectClose()

	// Asserts
	is.Nil(adapter.UnSync())
	is.Nil(mock.ExpectationsWereMet())
}
