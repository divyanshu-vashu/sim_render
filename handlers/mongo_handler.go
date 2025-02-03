package handlers

import (
    "go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
    db *mongo.Database
}

func NewMongoHandler(db *mongo.Database) *Handler {
    return &Handler{
        db: db,
    }
}