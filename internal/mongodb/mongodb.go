package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect(parentCtx context.Context, mongoURI string, mongoDatabase string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("MongoDB bağlantısı başarılı")
	return client, client.Database(mongoDatabase), nil
}
