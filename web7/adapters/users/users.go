package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = errors.New("user not found")

type User struct {
	ID           string    `bson:"_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
}

type Adapter struct {
	client *mongo.Client
}

func New(client *mongo.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) GetByUsername(ctx context.Context, username string) (*User, error) {
	collection := a.client.Database("meme9").Collection("users")
	var user User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return &user, nil
}

func (a *Adapter) GetByID(ctx context.Context, userID string) (*User, error) {
	collection := a.client.Database("meme9").Collection("users")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	var user User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return &user, nil
}

func (a *Adapter) GetByIDs(ctx context.Context, userIDs []string) (map[string]*User, error) {
	if len(userIDs) == 0 {
		return make(map[string]*User), nil
	}

	collection := a.client.Database("meme9").Collection("users")

	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(userIDs))
	for _, userID := range userIDs {
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			// Skip invalid IDs
			continue
		}
		objectIDs = append(objectIDs, objID)
	}

	if len(objectIDs) == 0 {
		return make(map[string]*User), nil
	}

	// Query all users at once
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	defer cursor.Close(ctx)

	users := make(map[string]*User)
	for cursor.Next(ctx) {
		var user User
		err = cursor.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("error decoding user: %w", err)
		}
		users[user.ID] = &user
	}

	return users, nil
}

func (a *Adapter) Create(ctx context.Context, user User) (*User, error) {
	collection := a.client.Database("meme9").Collection("users")

	insertDoc := bson.M{
		"username":      user.Username,
		"password_hash": user.PasswordHash,
		"created_at":    user.CreatedAt,
	}
	result, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	objID := result.InsertedID.(primitive.ObjectID)
	user.ID = objID.Hex()
	return &user, nil
}
