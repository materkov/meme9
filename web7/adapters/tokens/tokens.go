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
	client *mongo.Client
}

func NewAdapter(client *mongo.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) Create(ctx context.Context, token Token) (*Token, error) {
	collection := a.client.Database("meme9").Collection("tokens")

	insertDoc := bson.M{
		"token":      token.Token,
		"user_id":    token.UserID,
		"created_at": token.CreatedAt,
	}
	result, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		return nil, fmt.Errorf("error creating token: %w", err)
	}

	objID := result.InsertedID.(primitive.ObjectID)
	token.ID = objID.Hex()
	return &token, nil
}

func (a *Adapter) GetByValue(ctx context.Context, tokenValue string) (*Token, error) {
	collection := a.client.Database("meme9").Collection("tokens")
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
