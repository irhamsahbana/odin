// Package adapters are the glue between components and external sources.
package adapters

import (
	"fmt"
	"net/http"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type client interface {
	*sql.Driver | *http.Client | *resty.Client
}

// Driver - interface adapter.
type Driver[T client] interface {
	Open() (T, error)
	Connect() error
	Disconnect() error
}

// Adapter components for external sources.
type Adapter struct {
	HelloSqlite *sql.Driver
	HelloMysql  *sql.Driver
}

// Option is Adapter type return func.
type Option func(adapter *Adapter)

// Sync - configure all adapters.
func (a *Adapter) Sync(opts ...Option) {
	for o := range opts {
		opt := opts[o]
		opt(a)
	}
}

// UnSync - release all adapter connection.
func (a *Adapter) UnSync() error {
	var errs []string
	if a.HelloMysql != nil {
		log.Info().Msg("HelloMysql is closed")
		if err := a.HelloMysql.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if a.HelloSqlite != nil {
		log.Info().Msg("HelloSqlite is closed")
		if err := a.HelloSqlite.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}
