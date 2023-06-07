package adapters

import (
	"context"
	"fmt"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HelloMongo is data of instances.
type HelloMongo struct {
	NetworkDB
	driver *mongo.Client
}

// Open returns the client object for MongoDB database.
func (h *HelloMongo) Open() (*mongo.Client, error) {
	if h.driver == nil {
		return nil, fmt.Errorf("driver was failed to connected")
	}
	return h.driver, nil
}

// Connect establishes a connection to the MongoDB server.
func (h *HelloMongo) Connect() (err error) {
	clientOptions := options.Client().ApplyURI(h.dsn())

	if infrastructure.Envs.HelloMongo.AuthEnable {
		clientOptions.Auth = &options.Credential{
			Username: h.User,
			Password: h.Password,
		}
	}

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to create client for MongoDB database: %v", err)
	}

	h.driver = client
	return nil
}

// Disconnect closes the connection to the MongoDB server.
func (h *HelloMongo) Disconnect() error {
	return h.driver.Disconnect(context.Background())
}

// dsn returns the connection string for the MongoDB database.
func (h *HelloMongo) dsn() string {
	return fmt.Sprintf("mongodb://%s:%d/%s", h.Host, h.Port, h.Database)
}

// WithHelloMongo option function to assign on adapters.
func WithHelloMongo(driver Driver[*mongo.Client]) Option {
	return func(a *Adapter) {
		if err := driver.Connect(); err != nil {
			panic(err)
		}

		open, err := driver.Open()
		if err != nil {
			panic(err)
		}

		a.HelloMongo = open
		// define connection to the database
		a.PersistUsers = open.Database(driver.(*HelloMongo).Database)
	}
}
