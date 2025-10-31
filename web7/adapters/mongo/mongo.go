package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	CreatedAt time.Time `bson:"created_at"`
}

func (a *Adapter) GetAllPosts(ctx context.Context) ([]Post, error) {
	collection := a.Client.Database("meme9").Collection("posts")
	cursor, err := collection.Find(ctx, bson.D{})
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
