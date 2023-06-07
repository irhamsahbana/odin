// Package adapters provides adapters between the application and external services.
package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/infrastructure"
)

// MockMongoClient is a mock implementation of the mongo.Client interface.
type MockMongoClient struct {
	mock.Mock
}

// TestWithHelloMongo tests the WithHelloMongo function.
func TestWithHelloMongo(t *testing.T) {
	// Create a new assertion object using the testify library.
	is := assert.New(t)

	// Create a new mock mongo client.
	mockClient := &MockMongoClient{}

	// Create a new adapter instance.
	adapter := &Adapter{}

	// define env
	infrastructure.Configuration(
		infrastructure.WithPath("../.."),
		infrastructure.WithFilename("config.yaml"),
	).Initialize()

	// Call the adapter's Sync method with an Option object created by calling the WithHelloMongo function.
	adapter.Sync(
		WithHelloMongo(&HelloMongo{
			NetworkDB: NetworkDB{
				Host:     infrastructure.Envs.HelloMongo.Host,
				Port:     infrastructure.Envs.HelloMongo.Port,
				Database: infrastructure.Envs.HelloMongo.Database,
				User:     infrastructure.Envs.HelloMongo.User,
				Password: infrastructure.Envs.HelloMongo.Password,
			},
		}),
	)

	// Call the adapter's UnSync method and check that the error return value is nil.
	is.Nil(adapter.UnSync())

	// Check that all expectations set on the mockClient have been met.
	mockClient.AssertExpectations(t)
}
