package subscriptions

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Subscription struct {
	ID          string    `bson:"_id"`
	FollowerID  string    `bson:"follower_id"`
	FollowingID string    `bson:"following_id"`
	CreatedAt   time.Time `bson:"created_at"`
}

type Adapter struct {
	client       *mongo.Client
	databaseName string
}

func New(client *mongo.Client, databaseName string) *Adapter {
	return &Adapter{client: client, databaseName: databaseName}
}

func (a *Adapter) EnsureIndexes(ctx context.Context) error {
	collection := a.client.Database(a.databaseName).Collection("subscriptions")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "follower_id", Value: 1}, {Key: "following_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create subscription index: %w", err)
	}
	return nil
}

func (a *Adapter) Subscribe(ctx context.Context, followerID, followingID string) error {
	collection := a.client.Database(a.databaseName).Collection("subscriptions")

	// Don't allow self-subscription
	if followerID == followingID {
		return fmt.Errorf("cannot subscribe to yourself")
	}

	insertDoc := bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
		"created_at":   time.Now(),
	}
	_, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// Already subscribed, return nil (idempotent)
			return nil
		}
		return fmt.Errorf("error creating subscription: %w", err)
	}
	return nil
}

func (a *Adapter) Unsubscribe(ctx context.Context, followerID, followingID string) error {
	collection := a.client.Database(a.databaseName).Collection("subscriptions")

	_, err := collection.DeleteOne(ctx, bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	})
	if err != nil {
		return fmt.Errorf("error deleting subscription: %w", err)
	}
	return nil
}

func (a *Adapter) GetFollowing(ctx context.Context, followerID string) ([]string, error) {
	collection := a.client.Database(a.databaseName).Collection("subscriptions")

	cursor, err := collection.Find(ctx, bson.M{"follower_id": followerID})
	if err != nil {
		return nil, fmt.Errorf("error finding subscriptions: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	followingIDs := []string{}
	for cursor.Next(ctx) {
		var sub Subscription
		err = cursor.Decode(&sub)
		if err != nil {
			return nil, fmt.Errorf("error decoding subscription: %w", err)
		}
		followingIDs = append(followingIDs, sub.FollowingID)
	}

	return followingIDs, nil
}

func (a *Adapter) IsSubscribed(ctx context.Context, followerID, followingID string) (bool, error) {
	collection := a.client.Database(a.databaseName).Collection("subscriptions")

	count, err := collection.CountDocuments(ctx, bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	})
	if err != nil {
		return false, fmt.Errorf("error checking subscription: %w", err)
	}

	return count > 0, nil
}
