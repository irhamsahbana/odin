package adapters

import (
	"fmt"

	// mysql go lib.
	_ "github.com/go-sql-driver/mysql"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/rs/zerolog/log"
)

var HelloMysqlOpen = sql.Open // HelloMysqlOpen will invoke to test case.

// HelloMysql is data of instances.
type HelloMysql struct {
	driver *sql.Driver
}

// Open is open the connection of mysql.
func (h *HelloMysql) Open() (*sql.Driver, error) {
	if h.driver == nil {
		return nil, fmt.Errorf("driver was failed to connected")
	}
	return h.driver, nil
}

// Connect is connected the connection of mysql.
func (h *HelloMysql) Connect() (err error) {
	h.driver, err = HelloMysqlOpen(dialect.MySQL,
		"root:root1234@tcp(localhost:33067)/hello?parseTime=true")
	if err != nil {
		log.Error().Err(err).Msg("HelloMysqlOpen is failed to open")
		return err
	}
	return nil
}

// Disconnect is disconnect the connection of mysql.
func (h *HelloMysql) Disconnect() error {
	return h.driver.Close()
}

// WithHelloMysql option function to assign on adapters.
func WithHelloMysql(driver Driver[*sql.Driver]) Option {
	return func(a *Adapter) {
		if err := driver.Connect(); err != nil {
			panic(err)
		}
		open, err := driver.Open()
		if err != nil {
			panic(err)
		}
		a.HelloMysql = open
	}
}
