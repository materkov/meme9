package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUsernameExists = errors.New("username already exists")
	ErrTokenNotFound  = errors.New("token not found")
)

type User struct {
	ID           string    `bson:"_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"password_hash"`
	AvatarURL    string    `bson:"avatar_url,omitempty"`
	CreatedAt    time.Time `bson:"created_at"`
}

type Token struct {
	ID        string    `bson:"_id"`
	Token     string    `bson:"token"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

type MongoAdapter struct {
	client       *mongo.Client
	databaseName string
}

func New(client *mongo.Client, databaseName string) *MongoAdapter {
	return &MongoAdapter{client: client, databaseName: databaseName}
}

func (a *MongoAdapter) EnsureIndexes(ctx context.Context) error {
	usersCollection := a.client.Database(a.databaseName).Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := usersCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create username index: %w", err)
	}
	return nil
}

func (a *MongoAdapter) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	collection := a.client.Database(a.databaseName).Collection("users")
	var user User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return &user, nil
}

func (a *MongoAdapter) GetUserByID(ctx context.Context, userID string) (*User, error) {
	users, err := a.GetUsersByIDs(ctx, []string{userID})
	if err != nil {
		return nil, err
	}

	user, ok := users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (a *MongoAdapter) GetUsersByIDs(ctx context.Context, userIDs []string) (map[string]*User, error) {
	if len(userIDs) == 0 {
		return make(map[string]*User), nil
	}

	collection := a.client.Database(a.databaseName).Collection("users")

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
	defer func() { _ = cursor.Close(ctx) }()

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

func (a *MongoAdapter) CreateUser(ctx context.Context, user User) (string, error) {
	collection := a.client.Database(a.databaseName).Collection("users")

	insertDoc := bson.M{
		"username":      user.Username,
		"password_hash": user.PasswordHash,
		"created_at":    user.CreatedAt,
	}
	result, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	objID := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil
}

func (a *MongoAdapter) UpdateUserAvatar(ctx context.Context, userID, avatarURL string) error {
	collection := a.client.Database(a.databaseName).Collection("users")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrUserNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"avatar_url": avatarURL,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return fmt.Errorf("error updating avatar: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (a *MongoAdapter) CreateToken(ctx context.Context, token Token) (string, error) {
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

func (a *MongoAdapter) GetTokenByValue(ctx context.Context, tokenValue string) (*Token, error) {
	collection := a.client.Database(a.databaseName).Collection("tokens")
	var token Token
	err := collection.FindOne(ctx, bson.M{"token": tokenValue}).Decode(&token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("error finding token: %w", err)
	}
	return &token, nil
}
