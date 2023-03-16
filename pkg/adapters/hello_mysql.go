package adapters

import (
	"fmt"
	"net"
	"strconv"
	"time"

	// mysql go lib.
	"github.com/go-sql-driver/mysql"

	"entgo.io/ent/dialect"
	sqlEnt "entgo.io/ent/dialect/sql"
	"github.com/rs/zerolog/log"
)

var HelloMysqlOpen = sqlEnt.Open // HelloMysqlOpen will invoke to test case.

// HelloMysql is data of instances.
type HelloMysql struct {
	NetworkDB
	driver *sqlEnt.Driver
}

// Open is open the connection of mysql.
func (h *HelloMysql) Open() (*sqlEnt.Driver, error) {
	if h.driver == nil {
		return nil, fmt.Errorf("driver was failed to connected")
	}
	return h.driver, nil
}

// Connect is connected the connection of mysql.
func (h *HelloMysql) Connect() (err error) {
	h.driver, err = HelloMysqlOpen(dialect.MySQL, h.dsn())
	if err != nil {
		log.Error().Err(err).Msg("HelloMysqlOpen is failed to open")
		return err
	}

	if h.MaxIdleCons == 0 {
		h.driver.DB().SetMaxIdleConns(0)
	} else {
		h.driver.DB().SetMaxIdleConns(h.MaxIdleCons)
	}
	return nil
}

// Disconnect is disconnect the connection of mysql.
func (h *HelloMysql) Disconnect() error {
	return h.driver.Close()
}

func (h *HelloMysql) dsn() string {
	cfg := mysql.Config{
		User:                 h.User,
		Passwd:               h.Password,
		DBName:               h.Database,
		Timeout:              time.Second * time.Duration(h.ConnectionTimeout),
		ParseTime:            true,
		AllowNativePasswords: true,
		Params:               make(map[string]string),
	}
	if h.Host != "" {
		if h.Host[0] != '/' {
			cfg.Net = "tcp"
			cfg.Addr = h.Host

			if h.Port != 0 {
				cfg.Addr = net.JoinHostPort(h.Host, strconv.Itoa(int(h.Port)))
			}
		} else {
			cfg.Net = "unix"
			cfg.Addr = h.Host
		}
	}
	return cfg.FormatDSN()
}

// WithHelloMysql option function to assign on adapters.
func WithHelloMysql(driver Driver[*sqlEnt.Driver]) Option {
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
