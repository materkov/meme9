package tokens

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = errors.New("token not found")

type Token struct {
	ID        string    `bson:"_id"`
	Token     string    `bson:"token"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

type Adapter struct {
	client       *mongo.Client
	databaseName string
}

func New(client *mongo.Client, databaseName string) *Adapter {
	return &Adapter{client: client, databaseName: databaseName}
}

func (a *Adapter) Create(ctx context.Context, token Token) (string, error) {
	collection := a.client.Database(a.databaseName).Collection("tokens")

	insertDoc := bson.M{
		"token":      token.Token,
		"user_id":    token.UserID,
		"created_at": token.CreatedAt,
	}
	result, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		return "", fmt.Errorf("error creating token: %w", err)
	}

	objID := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil
}

func (a *Adapter) GetByValue(ctx context.Context, tokenValue string) (*Token, error) {
	collection := a.client.Database(a.databaseName).Collection("tokens")
	var token Token
	err := collection.FindOne(ctx, bson.M{"token": tokenValue}).Decode(&token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error finding token: %w", err)
	}
	return &token, nil
}

