package posts

import (
	"context"
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
}

type Adapter struct {
	client *mongo.Client
}

func New(client *mongo.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) GetAll(ctx context.Context) ([]Post, error) {
	collection := a.client.Database("meme9").Collection("posts")

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

func (a *Adapter) Add(ctx context.Context, post Post) (*Post, error) {
	collection := a.client.Database("meme9").Collection("posts")

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
