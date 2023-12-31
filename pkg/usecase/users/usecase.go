// Package users implement all logic.
package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/entity"
)

func (i *impl) Get(ctx context.Context, request entity.RequestGetUsers) (result entity.ResponseGetUsers, err error) {
	coll := i.adapter.PersistUsers.Collection("users")

	skip := (request.Page - 1) * request.Limit

	// Query options with skip and limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(request.Limit))

	// pagination
	result.Limit = request.Limit
	result.Page = request.Page
	var cursor *mongo.Cursor
	cursor, err = coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return result, err
	}
	defer func(c context.Context) {
		err = cursor.Close(c)
	}(ctx)

	// Iterate through the cursor to get each document.
	var documents = make([]entity.User, 0)
	for cursor.Next(ctx) {
		var document entity.User
		err := cursor.Decode(&document)
		if err != nil {
			return result, err
		}
		documents = append(documents, document)
	}

	result.Users = documents
	return result, nil
}

func (i *impl) Create(ctx context.Context, user entity.User) (entity.User, error) {
	coll := i.adapter.PersistUsers.Collection("users")

	user.CreatedAt = time.Now()

	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	// Retrieve the created document using the _id from the InsertOneResult
	var createdUser entity.User
	err = coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdUser)
	if err != nil {
		return entity.User{}, err
	}

	return createdUser, nil
}

func (i *impl) Update(ctx context.Context) error {
	uuid := uuid.New().String()

	i.adapter.ProducerHello.WriteMessages(ctx, kafka.Message{
		Key:   []byte("this-is-key-" + uuid),
		Value: []byte("this-is-value"),
		Time:  time.Now(),
	})

	return nil
}
