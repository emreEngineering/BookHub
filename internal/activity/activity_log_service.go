package activity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActivityLogger interface {
	Log(ctx context.Context, eventType string, message string, userID *int, metadata map[string]interface{}) error
	FindLatest(ctx context.Context, limit int64) ([]ActivityLog, error)
}

type MongoActivityLogService struct {
	collection *mongo.Collection
}

func NewMongoActivityLogService(db *mongo.Database) *MongoActivityLogService {
	return &MongoActivityLogService{
		collection: db.Collection("activity_logs"),
	}
}

func (s *MongoActivityLogService) Log(ctx context.Context, eventType string, message string, userID *int, metadata map[string]interface{}) error {
	if metadata == nil {
		metadata = map[string]interface{}{}
	}

	log := ActivityLog{
		Type:      eventType,
		Message:   message,
		UserID:    userID,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, log)
	return err
}

func (s *MongoActivityLogService) FindLatest(ctx context.Context, limit int64) ([]ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit)

	cursor, err := s.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	logs := []ActivityLog{}
	for cursor.Next(ctx) {
		var log ActivityLog
		err := cursor.Decode(&log)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
