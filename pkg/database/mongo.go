package database

import (
	"context"
	"log"
	"shifty-backend/configs"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB(cfg *configs.Config) *mongo.Client {

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Erro when create Mongo Client: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Can not ping to MongoDB: ", err)
	}
	log.Println("Connect to MongoDB successful!")
	return client
}
