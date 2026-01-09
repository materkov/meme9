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

type Like struct {
	ID        string    `bson:"_id"`
	PostID    string    `bson:"post_id"`
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

func (a *Adapter) EnsureIndexes(ctx context.Context) error {
	collection := a.client.Database(a.databaseName).Collection("likes")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "post_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create likes index: %w", err)
	}
	return nil
}

func (a *Adapter) Like(ctx context.Context, postID, userID string) error {
	collection := a.client.Database(a.databaseName).Collection("likes")

	insertDoc := bson.M{
		"post_id":    postID,
		"user_id":    userID,
		"created_at": time.Now(),
	}
	_, err := collection.InsertOne(ctx, insertDoc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// Already liked, return nil (idempotent)
			return nil
		}
		return fmt.Errorf("error creating like: %w", err)
	}
	return nil
}

func (a *Adapter) Unlike(ctx context.Context, postID, userID string) error {
	collection := a.client.Database(a.databaseName).Collection("likes")

	_, err := collection.DeleteOne(ctx, bson.M{
		"post_id": postID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("error deleting like: %w", err)
	}
	return nil
}

func (a *Adapter) IsLiked(ctx context.Context, postID, userID string) (bool, error) {
	collection := a.client.Database(a.databaseName).Collection("likes")

	count, err := collection.CountDocuments(ctx, bson.M{
		"post_id": postID,
		"user_id": userID,
	})
	if err != nil {
		return false, fmt.Errorf("error checking like: %w", err)
	}

	return count > 0, nil
}

func (a *Adapter) GetLikesCounts(ctx context.Context, postIDs []string) (map[string]int32, error) {
	if len(postIDs) == 0 {
		return make(map[string]int32), nil
	}

	collection := a.client.Database(a.databaseName).Collection("likes")

	// Use aggregation to count likes per post
	pipeline := []bson.M{
		{"$match": bson.M{"post_id": bson.M{"$in": postIDs}}},
		{"$group": bson.M{
			"_id":   "$post_id",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error aggregating likes: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	result := make(map[string]int32)
	for _, postID := range postIDs {
		result[postID] = 0 // Initialize all post IDs with 0
	}

	for cursor.Next(ctx) {
		var item struct {
			ID    string `bson:"_id"`
			Count int32  `bson:"count"`
		}
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("error decoding aggregation result: %w", err)
		}
		result[item.ID] = item.Count
	}

	return result, nil
}

func (a *Adapter) GetLikedByUser(ctx context.Context, userID string, postIDs []string) (map[string]bool, error) {
	if len(postIDs) == 0 {
		return make(map[string]bool), nil
	}

	collection := a.client.Database(a.databaseName).Collection("likes")

	cursor, err := collection.Find(ctx, bson.M{
		"user_id": userID,
		"post_id": bson.M{"$in": postIDs},
	})
	if err != nil {
		return nil, fmt.Errorf("error finding likes: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	result := make(map[string]bool)
	for _, postID := range postIDs {
		result[postID] = false // Initialize all post IDs as not liked
	}

	for cursor.Next(ctx) {
		var like Like
		if err := cursor.Decode(&like); err != nil {
			return nil, fmt.Errorf("error decoding like: %w", err)
		}
		result[like.PostID] = true
	}

	return result, nil
}

func (a *Adapter) GetLikers(ctx context.Context, postID, pageToken string, count int) ([]string, string, error) {
	collection := a.client.Database(a.databaseName).Collection("likes")

	filter := bson.M{"post_id": postID}

	if pageToken != "" {
		objID, err := primitive.ObjectIDFromHex(pageToken)
		if err == nil {
			filter["_id"] = bson.M{"$gt": objID}
		}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetLimit(int64(count) + 1)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", fmt.Errorf("error finding likers: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	// Use a struct that can properly decode ObjectID
	type LikeWithObjectID struct {
		ID     primitive.ObjectID `bson:"_id"`
		UserID string             `bson:"user_id"`
	}

	var likes []LikeWithObjectID
	for cursor.Next(ctx) {
		var like LikeWithObjectID
		if err := cursor.Decode(&like); err != nil {
			return nil, "", fmt.Errorf("error decoding like: %w", err)
		}
		likes = append(likes, like)
	}

	nextPageToken := ""
	hasMore := len(likes) > count
	if hasMore {
		likes = likes[:count]

		lastLike := likes[len(likes)-1]
		nextPageToken = lastLike.ID.Hex()
	}

	userIDs := make([]string, 0, len(likes))
	for _, like := range likes {
		userIDs = append(userIDs, like.UserID)
	}

	return userIDs, nextPageToken, nil
}
