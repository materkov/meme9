package posts

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

type Post struct {
	ID        string    `bson:"_id"`
	Text      string    `bson:"text"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
	Deleted   bool      `bson:"deleted,omitempty"`
}

type Adapter struct {
	client       *mongo.Client
	databaseName string
}

func New(client *mongo.Client, databaseName string) *Adapter {
	return &Adapter{client: client, databaseName: databaseName}
}

func (a *Adapter) GetAll(ctx context.Context) ([]Post, error) {
	collection := a.client.Database(a.databaseName).Collection("posts")

	// Filter out deleted posts
	filter := bson.M{"deleted": bson.M{"$ne": true}}

	// Sort by _id in descending order (newest first, ObjectID contains timestamp)
	opts := options.Find().SetSort(bson.D{bson.E{Key: "_id", Value: -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding posts: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

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

func (a *Adapter) GetByUserIDs(ctx context.Context, userIDs []string) ([]Post, error) {
	if len(userIDs) == 0 {
		return []Post{}, nil
	}

	collection := a.client.Database(a.databaseName).Collection("posts")

	// Filter by user IDs and exclude deleted posts
	filter := bson.M{
		"user_id": bson.M{"$in": userIDs},
		"deleted": bson.M{"$ne": true},
	}

	// Sort by _id in descending order (newest first, ObjectID contains timestamp)
	opts := options.Find().SetSort(bson.D{bson.E{Key: "_id", Value: -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding posts: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

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

var ErrNotFound = errors.New("post not found")

func (a *Adapter) GetByID(ctx context.Context, postID string) (*Post, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, ErrNotFound
	}

	collection := a.client.Database(a.databaseName).Collection("posts")
	var post Post
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

func (a *Adapter) Add(ctx context.Context, post Post) (*Post, error) {
	collection := a.client.Database(a.databaseName).Collection("posts")

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

func (a *Adapter) MarkAsDeleted(ctx context.Context, postID string) error {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil
	}

	collection := a.client.Database(a.databaseName).Collection("posts")
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"deleted": true}},
	)
	if err != nil {
		return fmt.Errorf("failed to mark post as deleted: %w", err)
	}
	return nil
}
