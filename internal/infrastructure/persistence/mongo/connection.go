package mongo

import (
	"context"

	"github.com/javito2003/shortener_url/internal/config"
	mongoDB "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDatabase() (*mongoDB.Client, error) {
	clientOptions := options.Client().ApplyURI(config.AppConfig.Mongo.URI)
	client, err := mongoDB.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	return client, err
}
