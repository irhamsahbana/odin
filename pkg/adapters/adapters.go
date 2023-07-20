// Package adapters are the glue between components and external sources.
package adapters

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	sqlEnt "entgo.io/ent/dialect/sql"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type client interface {
	*sqlEnt.Driver | *http.Client | *resty.Client | *mongo.Client | *kafka.Writer
}

// Driver - interface adapter.
type Driver[T client] interface {
	Open() (T, error)
	Connect() error
	Disconnect() error
}

// Adapter components for external sources.
type Adapter struct {
	HelloMysql    *sqlEnt.Driver
	PokemonResty  *resty.Client
	PokemonRest   *http.Client
	HelloPostgres *sqlEnt.Driver
	HelloSQLite   *sqlEnt.Driver
	HelloMongo    *mongo.Client
	PersistUsers  *mongo.Database
	ProducerHello *kafka.Writer
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
	if a.HelloSQLite != nil {
		log.Info().Msg("HelloSQLite is closed")
		if err := a.HelloSQLite.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if a.HelloPostgres != nil {
		log.Info().Msg("HelloPostgres is closed")
		if err := a.HelloPostgres.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if a.HelloMysql != nil {
		log.Info().Msg("HelloMysql is closed")
		if err := a.HelloMysql.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if a.HelloMongo != nil {
		log.Info().Msg("HelloMongo is closed")
		if err := a.HelloMongo.Disconnect(context.Background()); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}
