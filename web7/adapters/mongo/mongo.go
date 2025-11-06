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

var ErrTokenNotFound = errors.New("token not found")

type Adapter struct {
	Client *mongo.Client
}

func NewAdapter(ctx context.Context, uri string) (*Adapter, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Adapter{Client: client}, nil
}

type Post struct {
	ID        string    `bson:"_id"`
	Text      string    `bson:"text"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

type User struct {
	ID           string    `bson:"_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
}

type Token struct {
	ID        string    `bson:"_id"`
	Token     string    `bson:"token"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

func (a *Adapter) GetAllPosts(ctx context.Context) ([]Post, error) {
	collection := a.Client.Database("meme9").Collection("posts")

	// Sort by _id in descending order (newest first, ObjectID contains timestamp)
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding posts: %w", err)
	}
	defer cursor.Close(ctx)

	posts := []Post{}
	for cursor.Next(ctx) {
		var post Post
		err = cursor.Decode(&post)
		if err != nil {
			return nil, fmt.Errorf("error decoding post: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (a *Adapter) AddPost(ctx context.Context, post Post) (*Post, error) {
	collection := a.Client.Database("meme9").Collection("posts")

	// Insert without _id field to let MongoDB auto-generate it
	insertDoc := bson.M{
		"text":       post.Text,
		"user_id":    post.UserID,
		"created_at": post.CreatedAt,
	}
	result, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		return nil, fmt.Errorf("error adding post: %w", err)
	}

	objID := result.InsertedID.(primitive.ObjectID)
	post.ID = objID.Hex()
	return &post, nil
}

func (a *Adapter) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	collection := a.Client.Database("meme9").Collection("users")
	var user User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *Adapter) GetUserByID(ctx context.Context, userID string) (*User, error) {
	collection := a.Client.Database("meme9").Collection("users")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	var user User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *Adapter) GetUsersByIDs(ctx context.Context, userIDs []string) (map[string]*User, error) {
	if len(userIDs) == 0 {
		return make(map[string]*User), nil
	}

	collection := a.Client.Database("meme9").Collection("users")

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

func (a *Adapter) CreateUser(ctx context.Context, user User) (*User, error) {
	collection := a.Client.Database("meme9").Collection("users")

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

func (a *Adapter) CreateToken(ctx context.Context, token Token) (*Token, error) {
	collection := a.Client.Database("meme9").Collection("tokens")

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

func (a *Adapter) GetTokenByValue(ctx context.Context, tokenValue string) (*Token, error) {
	collection := a.Client.Database("meme9").Collection("tokens")
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
