package config

import (
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(GetMongoURI()))
    if err != nil {
        log.Fatal(err)
    }

    // Ping the database
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    MongoClient = client
    log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
    return MongoClient.Database("sim_render").Collection(collectionName)
}