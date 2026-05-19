package activity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityLog struct {
	ID        primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Type      string                 `json:"type" bson:"type"`
	Message   string                 `json:"message" bson:"message"`
	UserID    *int                   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
}
